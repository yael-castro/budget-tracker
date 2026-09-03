.PHONY: summary build

build:
	bash ./scripts/build.bash

summary: build
	bash ./scripts/summary.bash

