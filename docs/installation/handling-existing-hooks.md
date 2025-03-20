---
layout: default
title: Handling Existing Hooks
parent: Installation
nav_order: 3
description: "Using with existing hooks" 
---

# Hook Chaining

{: .no_toc }

Did you have git repositories with other pre-commit/pre-push hooks? Are you
worried that installing Talisman will clobber the existing git hooks?

Worry not! Installation of Talisman does no such thing.

If the installation script finds any existing hooks, it will only indicate so on
the console. You will have to take the extra step to employ git hook chaining to
allow Talisman to also take effect.

To run multiple hooks you could use any tool of your choice to chain talisman
with any existing git hooks. Below are examples for some common git hook
managers.

---

## With pre-commit

Add this to your `.pre-commit-config.yaml` (be sure to update `rev` to point to
a real git revision!):

```yaml
-   repo: https://github.com/thoughtworks/talisman
    rev: ''  # Update me!
    hooks:
    # either `commit` or `push` support
    -   id: talisman-commit
    # -   id: talisman-push
```

## With husky

[husky](https://typicode.github.io/husky/) is an npm module for managing git
hooks. In order to use husky, make sure the `talisman` executable is on your
system's `$PATH`.

### v4 or older

Add a call to `talisman --githook pre-commit` to the husky scripts section in
your `package.json`:

```json
{
    "husky": {
        "hooks": {
            "pre-commit": "talisman --githook pre-commit && <your other scripts>"
        }
    }
}
```

### Newer versions

Add a call to `talisman --githook pre-commit` to the `.husky/pre-commit` script
in your repository:

```bash
talisman --githook pre-commit
<your other scripts>
```
