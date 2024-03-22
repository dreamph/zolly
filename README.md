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
  - path: "/users"
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


Buy Me a Coffee
=======
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/dreamph)