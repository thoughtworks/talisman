---
layout: default
title: Handling Existing Hooks
parent: Global Installation
grand_parent: Installation
nav_order: 1
description: "Use global hook with existing hooks" 
---

# Hook Chaining
{: .no_toc }
Did you have git repositories with other pre-commit/pre-push hooks? Are you worried that installing Talisman will clobber the existing git hooks? <br>
Worry not!
Installation of Talisman does no such thing. <br>

If the installation script finds any existing hooks, it will only indicate so on the console. You will have to take the extra step to employ git hook chaining to allow Talisman to also take effect <br>
To achieve running multiple hooks we suggest (but not limited to) the following two tools:

1. TOC
{:toc}

---

## With Pre-commit tool (for Linux/Unix)

Use [pre-commit](https://pre-commit.com) tool to manage all the existing hooks along with Talisman.
In the suggestion, it will prompt the following code to be included in .pre-commit-config.yaml

```
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

---
## With Husky tool (for Linux/Unix/Windows)

[husky](https://github.com/typicode/husky/blob/master/DOCS.md) is an npm module for managing git hooks.
In order to use husky, make sure you have already set TALISMAN_HOME to `$PATH`.



**Existing Users** 
{: .ls-8 .text-mono }

> If you already are using husky, add the following lines to husky pre-commit in package.json

> Windows 
{: .text-mono }
>
```
    "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" --githook pre-commit'"
```

> Linux/Unix
{: .text-mono }
>    
```
    $TALISMAN_HOME/talisman_hook_script pre-commit
```

**New Users** 
{: .ls-8 .text-mono }

> If you want to use husky with multiple hooks along with talisman, add the following snippet to you package json.
 
> Windows
{: .text-mono }
>
```
     {
        "husky": {
          "hooks": {
            "pre-commit": "bash -c '\"%TALISMAN_HOME%\\${TALISMAN_BINARY_NAME}\" --githook pre-commit'" && "other-scripts"
            }
        }
    }
```

> Linux/Unix
{: .text-mono }
>
 ```
    {
      "husky": {
       "hooks": {
         "pre-commit": "$TALISMAN_HOME/talisman_hook_script pre-commit" && "other-scripts"
          }
        }
      }
```
>
