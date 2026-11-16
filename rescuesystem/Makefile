.DEFAULT_GOAL := image

# Resolve paths from this file, including when invoked with make -f.
RESCUESYSTEMDIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
IMAGE := $(RESCUESYSTEMDIR)/dist/rescuesystem.iso

.PHONY: image checksum clean up

# Deliberately avoid Git-backed flake capture of unrelated repository files.
image:
	@set -eu; \
	work=$$(mktemp -d "$${TMPDIR:-/tmp}/rescuesystem-build.XXXXXX"); \
	trap 'rm -rf -- "$$work"' EXIT; \
	trap 'exit 1' HUP INT TERM; \
	mkdir -p "$$work/source/rescuesystem"; \
	for file in flake.nix flake.lock; do \
		cp -- "$(RESCUESYSTEMDIR)/$$file" "$$work/source/$$file"; \
	done; \
	for file in banner.txt live-system.nix iso.nix; do \
		cp -- "$(RESCUESYSTEMDIR)/$$file" "$$work/source/rescuesystem/$$file"; \
	done; \
	nix build "path:$$work/source#rescuesystem" --out-link "$$work/result"; \
	test -s "$$work/result/iso/rescuesystem.iso"; \
	mkdir -p "$(dir $(IMAGE))"; \
	cp --remove-destination -- "$$work/result/iso/rescuesystem.iso" "$(IMAGE)"; \
	printf 'Image: %s\n' "$(IMAGE)"

up:
	nix flake update --flake "path:$(RESCUESYSTEMDIR)"

checksum: image
	@sha256sum "$(IMAGE)"

clean:
	rm -f -- "$(IMAGE)"
