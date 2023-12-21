TAG = $(shell git tag -l --points-at HEAD)
COMMIT = $(shell git show -s --abbrev=12 --pretty=format:%h HEAD)
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

LDFLAGS = -ldflags "-X 'app/system.Version=$(TAG)' -X 'app/system.Commit=$(COMMIT)' -X 'app/system.BuildAt=$(BUILD_TIME)'"

create_build_dir:
	mkdir -p ./_build

build: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./_build/hgraber-linux-arm64 ./cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-linux-amd64 ./cmd/server
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-windows-amd64.exe ./cmd/server

	tar -C ./_build -cf ./_build/hgraber.tar hgraber-linux-arm64 hgraber-linux-amd64 hgraber-windows-amd64.exe
	gzip -9f ./_build/hgraber.tar

run: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin  ./cmd/server
	./_build/hgraber-bin -p 8081
	
view: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin  ./cmd/server
	./_build/hgraber-bin -v -p 8081

debug: create_build_dir
	go build $(LDFLAGS) -trimpath -o ./_build/hgraber-bin  ./cmd/server
	./_build/hgraber-bin -stdfile-append -debug -debug-fullpath -h 127.0.0.1 -p 8080 -static="internal/service/webServer/static" --access-token=local-debug

demo: create_build_dir
	go build $(LDFLAGS) -trimpath -o ./_build/hgraber-bin  ./cmd/inmemory
	./_build/hgraber-bin -debug -h 127.0.0.1 -p 8081 --access-token=local-debug

fileserver: create_build_dir
	go build $(LDFLAGS) -trimpath -o ./_build/hgraber-fileserver  ./cmd/fileserver
	./_build/hgraber-fileserver -addr 127.0.0.1:8082 -token fs-local
