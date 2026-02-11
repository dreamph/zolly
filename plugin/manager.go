package plugin

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
)

// PluginInstance represents a single running plugin process.
type PluginInstance struct {
	Name   string
	Client *goplugin.Client
	Plugin MiddlewarePlugin
}

// Manager manages the lifecycle of all plugin processes.
type Manager struct {
	mu        sync.RWMutex
	instances map[string]*PluginInstance
	logger    hclog.Logger
}

// NewManager creates a new plugin manager.
func NewManager() *Manager {
	return &Manager{
		instances: make(map[string]*PluginInstance),
		logger: hclog.New(&hclog.LoggerOptions{
			Name:  "zolly-plugin",
			Level: hclog.Info,
		}),
	}
}

// LoadPlugin starts a plugin binary, performs the handshake, and calls Configure.
func (m *Manager) LoadPlugin(name string, binaryPath string, settings map[string]string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.instances[name]; exists {
		return fmt.Errorf("plugin %q already loaded", name)
	}

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          PluginMap,
		Cmd:              exec.Command(binaryPath),
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		Logger:           m.logger,
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("connect to plugin %q: %w", name, err)
	}

	raw, err := rpcClient.Dispense(PluginName)
	if err != nil {
		client.Kill()
		return fmt.Errorf("dispense plugin %q: %w", name, err)
	}

	p, ok := raw.(MiddlewarePlugin)
	if !ok {
		client.Kill()
		return fmt.Errorf("plugin %q does not implement MiddlewarePlugin", name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := p.Configure(ctx, &PluginSettings{Settings: settings})
	if err != nil {
		client.Kill()
		return fmt.Errorf("configure plugin %q: %w", name, err)
	}
	if !result.Success {
		client.Kill()
		return fmt.Errorf("plugin %q configuration failed: %s", name, result.ErrorMessage)
	}

	m.instances[name] = &PluginInstance{
		Name:   name,
		Client: client,
		Plugin: p,
	}

	log.Printf("[plugin] Loaded: %s (%s)\n", name, binaryPath)
	return nil
}

// GetPlugin returns a loaded plugin by name.
func (m *Manager) GetPlugin(name string) (MiddlewarePlugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inst, exists := m.instances[name]
	if !exists {
		return nil, fmt.Errorf("plugin %q not found", name)
	}
	return inst.Plugin, nil
}

// Shutdown kills all plugin processes.
func (m *Manager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, inst := range m.instances {
		log.Printf("[plugin] Stopping: %s\n", name)
		inst.Client.Kill()
	}
	m.instances = make(map[string]*PluginInstance)
}
