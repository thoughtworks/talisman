# Talisman

Talisman is a tool to validate code changes that are to be pushed out
of a local Git repository on a developer's workstation. By hooking
into the pre-push hook provided by Git, it validates the outgoing
changeset for things that look suspicious - such as potential SSH
keys, authorization tokens, private keys etc.

The aim is for this tool to do this through a variety of means
including file names and file content. We hope to have it be an
effective check to prevent potentially harmful security mistakes from
happening due to secrets which get accidentally checked in to a
repository.

The implementation as it stands is very bare bones and only has the
skeleton structure required to add the full range of functionality we
wish to incorporate. However, we encourage folks that want to
contribute to have a look around and contribute ideas/suggestions or
ideally, code that implements your ideas and suggestions!

## Installation

Talisman can either be installed into a single git repository, or as a global
[git hook template](https://git-scm.com/docs/git-init#_template_directory).

Talisman can be set up a as a pre-push or pre-commit hook on git repositories.


### Installation as a global hook template (recommended)
We recommend installing it as a git hook template, as that will cause
Talisman to be present, not only in your existing git repositories, but also in any new repository that you 'init' or
'clone'.

Use the [Global scripts Readme](global_install_scripts/Readme.md) to guide you through the installation process.

### Installation to a single project

```bash
# Download the talisman binary
curl https://thoughtworks.github.io/talisman/install.sh > ~/install-talisman.sh
chmod +x ~/install-talisman.sh
```

```bash
# Install to a single project (as pre-push hook)
cd my-git-project
~/install-talisman.sh
```

#### Usage with the [pre-commit](https://pre-commit.com) git hooks framework

Add this to your `.pre-commit-config.yaml` (be sure to update `rev` to point to
a real git revision!)

```yaml
-   repo: https://github.com/thoughtworks/talisman
    rev: ''  # Update me!
    hooks:
    # either `commit` or `push` support
    -   id: talisman-commit
    # -   id: talisman-push
```

## Talisman in action

After the installation is successful, Talisman will run checks for obvious secrets automatically before each push:

```bash
$ git push
The following errors were detected in danger.pem
         The file name "danger.pem" failed checks against the pattern ^.+\.pem$

error: failed to push some refs to 'git@github.com:jacksingleton/talisman-demo.git'
```

### Ignoring Files

If you're *really* sure you want to push that file, you can add it to
a `.talismanignore` file in the project root:

```bash
echo 'danger.pem' >> .talismanignore
```

Note that we can ignore files in a few different ways:

* If the pattern ends in a path separator, then all files inside a
  directory with that name are matched. However, files with that name
  itself will not be matched.

* If a pattern contains the path separator in any other location, the
  match works according to the pattern logic of the default golang
  glob mechanism.

* If there is no path separator anywhere in the pattern, the pattern
  is matched against the base name of the file. Thus, the pattern will
  match files with that name anywhere in the repository.

You can also disable only specific detectors.
For example, if your `init-env.sh` filename triggers a warning, you can only disable
this warning while still being alerted if other things go wrong (e.g. file content):
```bash
echo 'init-env.sh # ignore:filename,filesize' >> .talismanignore
```
Note: Here both filename and filesize detectors are ignored for init-env.sh, but
filecontent detector will still activate on `init-env.sh`

At the moment, you can ignore

* `filecontent`
* `filename`
* `filesize`


## Developing locally

To contribute to Talisman, you need a working golang development
environment. Check [this link](https://golang.org/doc/install) to help
you get started with that.

Talisman now uses go modules (GO111MODULE=on) to manage dependencies

Once you have go 1.11 installed and setup, clone the talisman repository. In your
working copy, fetch the dependencies by having go mod fetch them for
you.

```` GO111MODULE=on go mod vendor ````

To run tests ```` GO111MODULE=on go test -mod=vendor ./...  ````

To build Talisman, we can use [gox](https://github.com/mitchellh/gox):

```` gox -osarch="darwin/amd64 linux/386 linux/amd64" ````

Convenince scripts ```./build``` and ```./clean``` perform build and clean-up as mentioned above.

#### Contributing to Talisman

##### Releasing

* Follow the instructions at the end of 'Developing locally' to build the binaries
* Bump the [version in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L10) according to [semver](https://semver.org/) conventions
* Update the [expected hashes in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L16-L18) to match the new binaries you just created (`shasum -b -a256 ...`)
* Make release commit and tag with the new version prefixed by `v` (like `git tag v0.3.0`)
* Push your release commit and tag: `git push && git push --tags`
* [Create a new release in github](https://github.com/thoughtworks/talisman/releases/new), filling in the new commit tag you just created
* Update the install script hosted on github pages: `git checkout gh-pages`, `git checkout master -- install.sh`, `git commit -m ...`

The latest version will now be accessible to anyone who builds their own binaries, downloads binaries directly from github releases, or uses the install script from the website.
