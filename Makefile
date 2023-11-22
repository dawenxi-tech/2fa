
APP_NAME="2FA"
BUNDLE_ID="tech.dawenxi.2fa"
DIR_RELEASE="./dist/release"

macos: macos-app macos-dmg

macos-app:
	mkdir -p ${DIR_RELEASE}/macos
	rm -rf ${DIR_RELEASE}/macos/*
	go build -tags='RELEASE' -o ${DIR_RELEASE}/macos/binary/2fa github.com/dawenxi-tech/2fa
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
