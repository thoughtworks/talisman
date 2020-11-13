---
layout: default
title: How To Use
nav_order: 3
description: "How Talisman works" 
has_children: true
permalink: docs/how-to-use

---

# How To Use

Talisman acts as as a hook (pre-commit or pre-push) to your git repository on the developer's machine itself. 
<br> However, you can't always guarantee that all checkins made to the git repository came through a similar validation. Which is why you might also want to run Talisman to scan your complete git history.
<br> You will be able to find more details about how to configure and use them at:
1. [Talisman as a hook](./as-a-hook) 
2. [Talisman as a git scanner](./as-a-git-scanner)
3. [Other utilities](./as-a-utility)

## Validations
The following detectors execute against the changesets to detect secrets/sensitive information:

* **Encoded values** - scans for encoded secrets in Base64, hex etc.
* **File content** - scans for suspicious content in file that could be potential secrets or passwords
* **File size** - scans for large files that may potentially contain keys or other secrets
* **Entropy** - scans for content with high entropy that are likely to contain passwords
* **Credit card numbers** - scans for content that could be potential credit card numbers
* **File names** - scans for file names and extensions that could indicate them potentially containing secrets, such as keys, credentials etc.



### Define custom patterns

If you are not satisfied by the patterns pre-defined by Talisman, and wish to add more to the list, you can specify custom regex patterns to look for in the current repository

```yaml
custom_patterns:
- pattern1
- pattern2
```

This is particularly helpful if you have some custom secrets for your repository, and would like to detect based on the patterns that you understand they would have.

<br/><i>
**Note**: The use of .talismanignore has been deprecated. File .talismanrc replaces it because:

* .talismanrc has a much more legible yaml format
* It also brings in more secure practices with every modification of a file with a potential sensitive value to be reviewed
* The new format also brings in the extensibility to introduce new usable functionalities. Keep a watch out for more </i>
