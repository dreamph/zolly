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
wrk -t12 -c400 -d30s http://127.0.0.1:3001
Running 30s test @ http://127.0.0.1:3001
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.50ms    2.01ms 253.14ms   98.97%
    Req/Sec    12.74k     3.70k   44.26k    80.08%
  4567517 requests in 30.02s, 588.05MB read
  Socket errors: connect 155, read 109, write 0, timeout 0
Requests/sec: 152142.87
Transfer/sec:     19.59MB
```

Load Test with Zolly API Gateway
=======
``` sh
wrk -t12 -c400 -d30s http://127.0.0.1:3000/orders
Running 30s test @ http://127.0.0.1:3000/orders
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.97ms   29.01ms 625.74ms   93.84%
    Req/Sec     8.14k     2.85k   23.06k    70.18%
  2925285 requests in 30.10s, 376.62MB read
  Socket errors: connect 155, read 106, write 0, timeout 0
Requests/sec:  97173.74
Transfer/sec:     12.51MB
```

Buy Me a Coffee
=======
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/dreamph)