.PHONY: build_for_dev

build_for_dev:
	mkdir -p dist
	go build -o dist/cloud-init-decoder ./
