define setup
	go get -u github.com/jimmysawczuk/scm-status/...
	go get -u github.com/jteeuwen/go-bindata/...
	yarn
endef

define clean
	@rm -rf web/static/bin
	@rm -f web/static.go
endef

define build
	@echo 'Building...'

	scm-status -out=web/static/REVISION.json

	parcel build --global PowerMonitor -o app.js --public-url /bin -d web/static/bin --no-minify ./web/static/app.js
	go-bindata -debug -o web/static.go -pkg=web web/templates/... web/static/...
	go-bindata -debug -o tls.go -pkg=main tls/...
	go install -tags="debug" .
endef

define release
	@echo 'Building (release)...'

	scm-status -out=web/static/REVISION.json
	parcel build --global PowerMonitor -o app.js --public-url /bin -d web/static/bin ./web/static/app.js
	go-bindata -o web/static.go -pkg=web web/templates/... web/static/...
	go-bindata -debug -o tls.go -pkg=main tls/...
	go install -tags="release" .
endef

default: dev

setup:
	@$(setup)

dev:
	@$(clean)
	@$(build)

clean:
	@$(clean)

release:
	@$(clean)
	@$(release)

tls:
	mkdir -p tls
	mkcert localhost
	mv localhost.pem tls/certificate.pem
	mv localhost-key.pem tls/key.pem
