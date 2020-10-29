---
layout: default
title: Use As a Git Scanner
parent: How To Use
nav_order: 2
description: "Use As a Git Scanner" 
---

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

# HTML Reporting
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
