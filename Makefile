
.DEFAULT_GOAL := build
SHELL := bash

clean:
	rm -rf build/ dist/

build:
	set -eo pipefail
	go build -o ./build/elapsed .
	echo "built build/elapsed"

release:
	goreleaser release --rm-dist --parallelism 1 --skip-validate


.PHONY: clean build
