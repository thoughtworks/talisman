<div style="display: flex; justify-content: center;">
	<h1 align="center">
		<img class=logo align=bottom width="5%" height="5%" src="https://thoughtworks.github.io/talisman/logo.svg" />
		Talisman</h1>
</div>
<p align="center">A tool to detect and prevent secrets from getting checked in</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/thoughtworks/talisman)](https://goreportcard.com/report/thoughtworks/talisman) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/thoughtworks/talisman/issues) [![Build Status](https://travis-ci.org/thoughtworks/talisman.svg?branch=master)](https://travis-ci.org/thoughtworks/talisman)


## Table of Contents
- [What is Talisman?](#what-is-talisman)
- [Installation](#installation)
	- [As a global hook template (Recommended)](#installation-as-a-global-hook-template)
	- [To a single repository](#installation-to-a-single-project)
	- [As a CLI to find file types](#installation-as-a-cli)
- [Upgrading Talisman](#Upgrading)
- [Talisman in action](#talisman-in-action)
	- [Validations](#validations) 
	- [Ignoring files](#ignoring-files)
  - [Scanning Git hisotry](#scanning-git-history)
- [Uninstallation](#uninstallation)
	- [From a global hook template](#uninstallation-from-a-global-hook-template)
	- [From a single repository](#uninstallation-from-a-single-repository)   
- [Contributing to Talisman](#contributing-to-talisman)
	- [Developing locally](#developing-locally)
	- [Releasing](#releasing)  

# What is Talisman?
Talisman is a tool that installs a hook to your repository to ensure that potential secrets or sensitive information do not leave the developer's workstation. 

It validates the outgoing changeset for things that look suspicious - such as potential SSH
keys, authorization tokens, private keys etc.

# Installation

Talisman supports MAC OSX, Linux and Windows.

Talisman can be installed and used in one of three different ways:

1. As a git hook as a global [git hook template](https://git-scm.com/docs/git-init#_template_directory)
2. As a git hook into a single git repository
3. As a CLI with the `--pattern` argument to find files

Talisman can be set up as either a pre-commit or pre-push hook on the git repositories.

Find the instructions below.

## [Recommended approach]
## Installation as a global hook template 

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

2. If you do not have TALISMAN\_HOME set up in your `$PATH`, you will be asked an appropriate place to set it up. Choose the option number where you set the profile source on your machine.


  Remember to execute *source* on the path file or restart your terminal.
If you choose to set the `$PATH` later, please export TALISMAN\_HOME=$HOME/.talisman/bin to the path.


3. Choose a base directory where Talisman should scan for all git repositories, and setup a git hook (pre-commit or pre-push, as chosen in step 1) as a symlink.
  This script will not clobber pre-existing hooks. If you have existing hooks, [look for ways to chain Talisman into them.] (#handling-existing-hooks)


### Handling existing hooks
Installation of Talisman globally does not clobber pre-existing hooks on repositories. <br>
If the installation script finds any existing hooks, it will only indicate so on the console. <br>
To achieve running multiple hooks we suggest (but not limited to) the following two tools

#### 1. Pre-commit (Linux/Unix)
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

#### 2. Husky (Linux/Unix/Windows)
[husky](https://github.com/typicode/husky/blob/master/DOCS.md) is an npm module for managing git hooks.
In order to use husky, make sure you have already set TALISMAN_HOME to `$PATH`.
 
+ **Existing Users**
 
 If you already are using husky, add the following lines to husky pre-commit in package.json
 
 ###### Windows
 
 ``` 
    "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" -githook pre-commit'" 
```
 
 ###### Linux/Unix
 
 ```
    $TALISMAN_HOME/talisman_hook_script pre-commit
```
+ **New Users**

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

### Handling existing hooks
Talisman will need to be chained with any existing git hooks.You can use [pre-commit](https://pre-commit.com) git hooks framework to handle this.

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

## Installation as a CLI
1. Download the Talisman binary from the [Releases page](https://github.com/thoughtworks/talisman/releases) corresponding to your system type
2. Place the binary somewhere (either directly in your repository, or by putting it somewhere in your system and adding it to your `$PATH`)
3. Run talisman with the `--pattern` argument (matches glob-like patterns, [see more](https://github.com/bmatcuk/doublestar#patterns))

```bash
# finds all .go and .md files in the current directory (recursively) 
talisman --pattern="./**/*.{go,md}"
```
# Upgrading
To update Talisman to the latest release, run the following curl command:
```bash
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/update_talisman.bash > /tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash
```

# Talisman in action

After the installation is successful, Talisman will run checks for obvious secrets automatically before each commit or push (as chosen during installation). In case there are any security breaches detected, talisman will display a detailed report of the errors:

```bash
$ git push
Talisman Report:
+-----------------+-------------------------------------------------------------------------------+
|     FILE        |                                    ERRORS                                     |
+-----------------+-------------------------------------------------------------------------------+
| danger.pem      | The file name "danger.pem"                                                    |
|                 | failed checks against the                                                     |
|                 | pattern ^.+\.pem$                                                             |
+-----------------+-------------------------------------------------------------------------------+
| danger.pem      | Expected file to not to contain hex encoded texts such as:                    |
|                 | awsSecretKey=c64e8c79aacf5ddb02f1274db2d973f363f4f553ab1692d8d203b4cc09692f79 |
+-----------------+-------------------------------------------------------------------------------+
```

In the above example, the file *danger.pem* has been flagged as a security breach due to the following reasons:

* The filename matches one of the pre-configured patterns.
* The file contains an awsSecretKey which is scanned and flagged by Talisman

## Validations
The following detectors execute against the changesets to detect secrets/sensitive information:

* **Encoded values** - scans for encoded secrets in Base64, hex etc.
* **File content** - scans for suspicious content in file that could be potential secrets or passwords
* **File size** - scans for large files that may potentially contain keys or other secrets
* **Entropy** - scans for content with high entropy that are likely to contain passwords
* **Credit card numbers** - scans for content that could be potential credit card numbers
* **File names** - scans for file names and extensions that could indicate them potentially containing secrets, such as keys, credentials etc.


## Ignoring Files

If you're *really* sure you want to push that file, you can configure it into the `.talismanrc` file in the project root. The contents required for ignoring your failed files will be printed by Talisman on the console immediately after the Talisman Error report:


```bash
If you are absolutely sure that you want to ignore the above files from talisman detectors, consider pasting the following format in .talismanrc file in the project root
fileignoreconfig:
- filename: danger.pem
  checksum: cf97abd34cebe895417eb4d97fbd7374aa138dcb65b1fe7f6b6cc1238aaf4d48
  ignore_detectors: []
```
Entering this in the `.talismanrc` file will ensure that Talisman will ignore the `danger.pem` file as long as the checksum matches the value mentioned in the `checksum` field.  

Below is a detailed description of the various fields that can be configured into the `.talismanrc` file:

* `filename` : This field should mention the fully qualified filename.
* `checksum` : This field should always have the value specified by Talisman in the message displayed above. If at any point, a new change is made to the file, it will result in a new checksum and Talisman will scan the file again for any potential security threats.
* `ignore_detectors` : This field will disable specific detectors for a particular file.
For example, if your `init-env.sh` filename triggers a warning, you can only disable
this warning while still being alerted if other things go wrong (e.g. file content):

```bash
fileignoreconfig:
- filename: init-env.sh
  checksum: cf97abd34cebe895417eb4d97fbd7374aa138dcb65b1fe7f6b6cc1238aaf4d48
  ignore_detectors: [filename, filesize]
```

Note: Here both filename and filesize detectors are ignored for init-env.sh, but
filecontent detector will still activate on `init-env.sh`

At the moment, you can ignore

* `filecontent`
* `filename`
* `filesize`


## Checksum Calculator for .talismanrc

Talisman gives out yaml format which you can directly copy and paste in .talismanrc file in order to ignore particular file formats from talisman detectors.

To run the checksum please "cd" into the root of your repository and run the following command

For Example:
* `talisman --checksum="*.pem *.txt"`

1. This command finds all the .pem files in the respository and calculates collective checksum of all those files and outputs a yaml format for .talismanrc. In the same way it deals with the .txt files.
2. Multiple file names / patterns can be given with space seperation.

Example output:
```.talismanrc format for given file names / patterns
fileignoreconfig:
- filename: '*.pem'
  checksum: f731b26be086fd2647c40801630e2219ef207cb1aacc02f9bf0559a75c0855a4
  ignore_detectors: []
- filename: '*.txt'
  checksum: d9e9e94868d7de5b2a0706b8d38d0f79730839e0eb4de4e9a2a5a014c7c43f35
  ignore_detectors: []
```

## Scanning Git history

Talisman also scans the content present in the git history of the repository, this includes scanning of the files listed in the .talismanrc file as well.

To run the scanner please "cd" into the directory to be scanned and run the following command

* `talisman scan`

<i>Talisman currently does not support ignoring of files for scanning.</i>

# Uninstallation
The uninstallation process depends on how you had installed Talisman.
You could have chosen to install as a global hook template or at a single repository.

Please follow the steps below based on which option you had chosen at installation.

## Uninstallation from a global hook template
Run the following command on your terminal to uninstall talisman globally from your machine.

For pre-commit hook:

```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/uninstall.bash > /tmp/uninstall_talisman.bash && /bin/bash /tmp/uninstall_talisman.bash 
```

For pre-push hook:

```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/uninstall.bash > /tmp/uninstall_talisman.bash && /bin/bash /tmp/uninstall_talisman.bash pre-push
```

This will

1. ask you for the base dir of all your repos, find all git repos inside it and remove talisman hooks
2. remove talisman hook from .git-template 
3. remove talisman from the central install location ($HOME/.talisman/bin).<br>

<i>You will have to manually remove TALISMAN_HOME from your environment variables</i>

## Uninstallation from a single repository
When you installed Talisman, it must have created a pre-commit or pre-push hook (as selected) in your repository during installation. 

You can remove the hook manually by deleting the Talisman pre-commit or pre-push hook from .git/hooks folder in repository.

# Contributing to Talisman

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


## Releasing

* Follow the instructions at the end of 'Developing locally' to build the binaries
* Bump the [version in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L10) according to [semver](https://semver.org/) conventions
* Update the [expected hashes in install.sh](https://github.com/thoughtworks/talisman/blob/d4b1b1d11137dbb173bf681a03f16183a9d82255/install.sh#L16-L18) to match the new binaries you just created (`shasum -b -a256 ...`)
* Make release commit and tag with the new version prefixed by `v` (like `git tag v0.3.0`)
* Push your release commit and tag: `git push && git push --tags`
* [Create a new release in github](https://github.com/thoughtworks/talisman/releases/new), filling in the new commit tag you just created
* Update the install script hosted on github pages: `git checkout gh-pages`, `git checkout master -- install.sh`, `git commit -m ...`

The latest version will now be accessible to anyone who builds their own binaries, downloads binaries directly from github releases, or uses the install script from the website.
