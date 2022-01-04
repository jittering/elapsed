
.DEFAULT_GOAL := build
SHELL := bash

clean:
	rm -rf build/ dist/

build:
	set -eo pipefail
	go build -o ./build/elapsed .
	echo "built build/elapsed"

release:
	goreleaser release --rm-dist --skip-validate

check-style:
	goreleaser --snapshot --rm-dist && brew style ./dist/*.rb


.PHONY: clean build
