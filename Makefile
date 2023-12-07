

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

	$(MAKE) dep-tools

	$(MAKE) build

	$(MAKE) pack

ci-build: 
	# only runs in ci on a branch push
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
	# only runs in ci on a tag push !
	@echo ""
	@echo "ci-release called ..."
	@echo ""

	$(MAKE) ci-build

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
	@echo "Dep tools phase ..."

ifeq ($(OS_GO_OS),windows)
	@echo ""
	@echo "Detected Windows ..."
	# see: https://gioui.org/doc/install/windows
	@echo "None needed for Windows."
	@echo ""
endif

ifeq ($(OS_GO_OS),darwin)
	@echo ""
	@echo "Detected Darwin ..."
	# see: https://gioui.org/doc/install/macos

	@echo ""
	@echo "Installing create-dmg with brew, so that we can do packaging on Mac."
	@echo ""
	brew install create-dmg
	@echo ""
endif

ifeq ($(OS_GO_OS),linux)
	@echo ""
	@echo "Detected Linux ..."

	# from: https://github.com/g45t345rt/g45w/actions/runs/6285743746/job/17068447935

	@echo ""
	@echo "Installing various things, so that gio can build and run on Linux."
	@echo ""
	sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev libfuse2

	@echo ""
	@echo "Installing pkg2appimage with wget, so that we can do packaging with appimage on Linux."
	@echo ""
	wget -c $(shell wget -q https://api.github.com/repos/AppImageCommunity/pkg2appimage/releases -O - | grep "pkg2appimage-.*-x86_64.AppImage" | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
	chmod +x ./pkg2appimage-*.AppImage
	mkdir -p dist
	mv ./pkg2appimage-*.AppImage ./dist/

	@echo ""
endif
	
	@echo ""
	@echo "Installing icns maker and icns viewer using go install, so that we can generate icons."
	@echo ""
	# https://github.com/JackMordaunt/icns/releases/tag/v2.2.7
	go install github.com/jackmordaunt/icns/v2/cmd/icnsify@v2.2.7
	# only works on latest...
	#go install github.com/jackmordaunt/icns/cmd/preview@v2.2.7
	go install github.com/jackmordaunt/icns/cmd/preview@latest

	@echo ""
	@echo "Installing gio builder using go install, so that we can build gio apps."
	@echo ""
	# https://github.com/gioui/gio-cmd
	go install gioui.org/cmd/gogio@latest

	@echo ""
	@echo "Installing cross OS tree file printer using go install, so that we can see what we are making easily."
	@echo ""
	# https://github.com/a8m/tree
	go install github.com/a8m/tree/cmd/tree@latest

	@echo ""
	@echo "Installing cross OS go modules updater using go install, so that we can do golang module updates easily."
	@echo ""
	# https://github.com/oligot/go-mod-upgrade/releases/tag/v0.9.1
	go install github.com/oligot/go-mod-upgrade@v0.9.1

	@echo ""
	@echo "Dep tools phase done ..."
	@echo ""


### DIST

dist-list:
	@echo ""
	@echo "Dist folder has ..."
	tree  -l -a -C $(DIR_RELEASE)
	@echo ""

### MODS

mod-up:
	# for example: https://github.com/gioui/gio/releases/tag/v0.4.1
	go-mod-upgrade
mod-up-force:
	# so that in CI we can force an upgrade as part of the buuld.
	go-mod-upgrade -f
	
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

# TODO add versioning like go build -ldflags="-X 'github.com/dawenxi-tech/2fa/commands.Version=$(git describe --tags)'" -o "2fa-$(git describe --tags)-windows-amd64.exe"
# work out how to do with gio !

ifeq ($(OS_GO_OS),windows)
	@echo ""
	@echo "Detected Windows ..."
	# Windows cant build tray code: https://github.com/gedw99/2fa/actions/runs/7034294593/job/19142004038
	@echo "Skipping Windows until we support Windows tray ..."
	#$(MAKE) build-windows-all
	@echo ""
endif

ifeq ($(OS_GO_OS),darwin)
	@echo ""
	@echo "Detected Darwin so building ..."
	$(MAKE) build-macos-all
	@echo ""
endif

ifeq ($(OS_GO_OS),linux)
	@echo ""
	@echo "Detected Linux but"
	$(MAKE) build-linux-all
	@echo ""
endif
	@echo ""
	$(MAKE) dist-list
	@echo "Building phase done ..."


build-all: build-macos-all build-windows-all build-linux-all

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

	# NOTE: Not sure if we want the version to be part of the output name. I think not.
	gogio -target macos -arch $(MAC_ARCH) -appid $(BUNDLE_ID) -tags='RELEASE' -ldflags $(APP_VERSION_LDFLAGS) -icon $(APP_ICON) -o ${DIR_RELEASE}/macos/app/$(MAC_ARCH)/$(APP_NAME).app . 

	$(MAKE) dist-list

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
	gogio -target windows -arch $(WINDOWS_ARCH) -appid $(BUNDLE_ID) -ldflags $(APP_VERSION_LDFLAGS) -icon $(APP_ICON) -o ${DIR_RELEASE}/windows/exe/$(WINDOWS_ARCH)/$(APP_NAME).exe .

	$(MAKE) dist-list


### Linux
build-linux-all: build-linux-amd64 build-linux-arm64

build-linux-amd64:
	LINUX_ARCH=amd64 $(MAKE) build-linux

build-linux-arm64:
	LINUX_ARCH=arm64 $(MAKE) build-linux

build-linux:
	@echo ""
	@echo "Building Linux $(LINUX_ARCH) ..."

	rm -rf ${DIR_RELEASE}/linux/$(LINUX_ARCH)
	mkdir -p ${DIR_RELEASE}/linux/$(LINUX_ARCH)

	# gogio seems to not support linux desktop builds, so building with go.
	go build -tags='RELEASE' -ldflags $(APP_VERSION_LDFLAGS) -o ${DIR_RELEASE}/linux/$(LINUX_ARCH)/$(APP_NAME).exe . 

	cp 2fa-appimage.yml ${DIR_RELEASE}/linux/$(LINUX_ARCH)
	cp assets/2fa.png ${DIR_RELEASE}/linux/$(LINUX_ARCH)
	cd ${DIR_RELEASE}/linux/$(LINUX_ARCH); ../../../pkg2appimage-*.AppImage 2fa-appimage.yml

	$(MAKE) dist-list

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
	$(MAKE) dist-list
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

# to github tagged releases
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
	$(MAKE) dist-list
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

release-using-tag:
	# for dev to locally make a git tag and git tag, so that CI does a release.
	
