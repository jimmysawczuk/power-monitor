# power-monitor
[![Go Report Card](https://goreportcard.com/badge/github.com/jimmysawczuk/power-monitor)](https://goreportcard.com/report/github.com/jimmysawczuk/power-monitor)

![Screenshot](https://i.imgur.com/Z2q8Ry0.png)

## Setting up
You'll need to set up [PowerPanel for Linux](https://www.cyberpowersystems.com/product/software/powerpanel-for-linux/) on a server that's connected via USB to a CyberPower UPS.

You'll also need Go >= 1.8 and npm, as well as a globally installed version of bower, yarn and grunt (`npm install -g yarn bower grunt-cli`).

```bash
$ go get -u github.com/jimmysawczuk/power-monitor
$ cd $GOPATH/src/github.com/jimmysawczuk/power-monitor
$ make setup
$ make release
$ power-monitor
```

Finally, open a browser and navigate to `http://your-server:3000/` and you should be good to go.

## License

```
The MIT License (MIT)

Copyright (c) 2015-2017 Jimmy Sawczuk

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
