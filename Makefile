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
endef

define reset-deps
	@echo 'Clearing src/github.com/*...'
	rm -rf src/github.com
endef

define deps
	@echo 'Getting dependencies'

	@echo '  github.com/gin-gonic/gin'
	go get -u github.com/gin-gonic/gin
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

