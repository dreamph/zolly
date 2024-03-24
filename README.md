# zolly
[![GoDoc](https://godoc.org/github.com/imroc/req?status.svg)](https://godoc.org/github.com/imroc/req)

Zolly Api gateway
- High Performance
- Light weight
- Simple & Easy


Install
=======
``` sh
go install github.com/dreamph/zolly@latest
```


Configuration
=======
``` yml
server:
  port: 3000
services:
  - path: "/orders"
    stripPath: true
    timeout: 60
    servers:
      - "http://127.0.0.1:3001"
      - "http://127.0.0.1:3002"
```

Start Server
=======
``` sh
zolly -c config.yml
```

Load Test Direct to Backend
=======
``` sh
wrk -t12 -c100 -d60s http://127.0.0.1:3001
Running 1m test @ http://127.0.0.1:3001
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   644.53us  567.74us  42.95ms   95.24%
    Req/Sec    12.25k     1.25k   20.55k    75.78%
  8776698 requests in 1.00m, 1.10GB read
Requests/sec: 146270.41
Transfer/sec:     18.83MB
```

Load Test with Zolly API Gateway
=======
``` sh
wrk -t12 -c100 -d60s http://127.0.0.1:3000/orders
Running 1m test @ http://127.0.0.1:3000/orders
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.80ms    3.11ms  82.80ms   95.45%
    Req/Sec     6.01k     0.94k   12.34k    64.68%
  4315660 requests in 1.00m, 555.62MB read
Requests/sec:  71805.77
Transfer/sec:      9.24MB
```

Buy Me a Coffee
=======
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/dreamph)