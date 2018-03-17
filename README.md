# power-monitor
[![Go Report Card](https://goreportcard.com/badge/github.com/jimmysawczuk/power-monitor)](https://goreportcard.com/report/github.com/jimmysawczuk/power-monitor)

![Screenshot](https://i.imgur.com/9FDRr83.png)

## Setting up
You'll need to set up [PowerPanel for Linux](https://www.cyberpowersystems.com/product/software/powerpanel-for-linux/) on a server that's connected via USB to a CyberPower UPS.

You'll also need Go >= 1.9 and npm, as well as a globally installed verison of yarn and parcel (`npm install -g yarn parcel-bundler`).

```bash
$ go get -u github.com/jimmysawczuk/power-monitor
$ cd $GOPATH/src/github.com/jimmysawczuk/power-monitor
$ make setup
$ make tls
$ make release
$ power-monitor
```

Finally, open a browser and navigate to `https://your-server:3000/` and you should be good to go.
