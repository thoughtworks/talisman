# Contributing to Talisman

By contributing to Talisman, you agree to abide by the [code of conduct](CODE_OF_CONDUCT.md).

## How to start contributing

If you are not sure how to begin contributing to Talisman, have a look at the issues tagged under [good first issue](https://github.com/thoughtworks/talisman/labels/good%20first%20issue).

## Developing locally

To contribute to Talisman, you need a working golang development
environment. Check [this link](https://golang.org/doc/install) to help
you get started with that.

Talisman now uses go modules (GO111MODULE=on) to manage dependencies

Once you have go 1.11 installed and setup, clone the talisman repository. In your
working copy, fetch the dependencies by having go mod fetch them for
you.

```
GO111MODULE=on go mod vendor
```

To run tests `GO111MODULE=on go test -mod=vendor ./...`

To build Talisman, we can use [gox](https://github.com/mitchellh/gox):

```
gox -osarch="darwin/amd64 linux/386 linux/amd64"
```

Convenience scripts `./build` and `./clean` perform build and clean-up as mentioned above.

## Submitting a Pull Request

To send in a pull request

1. Fork the repo.
2. Create a new feature branch based off the master branch.
3. Provide the commit message with the the issue number and a proper description.
4. Ensure that all the tests pass.
5. Submit the pull request.

## Releasing

* Follow the instructions at the end of 'Developing locally' to build the binaries
* Bump the [version in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L10) according to [semver](https://semver.org/) conventions
* Update the [expected hashes in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L16-L18) to match the new binaries you just created (`shasum -b -a256 ...`)
* Make release commit and tag with the new version prefixed by `v` (like `git tag v0.3.0`)
* Push your release commit and tag: `git push && git push --tags`
* [Create a new release in github](https://github.com/thoughtworks/talisman/releases/new), filling in the new commit tag you just created
* Update the install script hosted on github pages: `git checkout gh-pages`, `git checkout master -- install.sh`, `git commit -m ...`

The latest version will now be accessible to anyone who builds their own binaries, downloads binaries directly from github releases, or uses the install script from the website.
