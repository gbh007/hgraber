create_build_dir:
	mkdir -p ./_build

build_arm64: create_build_dir
	# CGO_ENABLED=0 
	GOOS=linux GOARCH=arm64 go build -o ./_build/hgraber-arm64

build_amd64: create_build_dir
	# CGO_ENABLED=0 
	GOOS=linux GOARCH=amd64 go build -o ./_build/hgraber-amd64

run: build_amd64
	./_build/hgraber-amd64 -v