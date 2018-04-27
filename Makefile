GO_FMT     = gofmt -s -w -l .
BUILD_TIME = $(shell date +%Y-%m-%dT%H:%M:%s)
PKGROOT    = github.com/jucardi/go-infuse
VERSION   ?= git.commit-$(shell git rev-parse HEAD).local

vet:
	@go vet ./...

check:
	$(GO_FMT)
	@go vet ./...

format:
	$(GO_FMT)

test-deps:
	@echo "installing test dependencies..."
	@go get github.com/stretchr/testify/assert
	@go get github.com/smartystreets/goconvey/convey
	@go get github.com/axw/gocov/...
	@go get github.com/AlekSi/gocov-xml
	@go get gopkg.in/matm/v1/gocov-html

test: test-deps
	@echo "running test coverage..."
	@mkdir -p test-artifacts/coverage
	@gocov test ./... -v > test-artifacts/gocov.json
	@cat test-artifacts/gocov.json | gocov report
	@cat test-artifacts/gocov.json | gocov-xml > test-artifacts/coverage/coverage.xml
	@cat test-artifacts/gocov.json | gocov-html > test-artifacts/coverage/coverage.html

compile-all:
	@echo "compiling..."
	@rm -rf build
	@mkdir build
	@echo "building linux binary..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/jucardi/go-infuse/cmd/infuse/version.Version=$(VERSION) -X github.com/jucardi/go-infuse/cmd/infuse/version.Built=$(BUILD_TIME)" -o build/infuse-linux-amd64 ./cmd/infuse
	@echo "building macosx binary..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/jucardi/go-infuse/cmd/infuse/version.Version=$(VERSION) -X github.com/jucardi/go-infuse/cmd/infuse/version.Built=$(BUILD_TIME)" -o build/infuse-darwin-amd64 ./cmd/infuse
	@echo "building windows binary..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/jucardi/go-infuse/cmd/infuse/version.Version=$(VERSION) -X github.com/jucardi/go-infuse/cmd/infuse/version.Built=$(BUILD_TIME)" -o build/infuse-windows-amd64.exe ./cmd/infuse