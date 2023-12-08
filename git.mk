# git operations for devs on forks.

git-fork-init:
	# do this after git clone your fork.
	# config the upstream remote.
	$(OS_GIT_BIN_NAME) remote add upstream https://github.com/dawenxi-tech/2fa
	$(OS_GIT_BIN_NAME) remote -v

	# Get the sub modules...
	$(MAKE) dep-sub

git-fork-merge-upstream:
	# do this to catchup to upstream main repo.
	# # best to do this before you do a local push up to your own Github Repo, so that all changes from others are included.

	# pull from the upstream github repo and merge to this local fork.
	# https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork

	# https://everythingdevops.dev/how-to-avoid-merge-commits-when-syncing-a-fork/

	$(OS_GIT_BIN_NAME) remote -v
	$(OS_GIT_BIN_NAME) pull --rebase upstream main

	#$(OS_GIT_BIN_NAME) fetch upstream
	#$(OS_GIT_BIN_NAME) checkout main
	#$(OS_GIT_BIN_NAME) merge upstream/main

GIT_COMMIT_MESSAGE='chore'

git-fork-commit-push:
	# do this when you ready to push changes to your github repo and then make a PR.

	# example: make $GIT_COMMIT_MESSAGE='chore-test' git-fork-commit-push

	$(OS_GIT_BIN_NAME) add --all
	$(OS_GIT_BIN_NAME) commit -am $(GIT_COMMIT_MESSAGE)
	$(OS_GIT_BIN_NAME) push origin main --force
