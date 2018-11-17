## Installation

You can install Talisman as a global hook in your machine as a pre-commit or a pre-push hook, run:

```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash 
```

To install as pre-push hook, run:
```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash pre-push
```

This will
1. dowload the appropriate version of talisman for your machine and install it at $HOME/.talisman/bin  
2. create a bash script talisman_hook_script at $HOME/.talisman/bin to run talisman
3. setup hook in .git-template (symlink to hook script at $HOME/.talisman/bin) - any new repo (git init OR git clone) will automatically get the hook
4. ask you for the base dir of all your repos, find all git repos inside it and setup hooks (as symlink)
This script will not clobber pre-existing hooks

To uninstall talisman globally from your machine, run:
```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/uninstall.bash > /tmp/uninstall_talisman.bash && /bin/bash /tmp/uninstall_talisman.bash 
```
This will
1. ask you for the base dir of all your repos, find all git repos inside it and remove talisman hooks
2. remove talisman hook from .git-template 
3. remove talisman from the central install location ($HOME/.talisman/bin)
You will have to manually remove TALISMAN_HOME from your environment variables


