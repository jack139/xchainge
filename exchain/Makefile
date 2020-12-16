OUTPUT=exchain

BUILD_TAGS=cleveldb
LD_FLAGS=
BUILD_FLAGS=-ldflags "$(LD_FLAGS)"

CGO_ENABLED=1


all: build

build:
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o $(OUTPUT)

clean:
	@rm -f $(OUTPUT)
	@go clean
