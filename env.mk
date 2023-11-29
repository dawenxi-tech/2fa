
# OS
OS_GO_BIN_NAME=go
ifeq (uname),Windows)
	OS_GO_BIN_NAME=go.exe
endif
OS_GO_OS=$(shell $(OS_GO_BIN_NAME) env GOOS)
OS_GO_VERSION=$(shell $(OS_GO_BIN_NAME) env GOVERSION)
OS_GO_ARCH=$(shell $(OS_GO_BIN_NAME) env GOARCH)

# Packaging
APP_NAME=2FA
BUNDLE_ID=tech.someone.2fa
DIR_RELEASE=./dist/release
APP_ICON=./assets-backup/something.png

env-print:

	@echo ""
	@echo "-- .env --"
	@echo "--- OS ---"
	@echo "OS_GO_BIN_NAME:   $(OS_GO_BIN_NAME)"
	@echo "OS_GO_OS:         $(OS_GO_OS)"
	@echo "OS_GO_VERSION:    $(OS_GO_VERSION)"
	@echo "OS_GO_ARCH:       $(OS_GO_ARCH)"
	@echo ""

	@echo ""
	@echo "--- packaging ---"
	@echo "APP_NAME:      $(APP_NAME)"
	@echo "BUNDLE_ID:     $(BUNDLE_ID)"
	@echo "DIR_RELEASE:   $(DIR_RELEASE)"
	@echo "APP_ICON:      $(APP_ICON)"


