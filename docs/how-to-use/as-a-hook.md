---
layout: default
title: Use As a Hook
parent: How To Use
nav_order: 1
description: "Use As a Hook" 
---
# Talisman As a Hook

After the installation is successful, Talisman will run checks for obvious secrets automatically before each commit or push (as chosen during installation). In case there are any possible secret leakages detected, talisman will display a detailed report of the errors:

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

In the above example, the file *danger.pem* has been flagged as a secret leak due to the following reasons:

* The filename matches one of the pre-configured patterns.
* The file contains an awsSecretKey which is scanned and flagged by Talisman

If you have installed Talisman as a pre-commit hook, it will scan only the _diff_ within each commit. This means that it would only report errors for parts of the file that were changed.

In case you have installed Talisman as a pre-push hook, it will scan the complete file in which changes are made. As mentioned above, it is recommended that you use Talisman as a **pre-commit hook**.
