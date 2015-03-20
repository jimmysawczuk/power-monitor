GOPATH := ${PWD}

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

define reset-deps
	@echo 'Clearing src/github.com/*...'
	rm -rf src/github.com
	rm -rf src/web/static/bower
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
	@bower install

