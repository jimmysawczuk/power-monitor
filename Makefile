define setup
	rm -rf web/static/bower web/static/css web/static/js/bin node_modules
	go get -u github.com/jimmysawczuk/scm-status/...
	go get -u github.com/jteeuwen/go-bindata/...

	yarn
	bower install
	grunt
endef

define clean
	@rm -f web/static.go
endef

define build
	@echo 'Building...'

	scm-status -out=web/static/REVISION.json

	go-bindata -debug -o web/static.go -pkg=web web/templates/... web/static/...
	go-bindata -debug -o tls.go -pkg=main tls/...
	go install -tags="debug" .
endef

define release
	@echo 'Building (release)...'

	scm-status -out=web/static/REVISION.json
	grunt

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

release:
	@$(clean)
	@$(release)

tls:
	mkdir -p tls
	openssl req -newkey rsa:2048 -nodes -keyout tls/key.pem -x509 -days 3652 -out tls/certificate.pem -config ./tlsconfig.conf
