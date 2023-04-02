.PHONY: build

build:
	@if [ ! -d "build" ]; then mkdir "build"; fi

	@if [ ! -d "build/windows" ]; then mkdir "build/windows"; fi
	@if [ ! -d "build/windows/amd64" ]; then mkdir "build/windows/amd64"; fi
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/gitsum.exe ./cmd/gitsum
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/gitsum-importer.exe ./cmd/importer
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/gitsum-exporter.exe ./cmd/exporter

	@if [ ! -d "build/linux" ]; then mkdir "build/linux"; fi
	@if [ ! -d "build/linux/amd64" ]; then mkdir "build/linux/amd64"; fi
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/gitsum ./cmd/gitsum
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/gitsum-importer ./cmd/importer
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/gitsum-exporter ./cmd/exporter
