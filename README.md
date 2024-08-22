# zolly

Zolly Api gateway
- High Performance
- Light weight
- Simple & Easy
- Auto SSL with auto generate key


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

Start Server
=======
``` sh
zolly -c config.yml
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