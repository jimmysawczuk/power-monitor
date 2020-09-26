define setup
	go get -u github.com/jimmysawczuk/scm-status/...
	go get -u github.com/jimmysawczuk/tmpl/...
	go get -u github.com/go-bindata/go-bindata/...
endef

define build
	@echo 'Building...'

	scm-status -out=web/static/REVISION.json
	npm run build
	go-bindata -debug -o web/static.go -pkg=web web/templates/... web/static/...
	go-bindata -debug -o cmd/power-monitor/tls.go -pkg=main tls/...
	go install -tags="debug" ./...
endef

define release
	@echo 'Building (release)...'
	./deploy/release.bash
endef

default: dev

setup:
	@$(setup)

dev: tls
	$(build)

release:
	$(release)

tls:
	# Uses https://github.com/FiloSottile/mkcert
	mkdir -p tls
	mkcert localhost
	mv localhost.pem tls/certificate.pem
	mv localhost-key.pem tls/key.pem
