<div style="display: flex; justify-content: center;">
	<h1 align="center">
		<img class=logo align=bottom width="5%" height="5%" src="https://thoughtworks.github.io/talisman/logo.svg" />
		Talisman</h1>
</div>
<p align="center">A tool to detect and prevent secrets from getting checked in</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/thoughtworks/talisman)](https://goreportcard.com/report/thoughtworks/talisman) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/thoughtworks/talisman/issues)


## Table of Contents
- [What is Talisman?] (#what-is-talisman)
- [Installation](#installation)
	- [As a global hook template (Recommended)](#installation-as-a-global-hook-template)
	- [To a single repository](#installation-to-a-single-project) 
- [Talisman in action](#talisman-in-action)
	- [Validations](#validations) 
	- [Ignoring files](#ignoring-files)
- [Uninstallation](#uninstallation)
	- [From a global hook template](#uninstallation-from-a-global-hook-template)
	- [From a single repository](#uninstallation-from-a-single-project)   
- [Contributing to Talisman](#contributing-to-talisman)
	- [Developing locally](#developing-locally)
	- [Releasing](#releasing)  

# What is Talisman?
Talisman is a tool that installs a hook to your repository to ensure that potential secrets or sensitive information do not leave the developer's workstation. 

It validates the outgoing changeset for things that look suspicious - such as potential SSH
keys, authorization tokens, private keys etc.


# Installation

Talisman supports MAC OSX, Linux and Windows.

The recommended approach of installation is as a global
[git hook template](https://git-scm.com/docs/git-init#_template_directory). You could also choose to install it at a per-repo level instead.

Within each approach, Talisman can be set up as either a pre-commit or pre-push hook on the git repositories.

Find the instructions below.

## Installation as a global hook template 
### [Recommended approach]
We recommend installing Talisman as a git hook template, as that will cause
Talisman to be present, not only in your existing git repositories, but also in any new repository that you 'init' or
'clone'.

1. Run the following command on your terminal, to download and install the binary at $HOME/.talisman/bin

  As a pre-commit hook:

  ```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash 
```

  OR

  As a pre-push hook:

  ```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash pre-push
```

2. If you do not have TALISMAN\_HOME set up in your path, you will be asked an appropriate place to set it up. Choose the option number where you set the profile source on your machine.


  Remember to execute *source* on the path file or restart your terminal.
If you choose to set the path later, please export TALISMAN\_HOME=$HOME/.talisman/bin to the path.


3. Choose a base directory where Talisman should scan for all git repositories, and setup a git hook (pre-commit or pre-push, as chosen in step 1) as a symlink.
  This script will not clobber pre-existing hooks. If you have existing hooks, [look for ways to chain Talisman into them.] (#handling-existing-hooks)


#### Handling existing hooks
Installation of Talisman globally does not clobber pre-existing hooks on repositories. 
If the installation script finds any existing hooks, it will only indicate so on the console.
To achieve running multiple hooks we suggest the following two tools

##### 1. Pre-commit (Linux/Unix)
Use [pre-commit](https://pre-commit.com) tool to manage all the existing hooks along with Talisman.
In the suggestion, it will prompt the following code to be included in .pre-commit-config.yaml
```
    -   repo: local
        hooks:
        -   id: talisman-precommit
            name: talisman
            entry: bash -c 'if [ -n "${TALISMAN_HOME:-}" ]; then ${TALISMAN_HOME}/talisman_hook_script pre-commit; else echo "TALISMAN does not exist. Consider installing from https://github.com/thoughtworks/talisman . If you already have talisman installed, please ensure TALISMAN_HOME variable is set to where talisman_hook_script resides, for example, TALISMAN_HOME=${HOME}/.talisman/bin"; fi'
            language: system
            pass_filenames: false
            types: [text]
            verbose: true
```

##### 2. Husky (Linux/Unix/Windows)
[husky](https://github.com/typicode/husky/blob/master/DOCS.md) is an npm module for managing git hooks.
In order to use husky, make sure you set TALISMAN_HOME.
 
##### Existing Users
 If you already are using husky, add the following lines to husky pre-commit in package.json
###### Windows
``` 
    "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" -githook pre-commit'" 
```
###### Linux/Unix
```
    $TALISMAN_HOME/talisman_hook_script pre-commit
```
##### New Users
If you want to use husky with multiple hooks along with talisman, add the following snippet to you package json.
    
###### Windows
```
     {
        "husky": {
          "hooks": {
            "pre-commit": "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" -githook pre-commit'" && "other-scripts"
            }
        }
    }
```
###### Linux/Unix
```
    {
      "husky": {
       "hooks": {
         "pre-commit": "$TALISMAN_HOME/talisman_hook_script pre-commit" && "other-scripts"
          }
        }
      }
```


## Installation to a single project

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

### Usage with the [pre-commit](https://pre-commit.com) git hooks framework

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

After the installation is successful, Talisman will run checks for obvious secrets automatically before each commit or push (as chosen during installation):

```bash
$ git push
The following errors were detected in danger.pem
         The file name "danger.pem" failed checks against the pattern ^.+\.pem$

error: failed to push some refs to 'git@github.com:jacksingleton/talisman-demo.git'
```

### Validations
The following detectors execute against the changesets to detect secrets/sensitive information:

* **Encoded values** - scans for encoded secrets in Base64, hex etc.
* **File content** - scans for suspicious content in file that could be potential secrets or passwords
* **File size** - scans for large files that may potentially contain keys or other secrets
* **Entropy** - scans for content with high entropy that are likely to contain passwords
* **Credit card numbers** - scans for content that could be potential credit card numbers
* **File names** - scans for file names and extensions that could indicate them potentially containing secrets, such as keys, credentials etc.


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

## Uninstallation
To uninstall talisman globally from your machine, run:
```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/uninstall.bash > /tmp/uninstall_talisman.bash && /bin/bash /tmp/uninstall_talisman.bash 
```
This will
1. ask you for the base dir of all your repos, find all git repos inside it and remove talisman hooks
2. remove talisman hook from .git-template 
3. remove talisman from the central install location ($HOME/.talisman/bin)
You will have to manually remove TALISMAN_HOME from your environment variables


## Contributing to Talisman

### Developing locally

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


### Releasing

* Follow the instructions at the end of 'Developing locally' to build the binaries
* Bump the [version in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L10) according to [semver](https://semver.org/) conventions
* Update the [expected hashes in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L16-L18) to match the new binaries you just created (`shasum -b -a256 ...`)
* Make release commit and tag with the new version prefixed by `v` (like `git tag v0.3.0`)
* Push your release commit and tag: `git push && git push --tags`
* [Create a new release in github](https://github.com/thoughtworks/talisman/releases/new), filling in the new commit tag you just created
* Update the install script hosted on github pages: `git checkout gh-pages`, `git checkout master -- install.sh`, `git commit -m ...`

The latest version will now be accessible to anyone who builds their own binaries, downloads binaries directly from github releases, or uses the install script from the website.
