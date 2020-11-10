dev:
	@gin --path ${CURDIR}\src --all -i run main.go
release:
	@go build -o bin/main.exe src/main.go