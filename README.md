# zolly

Zolly API Gateway
- High Performance
- Lightweight
- Simple & Easy
- Auto SSL with auto generate key
- gRPC Plugin System (HashiCorp go-plugin)


Install
=======
``` sh
go install github.com/dreamph/zolly@latest
```


Configuration (Basic)
=======
``` yml
server:
  port: 3000
  bodyLimit: 100
  cors:
    enable: true
  log:
    enable: false
services:
  - path: "/orders"
    stripPath: true
    timeout: 60
    servers:
      - "http://localhost:3001"
  - path: "/products"
    stripPath: true
    timeout: 60
    servers:
      - "http://localhost:3002"
```

Configuration (Auto SSL)
=======
``` yml
server:
  port: 3000
  bodyLimit: 100
  cors:
    enable: true
  log:
    enable: false
  ssl:
    enable: true
    generateKey:
      enable: true
      keyConfig:
        commonName: "localhost"
        file: "./certs/server.p12"
        password: "password"
services:
  - path: "/orders"
    stripPath: true
    timeout: 60
    servers:
      - "http://localhost:3001"
  - path: "/products"
    stripPath: true
    timeout: 60
    servers:
      - "http://localhost:3002"
```

Configuration (With Plugins)
=======
``` yml
server:
  port: 3000
  cors:
    enable: true
  log:
    enable: false

plugins:
  - name: "auth"
    path: "./plugins/auth-plugin"
    settings:
      api_key: "my-secret-key"
      header_name: "X-API-Key"

services:
  - path: "/orders"
    stripPath: true
    timeout: 60
    plugins:
      - "auth"
    servers:
      - "http://localhost:3001"
  - path: "/products"
    stripPath: true
    timeout: 60
    plugins:
      - "auth"
      - "hmac"
    servers:
      - "http://localhost:3002"
```

Start Server
=======
``` sh
zolly -c config.yml
```

Plugins
=======

Zolly supports middleware plugins via [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin) (gRPC). Each plugin is a separate Go binary that communicates with the gateway over gRPC.

### How It Works

```
Client Request
  -> Plugin 1 (gRPC)
    -> Plugin 2 (gRPC)
      -> Proxy to backend
  -> Response
```

- Plugins are defined in `config.yml` with a name, binary path, and settings
- Each service can specify which plugins to apply (in order)
- Plugins can **continue** (pass request to next handler) or **abort** (return response immediately, e.g. 401)

### Writing a Plugin

``` go
package main

import (
	"context"

	zollyplugin "github.com/dreamph/zolly/plugin"
)

type MyPlugin struct{}

func (p *MyPlugin) Configure(_ context.Context, settings *zollyplugin.PluginSettings) (*zollyplugin.ConfigureResult, error) {
	// Read settings and initialize plugin
	return &zollyplugin.ConfigureResult{Success: true}, nil
}

func (p *MyPlugin) HandleRequest(_ context.Context, req *zollyplugin.Request) (*zollyplugin.Response, error) {
	// Inspect request and decide: continue or abort
	return &zollyplugin.Response{Action: zollyplugin.ActionContinue}, nil
}

func (p *MyPlugin) HealthCheck(_ context.Context) (*zollyplugin.HealthCheckResult, error) {
	return &zollyplugin.HealthCheckResult{Healthy: true, Message: "ok"}, nil
}

func main() {
	zollyplugin.Serve(&MyPlugin{})
}
```

### Plugin Interface

| Method | Description |
|---|---|
| `Configure` | Called once at startup with settings from config |
| `HandleRequest` | Called for every request. Return `ActionContinue` or `ActionAbort` |
| `HealthCheck` | Called periodically to verify plugin health |

### Request / Response

**Request** fields: `Method`, `Path`, `OriginalURL`, `Headers`, `Body`, `QueryParams`, `ClientIP`

**Response** fields:
- `Action` - `"continue"` (pass to next handler) or `"abort"` (return response)
- `StatusCode` - HTTP status code (used when abort)
- `Headers` - Response headers (abort) or modified request headers (continue)
- `Body` - Response body (abort) or modified request body (continue)

### Example Plugins

| Plugin | Description |
|---|---|
| [auth](examples/plugins/auth) | API key authentication via header |
| [hmac](examples/plugins/hmac) | HMAC-SHA256 signature validation |

### Build Plugins

``` sh
make build-plugin-auth    # Build auth plugin
make build-plugin-hmac    # Build hmac plugin
make build-plugins        # Build all plugins
```

### Benchmark
- MacBook Pro 2023
- Chip Apple M2 Max
- Memory 64 GB

[Benchmark Source](https://github.com/dreamph/zolly-bench)

### LoadTest without API Gateway

```shell
wrk -t12 -c100 -d60s http://127.0.0.1:8000/v1/hello
Running 1m test @ http://127.0.0.1:8000/v1/hello
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.06ms    5.16ms 229.90ms   98.87%
    Req/Sec    11.76k     2.72k   68.56k    92.35%
  8423482 requests in 1.00m, 1.03GB read
Requests/sec: 140187.71
Transfer/sec:     17.51MB
```

### LoadTest Zolly API Gateway

```shell
wrk -t12 -c100 -d60s http://127.0.0.1:8070/v1/hello
Running 1m test @ http://127.0.0.1:8070/v1/hello
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.63ms   16.56ms 376.43ms   95.09%
    Req/Sec     6.40k     2.90k   34.92k    82.16%
  4579541 requests in 1.00m, 572.13MB read
  Socket errors: connect 0, read 0, write 0, timeout 1
Requests/sec:  76203.34
Transfer/sec:      9.52MB
```

### LoadTest KrakenD API Gateway

```shell
wrk -t12 -c100 -d60s http://127.0.0.1:8090/v1/hello
Running 1m test @ http://127.0.0.1:8090/v1/hello
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.02ms    2.01ms  70.55ms   94.56%
    Req/Sec     4.47k   416.06     6.27k    70.93%
  3200251 requests in 1.00m, 717.22MB read
Requests/sec:  53322.81
Transfer/sec:     11.95MB
```

Buy Me a Coffee
=======
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/dreamph)
