TAG = $(shell git tag -l --points-at HEAD)
COMMIT = $(shell git show -s --abbrev=12 --pretty=format:%h HEAD)
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

LDFLAGS = -ldflags "-X 'app/system.Version=$(TAG)' -X 'app/system.Commit=$(COMMIT)' -X 'app/system.BuildAt=$(BUILD_TIME)'"
LDFLAGS_CGO = -ldflags "-linkmode external -extldflags -static -X 'app/system.Version=$(TAG)' -X 'app/system.Commit=$(COMMIT)' -X 'app/system.BuildAt=$(BUILD_TIME)'"

create_build_dir:
	mkdir -p ./_build

build: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./_build/hgraber-linux-arm64 ./cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-linux-amd64 ./cmd/server
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-windows-amd64.exe ./cmd/server

	CC=/usr/bin/aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build $(LDFLAGS_CGO) -o ./_build/hgraber-linux-arm64-cgo ./cmd/server
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(LDFLAGS_CGO) -o ./_build/hgraber-linux-amd64-cgo ./cmd/server
	CC=/usr/bin/x86_64-w64-mingw32-gcc  CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build $(LDFLAGS_CGO) -o ./_build/hgraber-windows-amd64-cgo.exe ./cmd/server

	tar -C ./_build -cf ./_build/hgraber.tar hgraber-linux-arm64 hgraber-linux-amd64 hgraber-windows-amd64.exe hgraber-linux-arm64-cgo hgraber-linux-amd64-cgo hgraber-windows-amd64-cgo.exe
	gzip -9f ./_build/hgraber.tar

build_arm64: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./_build/hgraber-arm64 ./cmd/server

	CC=/usr/bin/aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build $(LDFLAGS_CGO) -o ./_build/hgraber-linux-arm64-cgo ./cmd/server

build_amd64: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-amd64 ./cmd/server
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-amd64.exe ./cmd/server

	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(LDFLAGS_CGO) -o ./_build/hgraber-linux-amd64-cgo ./cmd/server
	CC=/usr/bin/x86_64-w64-mingw32-gcc  CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build $(LDFLAGS_CGO) -o ./_build/hgraber-windows-amd64-cgo.exe ./cmd/server

run: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin  ./cmd/server
	./_build/hgraber-bin -p 8081
	
view: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin  ./cmd/server
	./_build/hgraber-bin -v -p 8081

debug: create_build_dir
	go build $(LDFLAGS) -trimpath -o ./_build/hgraber-bin  ./cmd/server
	./_build/hgraber-bin -stdfile-append -debug -debug-fullpath -h 127.0.0.1 -p 8081 -static="internal/service/webServer/static" --access-token=local-debug

demo: create_build_dir
	go build $(LDFLAGS) -trimpath -o ./_build/hgraber-bin  ./cmd/inmemory
	./_build/hgraber-bin -debug -h 127.0.0.1 -p 8081 --access-token=local-debug
	