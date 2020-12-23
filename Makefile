default: build

prepare:
	docker build -t bexs/testbuild -f Dockerfile . --no-cache

## build: build application executable
.PHONY: build
build: prepare
	docker run   -v $(PWD)/dist:/app/dist -t bexs/testbuild go build -o /app/dist/service ./main.go
## test: run all tests
.PHONY: test
test: prepare
	docker run -t bexs/testbuild  go test -v ./...

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo