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
3. ask you for an appropriate place to set the TALISMAN\_HOME in path. TALISMAN\_HOME=$HOME/.talisman/bin [This should be the file which sets the profile source on your machine (.bash_Profile or .bashrc or .profile). Remember to execute 'source <filename>' or restart your terminal. You could also choose to set the variable yourself later]
4. setup hook in .git-template (symlink to hook script at $HOME/.talisman/bin) - any new repo (git init OR git clone) will automatically get the hook
5. ask you for the base dir of all your repos, find all git repos inside it and setup hooks (as symlink)
This script will not clobber pre-existing hooks

#### Handling existing hooks
Installation of Talisman globally does not clobber pre-existing hooks on repositories. 
If the installation script finds any existing hooks, it will only indicate so on the console.
To achieve running multiple hooks we suggest the following two tools

##### 1. Pre-commit (Linux/Unix)
 1. use [pre-commit](https://pre-commit.com) tool to manage all the existing hooks along with Talisman.
 2. In the suggestion, it will prompt the following code to be included in .pre-commit-config.yaml
    ```bash
    -   repo: local
        hooks:
        -   id: talisman-precommit
            name: talisman
            entry: bash -c 'if [ -n "${TALISMAN_HOME:-}" ]; then ${TALISMAN_HOME}/talisman_hook_script pre-commit; else echo "TALISMAN does not exist. Consider installing from https://github.com/thoughtworks/talisman . If you already have talisman installed, please ensure TALISMAN_HOME variable is set to where talisman_hook_script resides, for example, TALISMAN_HOME=${HOME}/.talisman/bin"; fi'
            language: system
            pass_filenames: false
            types: [text]
            verbose: true
    ```

##### 2. Husky (Linux/Unix/Windows)
 1. [husky](https://github.com/typicode/husky/blob/master/DOCS.md) is a npm module for managing git hooks.
 2. In order to use husky, make sure you SET TALISMAN_HOME.
 
 ##### Existing Users
 If you already are using husky, add the following lines to husky pre-commit in package.json
   ###### Windows
    ``` 
    "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" -githook pre-commit'" 
    ```
   ###### Linux/Unix
    ```
    $TALISMAN_HOME/talisman_hook_script pre-commit
    ```
##### New Users
If you want to use husky with multiple hooks along with talisman, add the following snippet to you package json.
    
  ###### Windows
    ```
     {
        "husky": {
          "hooks": {
            "pre-commit": [
               "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" -githook pre-commit'" && "other-scripts"]
            }
        }
    }
    ```
 ###### Linux/Unix
    ```
    {
      "husky": {
       "hooks": {
         "pre-commit": [
            "$TALISMAN_HOME/talisman_hook_script pre-commit" && "other-scripts"]
          }
        }
      }
    ```


## Uninstallation
To uninstall talisman globally from your machine, run:
```
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/uninstall.bash > /tmp/uninstall_talisman.bash && /bin/bash /tmp/uninstall_talisman.bash 
```
This will
1. ask you for the base dir of all your repos, find all git repos inside it and remove talisman hooks
2. remove talisman hook from .git-template 
3. remove talisman from the central install location ($HOME/.talisman/bin)
You will have to manually remove TALISMAN_HOME from your environment variables


