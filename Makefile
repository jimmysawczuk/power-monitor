define setup
	go get -u github.com/jimmysawczuk/scm-status/...
	go get -u github.com/jteeuwen/go-bindata/...
endef

define clean
	@rm -f web/static.go
endef

define build
	@echo 'Building...'
	go fmt ./...

	# bower install
	# grunt
	scm-status -out=web/static/REVISION.json

	go-bindata -debug -o web/static.go -pkg=web web/templates/... web/static/...
	go install .
endef

define release
	@echo 'Building (release)...'

	go fmt ./...

	bower install
	grunt
	scm-status -out=web/static/REVISION.json

	go-bindata -o web/static.go -pkg=web web/templates/... web/static/...
	go install .
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
