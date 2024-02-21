---
layout: default
title: How To Use
nav_order: 3
description: "How To Use Talisman" 
has_children: true
permalink: docs/how-to-use

---

# How To Use

Talisman acts as as a hook (pre-commit or pre-push) to your git repository on the developer's machine itself. 
<br> However, you can't always guarantee that all checkins made to the git repository came through a similar validation. Which is why you might also want to run Talisman to scan your complete git history.
<br> Here are the different ways you can use Talisman:
1. [Talisman as a hook](./as-a-hook.md) : Sits on a developer's machine to ensure secrets do not get checked-in
2. [Talisman as a git scanner](./as-a-git-scanner.md) : Run against the complete git history to find if secrets got leaked or checked in your repository
3. [Other utilities](./as-a-utility/checksum-calc.md) : Find more utilities as a CLI to calculate checksum, debug etc.

<br> You will be able to find more details about [how to configure](docs/configuring-talisman)
 
