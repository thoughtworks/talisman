<div align="center">
		<img class=logo align=bottom width="25%" height="95%" src="https://github.com/jaydeepc/talisman-html-report/raw/master/img/talisman.png" />
</div>
<h1 align="center">Talisman</h1>
<p align="center">A tool to detect and prevent secrets from getting checked in</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/thoughtworks/talisman)](https://goreportcard.com/report/thoughtworks/talisman)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/thoughtworks/talisman/issues)
[![Build Status](https://github.com/thoughtworks/talisman/actions/workflows/test-and-coverage.yml/badge.svg?branch=main)](https://github.com/thoughtworks/talisman/actions/workflows/test-and-coverage.yml)

## Table of Contents

- [What is Talisman?](#what-is-talisman)
- [Installation](#installation)
  - [Install onto path (recommended approach)](#install-onto-path-recommended-approach)
  - [Installation as a global hook template](#installation-as-a-global-hook-template)
- [Configuring a project](#configuring-a-project)
  - [Using with hook frameworks](#using-with-hook-frameworks)
    - [Pre-commit](#pre-commit)
    - [Husky](#husky)
  - [Directly invoking talisman](#directly-invoking-talisman)
- [Upgrading](#upgrading)
- [Talisman in action](#talisman-in-action)
  - [Validations](#validations)
  - [Ignoring Files](#ignoring-files)
    - [Interactive mode](#interactive-mode)
    - [Ignoring specific detectors](#ignoring-specific-detectors)
    - [Ignoring specific keywords](#ignoring-specific-keywords)
    - [Ignoring multiple files of same type (with wildcards)](#ignoring-multiple-files-of-same-type-with-wildcards)
    - [Ignoring files by specifying language scope](#ignoring-files-by-specifying-language-scope)
    - [Custom search patterns](#custom-search-patterns)
  - [Configuring severity threshold](#configuring-severity-threshold)
  - [Talisman as a CLI utility](#talisman-as-a-cli-utility)
    - [Interactive mode](#interactive-mode-1)
    - [Git history Scanner](#git-history-scanner)
    - [Checksum Calculator](#checksum-calculator)
- [Talisman HTML Reporting](#talisman-html-reporting)
  - [Sample Screenshots](#sample-screenshots)
- [Uninstallation](#uninstallation)
  - [Uninstallation from a global hook template](#uninstallation-from-a-global-hook-template)
  - [Uninstallation from a single repository](#uninstallation-from-a-single-repository)
- [Contributing to Talisman](#contributing-to-talisman)

# What is Talisman?

Talisman is a tool that scans git changesets to ensure that potential secrets
or sensitive information do not leave the developer's workstation.

It validates the outgoing changeset for things that look suspicious - such as
potential SSH keys, authorization tokens, private keys etc.

# Installation

Talisman supports MAC OSX, Linux and Windows.

Talisman can be installed and used in one of the following ways:

1. As a standalone executable
2. As a git hook as a global [git hook template](https://git-scm.com/docs/git-init#_template_directory) and a CLI utility (for git repo scanning)
3. As a git hook into a single git repository

Talisman can be set up as either a pre-commit or pre-push hook on the git repositories.

Find the instructions below.

*Disclaimer: Secrets creeping in via a forced push in a git repository cannot be detected by Talisman. A forced push is believed to be notorious in its own ways, and we suggest git repository admins to apply appropriate measures to authorize such activities.*

## Install onto path (recommended approach)

We recommend installing `talisman` onto your path so that it is available for
git hook frameworks and scripts. Pick the correct binary for your system from
our [Releases Page](https://github.com/thoughtworks/talisman/releases), or run
our [install script](https://github.com/thoughtworks/talisman/blob/main/install.sh):

```bash
bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/install.sh)"
```

Or set environment variable `INSTALL_LOCATION` to specify a custom location for
the binary:

```bash
INSTALL_LOCATION=/usr/local/bin bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/install.sh)"
```

Or set environment variable `VERSION` to a released tag to install a specific version::

```bash
VERSION=v1.36.0 bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/install.sh)"
```

Or using linuxbrew in Linux and homebrew in macOS by running the following
command in terminal:

```bash
brew install talisman
```

## Installation as a global hook template

We offer scripts that will install Talisman as a **pre-commit git hook template**, as that will cause
Talisman to be present, not only in your existing git repositories, but also in any new repository that you 'init' or
'clone'.

1. Run the following command on your terminal, to download and install the binary at $HOME/.talisman/bin

  As a pre-commit hook:

  ```
bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/global_install_scripts/install.bash)"
```

  OR

  As a pre-push hook:

  ```
bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/global_install_scripts/install.bash)" -- pre-push
```


2. If you do not have TALISMAN\_HOME set up in your `$PATH`, you will be asked an appropriate place to set it up. Choose the option number where you set the profile source on your machine.


  Remember to execute *source* on the path file or restart your terminal.
If you choose to set the `$PATH` later, please export TALISMAN\_HOME=$HOME/.talisman/bin to the path.


3. Choose a base directory where Talisman should scan for all git repositories, and set up a git hook (pre-commit or pre-push, as chosen in step 1) as a symlink.
  This script will not clobber pre-existing hooks. If you have existing hooks you can add talisman through a [hook framework](#using-with-hook-frameworks)

  - you can set SEARCH_ROOT environment variable with the path of the base directory before executing the installation so you don't need to input it manually during the installation

# Configuring a project

## Using with hook frameworks

Globally installing talisman as a hook will not clobber any existing hooks. If
the installation script finds any existing hooks, it will only indicate so on
the console. To run multiple hooks we suggest using a hook framework, such as
pre-commit or husky. These instructions assume that the talisman executable is
installed somewhere on your system's path.

### Pre-commit

Use [pre-commit](https://pre-commit.com) to manage existing hooks along with
Talisman. Reference our [pre-commit-hooks](.pre-commit-hooks.yaml) in your
`.pre-commit-config.yaml`:

```yaml
-   repo: https://github.com/thoughtworks/talisman
    rev: 'v1.32.2'  # Update me!
    hooks:
      # both pre-commit and pre-push supported
      # -   id: talisman-push
      - id: talisman-commit
```

### Husky

[husky](https://typicode.github.io/husky) is an npm module for managing hooks.
Add the following line to the husky pre-commit configuration in your
`package.json`:

```
talisman --githook pre-commit
```

## Directly invoking talisman

Once the talisman executable is [installed](#install-onto-path-recommended-approach)
you can configure a standalone pre-commit hook for a git repository:

```bash
cd my-git-project
echo "talisman -g pre-commit" >> .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

# Upgrading
Since release v0.4.4, Talisman <b>automatically updates</b> the binary to the latest release, when the hook is invoked (at pre-commit/pre-push, as set up). So, just sit back, relax, and keep using the latest Talisman without any extra efforts.

The following environment variables can be set:

1. TALISMAN_SKIP_UPGRADE: Set to true if you want to skip the automatic upgrade check. Default is false
2. TALISMAN_UPGRADE_CONNECT_TIMEOUT: Max connect timeout before the upgrade is cancelled(in seconds). Default is 10 seconds.

If at all you need to manually upgrade, here are the steps:
<br>[Recommended] Update Talisman binary and hook scripts to the latest release:

```bash
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/main/global_install_scripts/update_talisman.bash > /tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash
```


Update only Talisman binary by executing:

```bash
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/main/global_install_scripts/update_talisman.bash > /tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash talisman-binary
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
| danger.pem      | Expected file to not contain hex encoded texts such as:                       |
|                 | awsSecretKey=c64e8c79aacf5ddb02f1274db2d973f363f4f553ab1692d8d203b4cc09692f79 |
+-----------------+-------------------------------------------------------------------------------+
```

In the above example, the file *danger.pem* has been flagged as a security breach due to the following reasons:

* The filename matches one of the pre-configured patterns.
* The file contains an awsSecretKey which is scanned and flagged by Talisman

If you have installed Talisman as a pre-commit hook, it will scan only the _diff_ within each commit. This means that it would only report errors for parts of the file that were changed.

In case you have installed Talisman as a pre-push hook, it will scan the complete file in which changes are made. As mentioned above, it is recommended that you use Talisman as a **pre-commit hook**.

## Validations
The following detectors execute against the changesets to detect secrets/sensitive information:

* **Encoded values** - scans for encoded secrets in Base64, hex etc.
* **File content** - scans for suspicious content in file that could be potential secrets or passwords
* **File size** - scans for large files that may potentially contain keys or other secrets
* **Entropy** - scans for content with high entropy that are likely to contain passwords
* **Credit card numbers** - scans for content that could be potential credit card numbers
* **File names** - scans for file names and extensions that could indicate them potentially containing secrets, such as keys, credentials etc.


## Ignoring Files

If you're *really* sure you want to push that file, you can configure it into the `.talismanrc` file in the project root. The contents required for ignoring your failed files will be printed by Talisman on the console immediately after the Talisman Error Report:


```bash
If you are absolutely sure that you want to ignore the above files from talisman detectors, consider pasting the following format in .talismanrc file in the project root
fileignoreconfig:
- filename: danger.pem
  checksum: cf97abd34cebe895417eb4d97fbd7374aa138dcb65b1fe7f6b6cc1238aaf4d48
  ignore_detectors: []
```
Entering this in the `.talismanrc` file will ensure that Talisman will ignore the `danger.pem` file as long as the checksum matches the value mentioned in the `checksum` field.

### Interactive mode

**Available only for non-Windows users**

If it is too much of a hassle to keep copying content to .talismanrc everytime you encounter an error from Talisman, you could enable the interactive mode and let Talisman assist you in prompting the additions of the files to ignore.
Just follow the simple steps:
1. Open your bash profile where your environment variables are set (.bashrc, .bash_profile, .profile or any other location)
2. You will see `TALISMAN_INTERACTIVE` variable under `# >>> talisman >>>`
3. If not already set to true, add `export TALISMAN_INTERACTIVE=true`
4. Don't forget to save and source the file

That's it! Every time Talisman hook finds an error during pre-push/pre-commit, just follow the instructions as Talisman suggests.
Be careful to not ignore a file without verifying the content. You must be confident that no secret is getting leaked out.

### Ignoring specific detectors

Below is a detailed description of the various fields that can be configured into the `.talismanrc` file:

* `filename` : This field should mention the fully qualified filename.
* `checksum` : This field should always have the value specified by Talisman in the message displayed above. If at any point, a new change is made to the file, it will result in a new checksum and Talisman will scan the file again for any potential security threats.
* `ignore_detectors` : This field will disable specific detectors for a particular file.
For example, if your `init-env.sh` filename triggers a warning, you can only disable
this warning while still being alerted if other things go wrong (e.g. file content):


```yaml
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

### Ignoring specific keywords

Because some of your files might contain keywords such as `key` or `pass` that are not necessarily related to secrets, you might want to ignore these keywords to reduce the number of false positives.
This can be achieved by using the `allowed_patterns` field at the file level and/or at the repository level:

```yaml
fileignoreconfig:
- filename: test
  allowed_patterns: [key]
allowed_patterns:
- keyword
- pass
```

In the previous example, `key` is allowed in the `test` file, `keyword` and `pass` are allowed at the repository level.

The `allowed_patterns` field also supports Golang regular expressions. Here is a simple code example where Golang RegExp can be useful:

```sh
export AWS_ACCESS_KEY_ID = AKIAIO5FODNN7EXAMPLE
export AWS_ACCESS_KEY_ID=$(vault read -field=value path/to/aws-access-key-id)
```

By default, Talisman will alert for both lines. In the second line, we are extracting the AWS Access Key ID from Hashicorp Vault which doesn't expose the secret to the code. If this type of usage is common in your code, you might want to tell Talisman to not alert when you use a Vault. This can be achieved with a configuration like:

```yaml
allowed_patterns:
- export\ AWS[ \w]*KEY[ \w]*=.*vault\ read.*
```

### Ignoring multiple files of same type (with wildcards)

You can choose to ignore all files of a certain type, because you know they will always be safe, and you wouldn't want Talisman to scan them.

Steps:

1. Format a wildcard pattern for the files you want to ignore. For example, `*.lock`
2. Use the [checksum calculator](#checksum-calculator) to feed the pattern and attain a collective checksum. For example, `talisman --checksum="*.lock" `
3. Copy the fileconfig block, printed on console, to .talismanrc file.

If any of the files are modified, talisman will scan the files again, unless you re-calculate the new checksum and replace it in .talismanrc file.

### Ignoring files by specifying language scope

You can choose to ignore files by specifying the language scope for your project in your talismanrc.

```yaml
scopeconfig:
  - scope: go
  - scope: node
  - scope: images
  - scope: php
  - scope: python
```

Talisman is configured to ignore certain files based on the specified scope. For example, mentioning the node scope in the scopeconfig will prevent talisman from scanning files such as the yarn.lock or package-lock.json.

You can specify multiple scopes.

Currently .talismanrc only supports scopeconfig support for go, node, php and images. Other scopes will be added shortly.

### Custom search patterns

You can specify custom regex patterns to look for in the current repository

```yaml
custom_patterns:
- pattern1
- pattern2
```

<br/><i>
**Note**: The use of .talismanignore has been deprecated. File .talismanrc replaces it because:

* .talismanrc has a much more legible yaml format
* It also brings in more secure practices with every modification of a file with a potential sensitive value to be reviewed
* The new format also brings in the extensibility to introduce new usable functionalities. Keep a watch out for more </i>

## Configuring severity threshold

Each validation is associated with a severity
1. Low
2. Medium
3. High

You can specify a threshold in your .talismanrc:

```yaml
threshold: medium
```
This will report all Medium severity issues and higher (Potential risks that are below the threshold will be reported in the warnings)

1. A list of all risks with their severity level can be found in this [configuration file](detector/severity/severity_config.go).
2. By default, the threshold is set to low.
3. Any custom search patterns you add, are considered to be of high severity.

## Configuring custom severities

You can customize the [security levels](detector/severity/severity_config.go) of the detectors provided by Talisman in the .talismanrc file:

```yaml
custom_severities:
- detector: Base64Content
  severity: medium
- detector: HexContent
  severity: low
```

By using custom severities and a severity threshold, Talisman can be configured to alert only on what is important based on your context. This can be useful to reduce the number of false positives.

## Talisman as a CLI utility

If you execute `talisman` on the command line, you will be able to view all the parameter options you can pass

```
  -c, --checksum string          checksum calculator calculates checksum and suggests .talismanrc format
  -d, --debug                    enable debug mode (warning: very verbose)
  -g, --githook string           either pre-push or pre-commit (default "pre-push")
      --ignoreHistory            scanner scans all files on current head, will not scan through git commit history
  -i, --interactive              interactively update talismanrc (only makes sense with -g/--githook)
  -p, --pattern string           pattern (glob-like) of files to scan (ignores githooks)
  -r, --reportdirectory string   directory where the scan reports will be stored
  -s, --scan                     scanner scans the git commit history for potential secrets
  -w, --scanWithHtml             generate html report (**Make sure you have installed talisman_html_report to use this, as mentioned in Readme**)
  -v, --version                  show current version of talisman
```

### Interactive mode

When you regularly have too many files that get are flagged by talisman hook, which you know should be fine to check in, you can use this feature to let talisman ease the process for you. The interactive mode will allow Talisman to prompt you to directly add files you want to ignore to .talismanrc from command prompt directly.
To enable this feature, you need TALISMAN_INTERACTIVE variable to be set as true in your bash file.

You can invoke talisman in interactive mode by either of the 2 ways:
1.  Open your bash file, and add
```export TALISMAN_INTERACTIVE=true```
Don't forget to source the bash file for the variable to take effect!

2.  Alternatively, you can also invoke the interactive mode by using the CLI utility
(for using pre-commit hook)
```talisman -i -g pre-commit```

*Note*: If you use an IDE's Version Control integration for git operations, this feature will not work. You can still use the suggested filename and checksum to be entered in .talismanrc  file manually.

### Git history Scanner

You can now execute Talisman from CLI, and potentially add it to your CI/CD pipelines, to scan git history of your repository to find any sensitive content.
This includes scanning of the files listed in the .talismanrc file as well.

**Steps**:

 1. Get into the git directory path to be scanned `cd <directory to scan>`
 2. Run the scan command `talisman --scan`
  * Running this command will create a folder named <i>talisman_reports</i> in the root of the current directory and store the report files there.
  * You can also specify the location for reports by providing an additional parameter as <i>--reportdirectory</i> or <i>--rd</i>
<br>For example, `talisman --scan --reportdirectory=/Users/username/Desktop`

You can use the other options to scan as given above.


<i>Talisman currently does not support ignoring of files for scanning.</i>



### Checksum Calculator

Talisman Checksum calculator gives out yaml format which you can directly copy and paste in .talismanrc file in order to ignore particular file formats from talisman detectors.

To run the checksum please "cd" into the root of your repository and run the following command

For Example:
`talisman --checksum="*.pem *.txt"`

1. This command finds all the .pem files in the repository and calculates collective checksum of all those files and outputs a yaml format for .talismanrc. In the same way it deals with the .txt files.
2. Multiple file names / patterns can be given with space separation.

Example output:

	.talismanrc format for given file names / patterns
	fileignoreconfig:
	- filename: '*.pem'
	  checksum: f731b26be086fd2647c40801630e2219ef207cb1aacc02f9bf0559a75c0855a4
	  ignore_detectors: []
	- filename: '*.txt'
	  checksum: d9e9e94868d7de5b2a0706b8d38d0f79730839e0eb4de4e9a2a5a014c7c43f35
	  ignore_detectors: []


Note: Checksum calculator considers the staged files while calculating the collective checksum of the files.

# Talisman HTML Reporting
<i>Powered by 		<a href="https://jaydeepc.github.io/report-mine-website/"><img class=logo align=bottom width="10%" height="10%" src="https://github.com/jaydeepc/talisman-html-report/raw/master/img/logo_reportmine.png" /></a></i>

Talisman CLI tool `talisman` also comes with the capability to provide detailed and sharable HTML report. Once you have installed Talisman, please follow the steps mentioned in [talisman-html-report](https://github.com/jaydeepc/talisman-html-report), to install the reporting package in `.talisman` folder. To generate the html report, run:

* `talisman --scanWithHtml`

This will scan the repository and create a folder `talisman_html_report` under the scanned repository. We need to start an HTTP server inside this repository to access the report.Below is a recommended approach to start a HTTP server:

* `python -m SimpleHTTPServer <port> (eg: 8000)`

You can now access the report by navigating to:

`http://localhost:8000`

## Sample Screenshots

* Welcome

<img width="100%" height="70%" src="https://github.com/jaydeepc/talisman-html-report/raw/master/sample/summary.png" />

* Summary

<img width="100%" height="70%" src="https://github.com/jaydeepc/talisman-html-report/raw/master/sample/execution-summary.png" />

* Detailed Report

<img width="100%" height="70%" src="https://github.com/jaydeepc/talisman-html-report/raw/master/sample/detailed.png" />

* Error Report

<img width="100%" height="70%" src="https://github.com/jaydeepc/talisman-html-report/raw/master/sample/error-report.png" />

<i> **Note**: You don't have to start a server if you are running Talisman in CI or any other hosted environment </i>


# Uninstallation
The uninstallation process depends on how you had installed Talisman.
You could have chosen to install as a global hook template or at a single repository.

Please follow the steps below based on which option you had chosen at installation.

## Uninstallation from a global hook template
Run the following command on your terminal to uninstall talisman globally from your machine.

For pre-commit hook:

```
bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/global_install_scripts/uninstall.bash)"
```

For pre-push hook:

```
bash -c "$(curl --silent https://raw.githubusercontent.com/thoughtworks/talisman/main/global_install_scripts/uninstall.bash)" -- pre-push
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

To contribute to Talisman, have a look at our [contributing guide](contributing.md).
