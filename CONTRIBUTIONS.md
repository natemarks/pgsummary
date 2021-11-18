Feel free to add issues to the github project. I can accept PRs, but please open an issue and name the branch accordingly.

My release workflow is:
```shell
#commit changes to main
# build the executables in a directory tree like build/[COMMIT ID]
make build
# build/test/build/test...
# when I'm happy bump the version according to the type of change
make part=[major|minor|patch] bump
# now create a release tarball for the new version
make release
```
