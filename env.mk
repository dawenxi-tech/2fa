
# env.mk
# standardisation of environment for corss platform builds.

# golang - assumed to be installed.

OS_GO_BIN_NAME=go
ifeq ($(shell uname),Windows)
	OS_GO_BIN_NAME=go.exe
endif
OS_GO_OS=$(shell $(OS_GO_BIN_NAME) env GOOS)
OS_GO_VERSION=$(shell $(OS_GO_BIN_NAME) env GOVERSION)
OS_GO_ARCH=$(shell $(OS_GO_BIN_NAME) env GOARCH)


# git - assumed to be installed.

OS_GIT_BIN_NAME=git
ifeq ($(OS_GO_OS),windows)
	OS_GIT_BIN_NAME=git.exe
endif
OS_GIT_BIN_WHICH=$(shell which $(OS_GIT_BIN_NAME))
OS_GIT_BIN_WHICH_VERSION=$(shell $(OS_GIT_BIN_NAME) -v)

OS_GIT_FSPATH = $(shell pwd)
OS_GIT_SHA    = $(shell cd $(OS_GIT_FSPATH) && git rev-parse --short HEAD)
OS_GIT_TAG    = $(shell cd $(OS_GIT_FSPATH) && git describe --tags --abbrev=0 --exact-match 2>/dev/null)
OS_GIT_DIRTY  = $(shell cd $(OS_GIT_FSPATH) && test -n "`git status --porcelain`" && echo "dirty" || echo "clean")


# make - assumed to be installed.
OS_MAKE_BIN_NAME=make
ifeq ($(OS_GO_OS),windows)
	OS_MAKE_BIN_NAME=make.exe
endif
OS_MAKE_BIN_WHICH=$(shell which $(OS_MAKE_BIN_NAME))
OS_MAKE_BIN_WHICH_VERSION=$(MAKE_VERSION)


# dist
DIR_RELEASE=./dist/release

# version
APP_VERSION=$(shell git-describe-semver -dir $(OS_GIT_FSPATH) --fallback v0.0.0)
APP_VERSION_LDFLAGS="-X 'github.com/dawenxi-tech/2fa/main.Version=$(APP_VERSION)'"

# Env - create a .env to overide these values.
APP_NAME=2FA
BUNDLE_ID=tech.someone.2fa
APP_ICON=./assets-backup/something.png

# Github CI env variables: https://docs.github.com/en/actions/learn-github-actions/variables

env-init:
	@echo ""
	@echo "Installing cross OS go semver tool using go install, so that we can get verion info easily."
	@echo ""
	# https://github.com/choffmeister/git-describe-semver/releases/tag/v0.3.11
	go install github.com/choffmeister/git-describe-semver@v0.3.11

env-print: env-init

	@echo ""
	@echo "-- .env --"

	@echo "--- os ---"
	@echo "---- bin ----"
	@echo "OS_GO_BIN_NAME:              $(OS_GO_BIN_NAME)"
	@echo "---- var ----"
	@echo "OS_GO_OS:                    $(OS_GO_OS)"
	@echo "OS_GO_VERSION:               $(OS_GO_VERSION)"
	@echo "OS_GO_ARCH:                  $(OS_GO_ARCH)"
	
	@echo ""
	@echo "--- git ---"
	@echo "---- bin ----"
	@echo "OS_GIT_BIN_NAME:             $(OS_GIT_BIN_NAME)"
	@echo "OS_GIT_BIN_WHICH:            $(OS_GIT_BIN_WHICH)"
	@echo "OS_GIT_BIN_WHICH_VERSION:    $(OS_GIT_BIN_WHICH_VERSION)"
	@echo "---- var ----"
	@echo "OS_GIT_FSPATH:               $(OS_GIT_FSPATH)"
	@echo "OS_GIT_SHA:                  $(OS_GIT_SHA)"
	@echo "OS_GIT_TAG:                  $(OS_GIT_TAG)"
	@echo "OS_GIT_DIRTY:                $(OS_GIT_DIRTY)"
	@echo ""
	@echo ""
	@echo "--- release ---"
	@echo "DIR_RELEASE:                 $(DIR_RELEASE)"
	@echo ""
	@echo ""
	@echo "--- version ---"
	@echo "APP_VERSION:                 $(APP_VERSION)"
	@echo "APP_VERSION_LDFLAGS:         $(APP_VERSION_LDFLAGS)"
	@echo ""
	@echo ""
	@echo "--- packaging ( override in .env ) ---"
	@echo "APP_NAME:                    $(APP_NAME)"
	@echo "BUNDLE_ID:                   $(BUNDLE_ID)"
	@echo "APP_ICON:                    $(APP_ICON)"
	@echo ""

