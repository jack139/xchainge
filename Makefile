BUILD=build
BUILD_TAGS=cleveldb
CGO_ENABLED=1

all: build

build:
	CGO_ENABLED=$(CGO_ENABLED) go build -o $(BUILD)/xchain -tags '$(BUILD_TAGS)' xchain/main.go

clean:
	@rm -rf $(BUILD)
	@go clean
