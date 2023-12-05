

APP_NAME=2FA
BUNDLE_ID=tech.dawenxi.2fa
DIR_RELEASE=./dist/release
APP_ICON=./assets-backup/2fa.png

# env
# to override above variables.
include env.mk


all:
	# .env
	$(MAKE) env-print

	$(MAKE) dep-sub

	$(MAKE) build

ci-build: 
	# only run in ci on a branch push
	@echo ""
	@echo "ci-build called ..."
	@echo ""

	# .ci
	@echo ""
	@echo "-- .ci --"
	@echo "--- env ---"
	@echo "RUNNER_OS:        $(RUNNER_OS)"
	@echo "RUNNER_ARCH:      $(RUNNER_ARCH)"

	$(MAKE) all

ci-release: 
	# only run in ci on a tag push !
	@echo ""
	@echo "ci-release called ..."
	@echo ""

	$(MAKE) ci-build

	$(MAKE) pack

	$(MAKE) release


dep-sub:
	@echo ""
	@echo "Installing gio sub module ..."
	git submodule update --init --recursive

	# and upgrade it. 
	git submodule update --remote
	@echo ""

dep-tools:
	@echo ""
	@echo "Installing tools ..."

	# icns maker doing png to icns
	# icns viewer for checking they are ok.
	# https://github.com/JackMordaunt/icns/releases/tag/v2.2.7
	go install github.com/jackmordaunt/icns/v2/cmd/icnsify@v2.2.7
	# only works on latest...
	#go install github.com/jackmordaunt/icns/cmd/preview@v2.2.7
	go install github.com/jackmordaunt/icns/cmd/preview@latest

	# gio command for building cross platform
	# https://github.com/gioui/gio-cmd
	go install gioui.org/cmd/gogio@latest

	# simple file listing help
	# https://github.com/a8m/tree
	go install github.com/a8m/tree/cmd/tree@latest

	# easy way to migrate goalng code to latest dependencies
	# https://github.com/oligot/go-mod-upgrade/releases/tag/v0.9.1
	go install github.com/oligot/go-mod-upgrade@v0.9.1



	@echo ""
### MODS

mod-up:

	# for example: https://github.com/gioui/gio/releases/tag/v0.4.1
	go-mod-upgrade

### ASSETS

assets-convert:
	@echo ""
	@echo "Asset conversion ..."
	# First we copy the PNG we want up to assets folder ( which the build and packaging uses as truth )
	cp $(APP_ICON) ./assets/2fa.png

	# Then, we do the conversion of the PNG to ICNS
	icnsify --input ./assets/2fa.png --output ./assets/2fa.icns
	@echo ""

assets-preview:
	# Lets check if the conversion worked and check the different resolutions.
	preview $(PWD)/assets/2fa.icns



### BUILD

build:
	@echo ""
	@echo "Building phase ..."

ifeq ($(OS_GO_OS),windows)
	@echo ""
	@echo "Detected Windows ..."
	$(MAKE) dep-tools

	# Windows cant build tray code: https://github.com/gedw99/2fa/actions/runs/7034294593/job/19142004038
	@echo "Skipping Windows until we support Windows tray ..."
	#$(MAKE) build-windows-all
	@echo ""
endif

ifeq ($(OS_GO_OS),darwin)
	@echo ""
	@echo "Detected Darwin so building ..."
	$(MAKE) dep-tools
	$(MAKE) build-macos-all
	@echo ""
endif

ifeq ($(OS_GO_OS),linux)
	@echo ""
	@echo "Detected Linux but we have no Linux support yet, so skipping ..."
	@echo ""
endif 


build-list:
	@echo ""
	@echo "Build produced ..."
	tree  -l -a -C $(DIR_RELEASE)
	@echo ""

build-all: build-macos-all build-windows-all 

build-all-del:
	rm -rf $(DIR_RELEASE)

#### macos

build-macos-all: build-macos-amd64 build-macos-arm64

build-macos-amd64:
	MAC_ARCH=amd64 $(MAKE) build-macos

build-macos-arm64:
	MAC_ARCH=arm64 $(MAKE) build-macos

build-macos:
	@echo ""
	@echo "Building Darwin $(MAC_ARCH) ..."

	rm -rf ${DIR_RELEASE}/macos/$(MAC_ARCH)
	#TODO: release tag. cant see how to do it with gio command yet..
	gogio -target macos -arch $(MAC_ARCH) -appid $(BUNDLE_ID) -icon $(APP_ICON) -o ${DIR_RELEASE}/macos/app/$(MAC_ARCH)/$(APP_NAME).app . 

	$(MAKE) build-list

#### windows 

build-windows-all: build-windows-amd64 build-windows-arm64

build-windows-amd64:
	WINDOWS_ARCH=amd64 $(MAKE) build-windows

build-windows-arm64:
	WINDOWS_ARCH=arm64 $(MAKE) build-windows

build-windows:
	@echo ""
	@echo "Building Windows $(WINDOWS_ARCH) ..."
	rm -rf ${DIR_RELEASE}/windows/$(WINDOWS_ARCH)
	gogio -target windows -arch $(WINDOWS_ARCH) -appid $(BUNDLE_ID) -icon $(APP_ICON) -o ${DIR_RELEASE}/windows/exe/$(WINDOWS_ARCH)/$(APP_NAME).exe .

	$(MAKE) build-list

### RUN

run:
	# With gio its best to run off a .app or .exe, rather an using ``` go run . ```, 
	# so that your are seeing all icons and other things work.
	@echo ""
	@echo "Running. Assume you done a build already .."

ifeq ($(OS_GO_OS),windows)
	@echo ""
	@echo "Detected Windows ..."
	
	@echo ""
endif

ifeq ($(OS_GO_OS),darwin)
	@echo ""
	@echo "Detected Darwin ..."

	@echo ""
endif

ifeq ($(OS_GO_OS),linux)
	@echo ""
	@echo "Detected Linux ..."
	@echo ""
endif

run-macos-amd64:
	MAC_ARCH=amd64 $(MAKE) run-macos

run-macos-arm64:
	MAC_ARCH=arm64 $(MAKE) run-macos

run-macos:
	open ${DIR_RELEASE}/macos/app/$(MAC_ARCH)/$(APP_NAME).app


### PACKAGE

# into DMG and MSI ( must be run on correct OS )


pack:
	@echo ""
	@echo "Packaging phase ..."

ifeq ($(OS_GO_OS),windows)
	@echo ""
	@echo "Detected Windows ... Add later when ready."
	@echo ""
endif

ifeq ($(OS_GO_OS),darwin)
	@echo ""
	@echo "Detected Darwin ..."
	$(MAKE) pack-macos-all
	@echo ""
endif

ifeq ($(OS_GO_OS),linux)
	@echo ""
	@echo "Detected Linux ... Add later when ready."
	@echo ""
endif
	$(MAKE) build-list
	@echo ""
	@echo "Packaging phase done ..."
	@echo ""

pack-all: pack-macos-all pack-windows-all

pack-macos-all: pack-macos-amd64 pack-macos-arm64
	
pack-macos-amd64:
	MAC_ARCH=amd64 $(MAKE) pack-macos

pack-macos-arm64:
	MAC_ARCH=arm64 $(MAKE) pack-macos

pack-macos:
	# need this
	brew install create-dmg

	rm -rf ${DIR_RELEASE}/macos/dmg/$(MAC_ARCH)

	create-dmg \
          --volname "2FA Installer" \
          --volicon "./assets/2fa.icns" \
          --window-pos 200 120 \
          --window-size 800 400 \
          --icon-size 100 \
          --icon "2FA.app" 200 190 \
          --hide-extension "2FA.app" \
          --app-drop-link 600 185 \
          "2FA-Installer.dmg" \
          ${DIR_RELEASE}/macos/app/$(MAC_ARCH)
	mkdir -p ${DIR_RELEASE}/macos/dmg/$(MAC_ARCH)

	mv 2FA-Installer.dmg ${DIR_RELEASE}/macos/dmg/$(MAC_ARCH)			

pack-windows-all:
	# todo when windows packager worked out.



### RELEASE
# need https://github.com/cli/cli which is called "gh". Yeah well done with naming :)


release:
	@echo ""
	@echo "Release phase ..."

ifeq ($(OS_GO_OS),windows)
	@echo ""
	@echo "Detected Windows ..."
	$(MAKE) release-windows
	@echo ""
endif

ifeq ($(OS_GO_OS),darwin)
	@echo ""
	@echo "Detected Darwin ..."
	$(MAKE) release-macos
	@echo ""
endif

ifeq ($(OS_GO_OS),linux)
	@echo ""
	@echo "Detected Linux ..."
	$(MAKE) release-linux
	@echo ""
endif
	$(MAKE) build-list
	@echo ""
	@echo "Release phase done ..."
	@echo ""

release-macos:
	brew install gh
release-windows:
	# no idea what to do here for windows ...
	#scoop install gh
release-linux:
	# no idea what to do here for linux...
	#brew install gh



### OLD below. Kept for refernce until above works...

macos: macos-app macos-dmg

macos-app:
	mkdir -p ${DIR_RELEASE}/macos
	rm -rf ${DIR_RELEASE}/macos/*
	go build -tags='RELEASE' -ldflags="-s -w" -o ${DIR_RELEASE}/macos/binary/2fa github.com/dawenxi-tech/2fa
	go run cmd/macapp/macapp.go -assets ${DIR_RELEASE}/macos/binary \
                                		-bin 2fa \
                                		-icon ./assets/2fa.png \
                                		-identifier ${BUNDLE_ID} \
                                		-name ${APP_NAME} \
                                		-o dist/release/macos/app

macos-dmg:
	rm -rf ${DIR_RELEASE}/macos/dmg
	create-dmg \
          --volname "2FA Installer" \
          --volicon "./assets/2fa.icns" \
          --window-pos 200 120 \
          --window-size 800 400 \
          --icon-size 100 \
          --icon "2FA.app" 200 190 \
          --hide-extension "2FA.app" \
          --app-drop-link 600 185 \
          "2FA-Installer.dmg" \
          ${DIR_RELEASE}/macos/app
	mkdir ${DIR_RELEASE}/macos/dmg
	mv 2FA-Installer.dmg ${DIR_RELEASE}/macos/dmg/
