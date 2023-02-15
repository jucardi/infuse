GO_FMT     = gofmt -s -w -l .
BUILD_TIME = $(shell date +%Y-%m-%dT%H:%M:%s)
CMDROOT    = github.com/jucardi/infuse/cmd/infuse
COMMIT_ID  = $(shell git rev-parse HEAD)
LOCAL_VER ?= git.commit-$(shell git rev-parse HEAD).local
VERSION    = 1.0.12

deps:
	@go get ./...

vet:
	@go vet ./...

format:
	$(GO_FMT)

test:
	@echo "running test coverage..."
	@mkdir -p test-artifacts/coverage
	@go test -mod=vendor ./... -v -coverprofile test-artifacts/cover.out
	@go tool cover -func test-artifacts/cover.out

compile-all: deps
	@echo "compiling..."
	@rm -rf build
	@mkdir build
	@echo "building x86_64 linux binary..."
	@GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=$(VERSION) -X $(CMDROOT)/version.Built=$(BUILD_TIME) -X $(CMDROOT)/version.Commit=$(COMMIT_ID) -X $(CMDROOT)/version.Arch=x86_64" -o build/infuse-Linux-x86_64 ./cmd/infuse
	@shasum -a 256 build/infuse-Linux-x86_64 >> build/infuse-Linux-x86_64.sha256
	@echo "building arm64 linux binary..."
	@GOOS=linux GOARCH=arm64 go build -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=$(VERSION) -X $(CMDROOT)/version.Built=$(BUILD_TIME) -X $(CMDROOT)/version.Commit=$(COMMIT_ID) -X $(CMDROOT)/version.Arch=arm64" -o build/infuse-Linux-arm64 ./cmd/infuse
	@shasum -a 256 build/infuse-Linux-arm64 >> build/infuse-Linux-arm64.sha256
	@echo "building x86_64 macosx binary..."
	@GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=$(VERSION) -X $(CMDROOT)/version.Built=$(BUILD_TIME) -X $(CMDROOT)/version.Commit=$(COMMIT_ID) -X $(CMDROOT)/version.Arch=x86_64" -o build/infuse-Darwin-x86_64 ./cmd/infuse
	@shasum -a 256 build/infuse-Darwin-x86_64 >> build/infuse-Darwin-x86_64.sha256
	@echo "building arm64 macosx binary..."
	@GOOS=darwin GOARCH=arm64 go build -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=$(VERSION) -X $(CMDROOT)/version.Built=$(BUILD_TIME) -X $(CMDROOT)/version.Commit=$(COMMIT_ID) -X $(CMDROOT)/version.Arch=arm64" -o build/infuse-Darwin-arm64 ./cmd/infuse
	@shasum -a 256 build/infuse-Darwin-arm64 >> build/infuse-Darwin-arm64.sha256
	@echo "building x86_64 windows binary..."
	@GOOS=windows GOARCH=amd64 go build -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=$(VERSION) -X $(CMDROOT)/version.Built=$(BUILD_TIME) -X $(CMDROOT)/version.Commit=$(COMMIT_ID) -X $(CMDROOT)/version.Arch=x86_64" -o build/infuse-Windows-x86_64.exe ./cmd/infuse
	@shasum -a 256 build/infuse-Windows-x86_64.exe >> build/infuse-Windows-x86_64.exe.sha256
	@echo "building arm64 windows binary..."
	@GOOS=windows GOARCH=arm64 go build -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=$(VERSION) -X $(CMDROOT)/version.Built=$(BUILD_TIME) -X $(CMDROOT)/version.Commit=$(COMMIT_ID) -X $(CMDROOT)/version.Arch=arm64" -o build/infuse-Windows-arm64.exe ./cmd/infuse
	@shasum -a 256 build/infuse-Windows-arm64.exe >> build/infuse-Windows-arm64.exe.sha256

install:
	@go install -mod=vendor -ldflags "-X $(CMDROOT)/version.Version=git.commit-$(COMMIT_ID).local -X $(CMDROOT)/version.Built=$(BUILD_TIME)" ./cmd/infuse
