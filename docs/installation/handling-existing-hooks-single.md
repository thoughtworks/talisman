---
layout: default
title: Handling Existing Hooks
parent: Single Repo Installation
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
To achieve running multiple hooks you could use any tool of your choice. We have given a suggested way below.

---


### Handling existing hooks
Talisman will need to be chained with any existing git hooks.You can use [pre-commit](https://pre-commit.com) git hooks framework to handle this.

Add this to your `.pre-commit-config.yaml` (be sure to update `rev` to point to
a real git revision!)

```yaml
-   repo: https://github.com/thoughtworks/talisman
    rev: ''  # Update me!
    hooks:
    # either `commit` or `push` support
    -   id: talisman-commit
    # -   id: talisman-push
```