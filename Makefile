define setup
	go get -u github.com/jimmysawczuk/scm-status/...
	go get -u github.com/jteeuwen/go-bindata/...
	yarn
endef

define clean
	rm -rf web/static/bin web/static.go tls.go tls/ deploy/
endef

define build
	@echo 'Building...'

	scm-status -out=web/static/REVISION.json
	parcel build --detailed-report --global PowerMonitor -o app.js --public-url /bin -d web/static/bin --no-minify ./web/static/app.js
	go-bindata -debug -o web/static.go -pkg=web web/templates/... web/static/...
	go-bindata -debug -o cmd/power-monitor/tls.go -pkg=main tls/...
	go install -tags="debug" ./...
endef

define release
	@echo 'Building (release)...'

	scm-status -out=web/static/REVISION.json
	parcel build --global PowerMonitor -o app.js --public-url /bin -d web/static/bin ./web/static/app.js
	go-bindata -o web/static.go -pkg=web web/templates/... web/static/...
	go-bindata -o cmd/power-monitor/tls.go -pkg=main tls/...
	go build -tags="release" ./...
endef

default: dev

setup:
	@$(setup)

dev: tls
	$(build)

clean:
	$(clean)

release: clean tls
	$(release)

tls:
	# Uses https://github.com/FiloSottile/mkcert
	mkdir -p tls
	mkcert localhost
	mv localhost.pem tls/certificate.pem
	mv localhost-key.pem tls/key.pem
