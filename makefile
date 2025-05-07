.PHONY: toolbar cli copy-cli

all: cli copy-cli toolbar

MAKE := make

toolbar:
	$(MAKE) -C toolbar build

cli:
	$(MAKE) -C cmd/cli bump-patch
	$(MAKE) -C cmd/cli build

copy-cli:
	cp cmd/cli/fyodorov ./fyodorov
	cp cmd/cli/fyodorov ~/.fyodorov/fyodorov

build: cli toolbar