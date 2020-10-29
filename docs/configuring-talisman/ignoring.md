---
layout: default
title: Ignore From Scan
parent: Configuring Talisman
nav_order: 1
description: "Ignore From Scan" 
---

# Ignoring Files

If you're *really* sure you want to push that file, you can configure it into the `.talismanrc` file in the project root. The contents required for ignoring your failed files will be printed by Talisman on the console immediately after the Talisman Error Report:


```bash
If you are absolutely sure that you want to ignore the above files from talisman detectors, consider pasting the following format in .talismanrc file in the project root
fileignoreconfig:
- filename: danger.pem
  checksum: cf97abd34cebe895417eb4d97fbd7374aa138dcb65b1fe7f6b6cc1238aaf4d48
  ignore_detectors: []
```
Entering this in the `.talismanrc` file will ensure that Talisman will ignore the `danger.pem` file as long as the checksum matches the value mentioned in the `checksum` field.

## Ignoring specific detectors

Below is a detailed description of the various fields that can be configured into the `.talismanrc` file:

* `filename` : This field should mention the fully qualified filename.
* `checksum` : This field should always have the value specified by Talisman in the message displayed above. If at any point, a new change is made to the file, it will result in a new checksum and Talisman will scan the file again for any potential security threats. If needed, you can also [calculate the checksum](../how-to-use/as-a-utility/checksum-calc) again.
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

## Ignoring specific keywords

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

## Ignoring multiple files of same type (with wildcards)

You can choose to ignore all files of a certain type, because you know they will always be safe, and you wouldn't want Talisman to scan them.

Steps:

1. Format a wildcard pattern for the files you want to ignore. For example, `*.lock`
2. Use the [checksum calculator](#checksum-calculator) to feed the pattern and attain a collective checksum. For example, `talisman --checksum="*.lock" `
3. Copy the fileconfig block, printed on console, to .talismanrc file.

If any of the files are modified, talisman will scan the files again, unless you re-calculate the new checksum and replace it in .talismanrc file.

## Ignoring files by specifying language scope

You can choose to ignore files by specifying the language scope for your project in your talismanrc.

```yaml
scopeconfig:
  - scope: go
  - scope: node
```

Talisman is configured to ignore certain files based on the specified scope. For example, mentioning the node scope in the scopeconfig will prevent talisman from scanning files such as the yarn.lock or package-lock.json.

You can specify multiple scopes.

Currently .talismanrc only supports scopeconfig support for go and node. Other scopes will be added shortly.
