# 2fa

A 2FA desktop application.

There is a systems tray and a main window.

The system tray is a way to keep the App running in the backgorund. Only Mac for now with Winodws and Linux soon. Jump in if oyu wanrt to help or test.

## Usage

Currently works on Mac, with Windows to come.

## Development

All OS's need golang installed, we are at go version 1.21.5.

Also some tools that we install require:

- Mac needs brew installed, so that any brew based tools can be installed.
- Windows needs ? installed.
- Linux needed ? installed.

## Build

The Makefile is used for local and Github workflow builds so that we have a single source of truth.

All builds are automatically versioned using semver
based on the git hash or tag tag or a combination of gthe two.

Currently the following is built:

- Mac amd64 and arm64 ( DMG with app inside)
- Windows amd64 and arm64 ( only exe )
- Linus amd64 and arm64 ( not sure...)


```sh
# Does everything.
make all
```

```sh
# Get the gio sub module.
make dep-sub 
```

```sh
# install and tooling need for your OS
make dep-tools
```

```sh
# build for your OS
make buuld
```


```sh
# package for your os.
make pack
```


```sh
# Called by Github action only and builds for all OS using a simple github workflow Matrix.
make ci-build
make ci-release
```

## Releases

There is no automated releases yet.

## Updating

There is no automatic updating yet.


