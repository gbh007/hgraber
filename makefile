TAG = $(shell git tag -l --points-at HEAD)
COMMIT = $(shell git show -s --abbrev=12 --pretty=format:%h HEAD)
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

LDFLAGS = -ldflags "-X 'app/system.Version=$(TAG)' -X 'app/system.Commit=$(COMMIT)' -X 'app/system.BuildAt=$(BUILD_TIME)'"

create_build_dir:
	mkdir -p ./_build

build: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./_build/hgraber-linux-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-linux-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-windows-amd64.exe
	tar -C ./_build -cf ./_build/hgraber.tar hgraber-linux-arm64 hgraber-linux-amd64 hgraber-windows-amd64.exe
	gzip -9f ./_build/hgraber.tar

build_arm64: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./_build/hgraber-arm64

build_amd64: create_build_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./_build/hgraber-amd64.exe

run: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin 
	./_build/hgraber-bin -p 8081
	
view: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin 
	./_build/hgraber-bin -v -p 8081

debug: create_build_dir
	go build $(LDFLAGS) -o ./_build/hgraber-bin 
	./_build/hgraber-bin -stdfile-append -debug -debug-fullpath -h 127.0.0.1 -p 8081 -static="service/webServer/static" --access-token=local-debug
	