BUILD := $(CURDIR)/build

build:
	mkdir -p $(BUILD)
	cd pkg/consumer/main && go build -o $(BUILD)/cmdhandler-consumer
	cd pkg/producer/main && go build -o $(BUILD)/cmdhandler-producer
tests:
	go test ./... -v
clean:
	rm -rf $(BUILD)

.PHONY: build tests clean
