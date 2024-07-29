# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=multiAdaptive-cli
BUILD_DIR=build

# All targets
.PHONY: all clean build

all: build

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): main.go
	@echo "Building Go application..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Copying srs files..."
	@mkdir -p $(BUILD_DIR)

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
