dev:
	@gin --path $(pwd)\src --all -i run main.go
dev-windows:
	@gin --path ${CURDIR}\src --all -i run main.go
release-windows:
	@go build -o bin/main.exe src/main.go
release:
	@go build -o bin/main src/main.go
install:
	@go mod vendor