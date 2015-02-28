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

default: dev

dev:
	@$(reset)
	@$(fmt)
	@$(build)

fmt:
	@$(reset)
	@$(fmt)
