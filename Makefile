
macos:
	mkdir -p dist/release/macos
	rm -rf dist/release/macos
	GOOS=darwin GOARCH=arm64 go build -tags='RELEASE' -o dist/release/macos/2fa github.com/dawenxi-tech/2fa
	go run cmd/macapp/macapp.go -assets ./dist/release/macos \
                                		-bin 2fa \
                                		-icon ./assets/2fa.png \
                                		-identifier tech.dawenxi.2fa \
                                		-name "2FA" \
                                		-dmg "2FA.dmg" \
                                		-o dist/release/macos-release