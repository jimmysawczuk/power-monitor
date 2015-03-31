export PATH := $(shell npm bin):$(PATH)
export GOPATH := ${PWD}

define reset
	@rm -rf bin pkg
	@mkdir -p bin
	@mkdir -p pkg
endef

define fmt
	@echo 'Running gofmt...';
	find . -type f -name "*.go" | xargs gofmt -w
endef

define build
	@echo 'Building...'

	go install power-monitor
	scm-status -out=REVISION.json
	cp REVISION.json src/web/static/
endef

define release
	@echo 'Packaging assets for release...'
	grunt
endef

define reset-deps
	@echo 'Clearing src/github.com/*, node_modules/*, cached web resources...'
	rm -rf src/github.com
	rm -rf node_modules
	rm -rf src/web/static/bower
	rm -rf src/web/static/css
	rm -rf src/web/static/js/bin
endef

define deps
	@echo 'Getting dependencies'

	@echo '  github.com/gin-gonic/gin'
	go get -u github.com/gin-gonic/gin

	@echo '  github.com/gin-gonic/contrib/static'
	go get -u github.com/gin-gonic/contrib/static

	@echo '  github.com/gin-gonic/contrib/gzip'
	go get -u github.com/gin-gonic/contrib/gzip
endef

default: dev

release:
	@$(reset)
	@$(fmt)
	@$(release)
	@$(build)

dev:
	@$(reset)
	@$(fmt)
	@$(build)

fmt:
	@$(reset)
	@$(fmt)

setup:
	@$(reset-deps)
	@$(deps)
	@npm install
	@bower install
