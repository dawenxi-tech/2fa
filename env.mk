
# os
OS_GO_BIN_NAME=go
ifeq (uname),Windows)
	OS_GO_BIN_NAME=go.exe
endif
OS_GO_OS=$(shell $(OS_GO_BIN_NAME) env GOOS)
OS_GO_VERSION=$(shell $(OS_GO_BIN_NAME) env GOVERSION)
OS_GO_ARCH=$(shell $(OS_GO_BIN_NAME) env GOARCH)

# git
# assumes git is installed

OS_GIT_BIN_NAME=git
ifeq ($(OS_GO_OS),windows)
	OS_GIT_BIN_NAME=git.exe
endif
OS_GIT_BIN_WHICH=$(shell which $(OS_GIT_BIN_NAME))
OS_GIT_BIN_WHICH_VERSION=$(shell $(OS_GIT_BIN_NAME) -v)

OS_GIT_SHA    = $(shell cd $(OS_GIT_FSPATH) && git rev-parse --short HEAD)
OS_GIT_TAG    = $(shell cd $(OS_GIT_FSPATH) && git describe --tags --abbrev=0 --exact-match 2>/dev/null)
OS_GIT_DIRTY  = $(shell cd $(OS_GIT_FSPATH) && test -n "`git status --porcelain`" && echo "dirty" || echo "clean")


# make
# assumes make is installed
OS_MAKE_BIN_NAME=make
ifeq ($(OS_GO_OS),windows)
	OS_MAKE_BIN_NAME=make.exe
endif
OS_MAKE_BIN_WHICH=$(shell which $(OS_MAKE_BIN_NAME))
OS_MAKE_BIN_WHICH_VERSION=$(MAKE_VERSION)

# Packaging
APP_NAME=2FA
BUNDLE_ID=tech.someone.2fa
DIR_RELEASE=./dist/release
APP_ICON=./assets-backup/something.png

# Github CI env variables: https://docs.github.com/en/actions/learn-github-actions/variables

env-print:

	@echo ""
	@echo "-- .env --"

	@echo "--- os ---"
	@echo "---- bin ----"
	@echo "OS_GO_BIN_NAME:             $(OS_GO_BIN_NAME)"
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
	@echo "---- cvar ----"
	@echo "OS_GIT_SHA:                  $(OS_GIT_SHA)"
	@echo "OS_GIT_TAG:                  $(OS_GIT_TAG)"
	@echo "OS_GIT_DIRTY:                $(OS_GIT_DIRTY)"

	@echo ""
	@echo "--- packaging ---"
	@echo "APP_NAME:                    $(APP_NAME)"
	@echo "BUNDLE_ID:                   $(BUNDLE_ID)"
	@echo "DIR_RELEASE:                 $(DIR_RELEASE)"
	@echo "APP_ICON:                    $(APP_ICON)"
	@echo ""

