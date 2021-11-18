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
	find . -type l -name $@ -exec rm -f {} \;
	mkdir -p build/$(COMMIT)/linux/amd64 build/linux/amd64
	env GOOS=linux GOARCH=amd64 \
	go build  -v -o build/$(COMMIT)/linux/amd64/$@ ${PKG}/cmd/$@
	ln -s $(CDIR)/build/$(COMMIT)/linux/amd64/$@ $(CDIR)/build/linux/amd64/$@
	mkdir -p build/$(COMMIT)/darwin/amd64 build/darwin/amd64
	env GOOS=darwin GOARCH=amd64 \
	go build  -v -o build/$(COMMIT)/darwin/amd64/$@ ${PKG}/cmd/$@
	ln -s $(CDIR)/build/$(COMMIT)/darwin/amd64/$@ $(CDIR)/build/darwin/amd64/$@
	echo $@

build: git-status ${EXECUTABLES}

release:  git-status build ## Build release versions
	echo "VERSION: $(VERSION)" > ./build/$(COMMIT)/version.txt
	echo "COMMIT: $(COMMIT)" >> ./build/$(COMMIT)/version.txt
	tar -C ./build/$(COMMIT) -czvf pgsummary-$(VERSION).tar.gz .

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