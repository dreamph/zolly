# zolly

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
    Latency     6.27ms   21.91ms 384.45ms   93.90%
    Req/Sec     8.16k     4.82k  107.65k    87.96%
  5801723 requests in 1.00m, 746.95MB read
  Socket errors: connect 0, read 0, write 0, timeout 1
Requests/sec:  96542.32
Transfer/sec:     12.43MB
```

Buy Me a Coffee
=======
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/dreamph)