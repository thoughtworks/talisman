---
layout: default
title: Uninstallation
nav_order: 6
permalink: docs/uninstallation
---

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