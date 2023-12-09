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
# Does everything to buld and package the app.
make all

# To setup your fork. 
# Run this after you have done a git clone of your remote fork to your local laptop.
make git-fork-init

# To fetch and rebase upstream to your local fork. 
# Run this before you push, so that you have everyones else changes rebased onto your repoö
make git-fork-merge-upstream

# To commit and force push your local fork to your remote github fork.
# Run this when your want test your changes in CI, and then PR ( usinfg the wbe gui ) to the remote Upstream repo.
make GIT_COMMIT_MESSAGE='git-test' git-fork-commit-push

# Get the gio sub module.
make dep-sub 

# Install all tooling need for your OS
make dep-tools

# Bulld for your OS.
make buuld

# Package for your OS.
make pack

# Called by Github action to build for all Os's.
# Run this to do locally what will happen in Github CI.
make ci-build
make ci-release
```

## Releases

There is no automated releases yet.

## Updating

There is no automatic updating yet.


