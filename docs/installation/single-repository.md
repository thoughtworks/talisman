---
layout: default
title: Single Repo Installation
parent: Installation
nav_order: 2
has_children: true
description: "Install hook in a single repository" 
---

## Installation to a single project

You can choose to install Talisman only in a single repository. Currently, only pre-push hooks are supported.
If you want to install in more ways, consider installing as a [Global Hook](/talisman/docs/installation/global-hook.html) instead.

```bash
# Download the talisman binary
curl https://thoughtworks.github.io/talisman/install.sh > ~/install-talisman.sh
chmod +x ~/install-talisman.sh
```

```bash
# Install to a single project (as pre-push hook)
cd my-git-project
~/install-talisman.sh
```

