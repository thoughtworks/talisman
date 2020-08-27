---
layout: default
title: Talisman In Action
nav_order: 4
description: "How Talisman works" 
permalink: /details
---

# Talisman In Action

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
```

Talisman is configured to ignore certain files based on the specified scope. For example, mentioning the node scope in the scopeconfig will prevent talisman from scanning files such as the yarn.lock or package-lock.json.

You can specify multiple scopes.

Currently .talismanrc only supports scopeconfig support for go and node. Other scopes will be added shortly.

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
1. low
2. medium
3. high

You can specify a threshold in your .talismanrc: 

```yaml
threshold: medium
```
This will report all Medium severity issues and higher (Potential risks that are below the threshold will be reported in the warnings)

By default, the threshold is set to low

## Talisman as a CLI utility

If you execute `talisman` on the command line, you will be able to view all the parameter options you can pass

```
  -c, --checksum string          checksum calculator calculates checksum and suggests .talismanrc format
  -d, --debug                    enable debug mode (warning: very verbose)
  -g, --githook string           either pre-push or pre-commit (default "pre-push")
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
  * You can also specify the location for reports by providing an additional parameter as <i>--reportDirectory</i> or <i>--rd</i>
<br>For example, `talisman --scan --reportdirectory=/Users/username/Desktop`

You can use the other options to scan as given above.


<i>Talisman currently does not support ignoring of files for scanning.</i>



### Checksum Calculator

Talisman Checksum calculator gives out yaml format which you can directly copy and paste in .talismanrc file in order to ignore particular file formats from talisman detectors.

To run the checksum please "cd" into the root of your repository and run the following command

For Example:
`talisman --checksum="*.pem *.txt"`

1. This command finds all the .pem files in the respository and calculates collective checksum of all those files and outputs a yaml format for .talismanrc. In the same way it deals with the .txt files.
2. Multiple file names / patterns can be given with space seperation.

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

This will scan the repository and create a folder `talisman_html_report` under the the scanned repository. We need to start an HTTP server inside this repository to access the report.Below is a recommended approach to start a HTTP server:

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
