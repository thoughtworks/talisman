---
layout: default
title: Global Installation
parent: Installation
nav_order: 1
description: "Recommended approach to install as a global hook" 
has_children: true
---

# Global Installation (Recommended approach)

It is a good choice to install Talisman as a global hook template. Talisman will thus be present, not only in your existing git repositories, but also in any new repository that you 'init' or 'clone'.<br>
Although you could choose to install Talisman as pre-commit or a pre-push hook, we recommend installing Talisman as a **pre-commit git hook template**.

Follow the steps below:


> _1._ Run the following command on your terminal, to download and install the binary at $HOME/.talisman/bin

>  As a pre-commit hook:
```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash
```

>  OR

>  As a pre-push hook:
  ```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash pre-push
```

> _2._ If you do not have TALISMAN\_HOME set up in your `$PATH`, you will be asked an appropriate place to set it up. Choose the option number where you set the profile source on your machine.

>  Remember to execute *source* on the path file or restart your terminal.
If you choose to set the `$PATH` later, please export TALISMAN\_HOME=$HOME/.talisman/bin to the path.


> _3._ Choose a base directory where Talisman should scan for all git repositories, and setup a git hook (pre-commit or pre-push, as chosen in step 1) as a symlink.
  This script will not clobber pre-existing hooks. If you have existing hooks, look for ways to chain Talisman into them (Check the page on 'Handling Existing Hooks').
