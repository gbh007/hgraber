create_build_dir:
	mkdir -p ./_build

build: create_build_dir
	GOOS=linux GOARCH=arm64 go build -o ./_build/hgraber-linux-arm64
	GOOS=linux GOARCH=amd64 go build -o ./_build/hgraber-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o ./_build/hgraber-windows-amd64.exe
	tar -C ./_build -cf ./_build/hgraber.tar hgraber-linux-arm64 hgraber-linux-amd64 hgraber-windows-amd64.exe
	tar -rf ./_build/hgraber.tar ./static
	gzip -9f ./_build/hgraber.tar

build_arm64: create_build_dir
	# CGO_ENABLED=0 
	GOOS=linux GOARCH=arm64 go build -o ./_build/hgraber-arm64

build_amd64: create_build_dir
	# CGO_ENABLED=0 
	GOOS=linux GOARCH=amd64 go build -o ./_build/hgraber-amd64
	GOOS=windows GOARCH=amd64 go build -o ./_build/hgraber-amd64.exe

run: build_amd64
	./_build/hgraber-amd64 -v