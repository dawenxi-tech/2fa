

APP_NAME=2FA
BUNDLE_ID=tech.dawenxi.2fa
DIR_RELEASE=./dist/release
APP_ICON=./assets-backup/2fa.png

# env
# override the above environment varaibles as needed
include env.mk


ci-all: env-print
	@echo ""
	@echo "ci-all called"
	@echo ""

	# Now build based on the CI env variable.
	@echo "RUNNER_OS:        $(RUNNER_OS)" 
	@echo "RUNNER_ARCH:      $(RUNNER_ARCH)" 

dep-sub:
	# to pull in gio
	git submodule update --init --recursive
	

dep-tools:
	# icns viewer
	go install github.com/jackmordaunt/icns/cmd/preview@latest

	# icns maker doing png to icns
	go install github.com/jackmordaunt/icns/v2/cmd/icnsify@latest

	# gio command for building cross platform
	go install gioui.org/cmd/gogio@latest

### ASSETS

assets-convert:
	# First we copy the PNG we want up to assets folder ( which the packaing uses as truth)
	cp $(APP_ICON) ./assets/2fa.png

	# Then, we do the conversion of the PNG to ICNS
	icnsify --input ./assets/2fa.png --output ./assets/2fa.icns
assets-preview:
	# Lets check if the conversion worked and check the diffeerent resolutions.
	preview $(PWD)/assets/2fa.icns



### BUILD 

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
	rm -rf ${DIR_RELEASE}/macos/$(MAC_ARCH)
	#TODO: release tag. cant see how to do it with gio command yet..
	gogio -target macos -arch $(MAC_ARCH) -appid $(BUNDLE_ID) -icon $(APP_ICON) -o ${DIR_RELEASE}/macos/app/$(MAC_ARCH)/$(APP_NAME).app . 

#### windows 

build-windows-all: build-windows-amd64 build-windows-arm64

build-windows-amd64:
	WINDOWS_ARCH=amd64 $(MAKE) build-windows

build-windows-arm64:
	WINDOWS_ARCH=arm64 $(MAKE) build-windows

build-windows:
	rm -rf ${DIR_RELEASE}/windows/$(WINDOWS_ARCH)
	gogio -target windows -arch $(WINDOWS_ARCH) -appid $(BUNDLE_ID) -icon $(APP_ICON) -o ${DIR_RELEASE}/windows/exe/$(WINDOWS_ARCH)/$(APP_NAME).exe .



### PACKAGE

# into DMG and MSI ( must be run on correct OS )

pack-all: pack-macos pack-windows

pack-macos:
	# arm64 for now. Refactor like the build is, when working.

	rm -rf ${DIR_RELEASE}/macos/dmg/arm64

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
          ${DIR_RELEASE}/macos/app/arm64
	mkdir -p ${DIR_RELEASE}/macos/dmg/arm64

	mv 2FA-Installer.dmg ${DIR_RELEASE}/macos/dmg/arm64			

pack-windows:
	# todo when windows packager worked out.



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
