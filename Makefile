.DEFAULT_GOAL := help

# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))
PKG := github.com/natemarks/pgsummary
VERSION := 0.0.3
COMMIT := $(shell git rev-parse HEAD)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
CDIR = $(shell pwd)
EXECUTABLES := pgreport pgcompare
GOOS := linux darwin
GOARCH := amd64

CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DEFAULT_BRANCH := main

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

clean-venv: ## re-create virtual env
	rm -rf .venv
	python3 -m venv .venv
	( \
       source .venv/bin/activate; \
       pip install --upgrade pip setuptools; \
    )

${EXECUTABLES}:
	@for o in $(GOOS); do \
	  for a in $(GOARCH); do \
        echo "$(COMMIT)/$${o}/$${a}" ; \
        mkdir -p build/$(COMMIT)/$${o}/$${a} ; \
        echo "VERSION: $(VERSION)" > build/$(COMMIT)/$${o}/$${a}/version.txt ; \
        echo "COMMIT: $(COMMIT)" >> build/$(COMMIT)/$${o}/$${a}/version.txt ; \
        env GOOS=$${o} GOARCH=$${a} \
        go build  -v -o build/$(COMMIT)/$${o}/$${a}/$@ ${PKG}/cmd/$@ ; \
	  done \
    done ; \

build: git-status ${EXECUTABLES}
	rm -f build/current
	ln -s $(CDIR)/build/$(COMMIT) $(CDIR)/build/current

release:
	mkdir -p release/$(VERSION)
	@for o in $(GOOS); do \
	  for a in $(GOARCH); do \
        tar -C ./build/$(COMMIT)/$${o}/$${a} -czvf release/$(VERSION)/pgsummary_$(VERSION)_$${o}_$${a}.tar.gz . ; \
	  done \
    done ; \

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

static: vet lint test

clean:
	-@rm ${OUT} ${OUT}-v*


bump: clean-venv  ## bump version in main branch
ifeq ($(CURRENT_BRANCH), $(DEFAULT_BRANCH))
	( \
	   source .venv/bin/activate; \
	   pip install bump2version; \
	   bump2version $(part); \
	)
else
	@echo "UNABLE TO BUMP - not on Main branch"
	$(info Current Branch: $(CURRENT_BRANCH), main: $(DEFAULT_BRANCH))
endif


git-status: ## require status is clean so we can use undo_edits to put things back
	@status=$$(git status --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "Error - working directory is dirty. Commit those changes!"; \
		exit 1; \
	fi


.PHONY: build release static upload vet lint