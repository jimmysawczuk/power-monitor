define build
	@echo 'Building...'

	go install .
	scm-status -out=web/static/REVISION.json
endef

define release
	@echo 'Packaging assets for release...'
	grunt
endef

default: dev

dev:
	@$(build)

release:
	@$(release)
	@$(build)
