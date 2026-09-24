.PHONY: summary build

build:
	bash ./scripts/build.bash

summary: build
	./build/budget-tracker

