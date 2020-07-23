DATE := $(shell date --iso-8601=seconds)

build/slab_exporter: test | vendor
	@go build -o ./build/slab_exporter \
		-ldflags "-X main.BuildDate=${DATE}" \
		./slab_exporter

go.mod:
	@GO111MODULE=on go mod tidy

go.sum: | go.mod
	@GO111MODULE=on go mod verify

vendor: | go.sum
	@GO111MODULE=on go mod vendor 

.PHONEY: install
install: | vendor
	@go install -v -ldflags "-X main.BuildDate=${DATE} ./slab_exporter

test: | vendor
	@go test -v -race -cover ./...

.PHONEY: clean
clean:
	rm -rf build vendor
