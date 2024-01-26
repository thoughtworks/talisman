# Contributing to Talisman

By contributing to Talisman, you agree to abide by the [code of conduct](CODE_OF_CONDUCT.md).

## How to start contributing

If you are not sure how to begin contributing to Talisman, have a look at the issues tagged under [good first issue](https://github.com/thoughtworks/talisman/labels/good%20first%20issue).

## Developing locally

To contribute to Talisman, you need a working golang development environment.
Check [this link](https://golang.org/doc/install) to help you get started.

Once you have go installed and set up, clone the talisman repository. In your
working copy, fetch the dependencies by having go mod fetch them for you:

```
go mod vendor
```

Run the tests:

```
go test ./...
```

Build talisman:

```
go build -o dist/talisman -ldflags="-s -w" talisman/cmd
```

To build for multiple platforms we use [GoReleaser](https://goreleaser.com/):

```
goreleaser build --snapshot --clean
```

## Submitting a Pull Request

To send in a pull request

1. Fork the repo.
2. Create a new feature branch based off the master branch.
3. Provide the commit message with the issue number and a proper description.
4. Ensure that all the tests pass.
5. Submit the pull request.

## Updating Talisman GitHub Pages

1. Checkout a new branch from gh-pages
2. Navigate to the docs/ folder and update the files
3. See instructions for checking locally [here](https://github.com/thoughtworks/talisman/blob/gh-pages/README.md).
4. Raise a pull request against the branch gh-pages

## Releasing

1. Tag the commit to be released with the next version according to
[semver](https://semver.org/) conventions
2. Push the tag to trigger the GitHub Actions Release pipeline
3. Approve the [drafted GitHub Release](https://github.com/thoughtworks/talisman/releases)

The latest version will now be accessible to anyone who builds their own
binaries, downloads binaries directly from GitHub Releases or homebrew, or uses
the install script from the website.
