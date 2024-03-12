.PHONY: toolbar cli

toolbar:
	$(MAKE) -C toolbar build

cli:
	$(MAKE) -C cmd/cli build

build: cli toolbar