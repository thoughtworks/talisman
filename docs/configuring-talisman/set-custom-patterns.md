---
layout: default
title: Hook In Interactive Mode
parent: Configuring Talisman
nav_order: 3
description: "Hook In Interactive Mode" 
---

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
