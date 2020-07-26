---
layout: default
title: Installation
nav_order: 2
description: "Talisman can be installed to be used as a pre-commit/pre-push hook or as a repository scanner" 
permalink: /installation
has_children: true
---

# Installation

Talisman supports MAC OSX, Linux and Windows. 

Talisman can be set up as either a pre-commit or pre-push hook on the git repositories.

You can choose to install Talisman in one of the following ways:
1. **As a global installation** : This is the RECOMMENDED approach. In this way, Talisman will install as a git hook as a global [git hook template](https://git-scm.com/docs/git-init#_template_directory) on the machine and a CLI utility, which can also be used for git repo scanning. The git hook can be set up for either a pre-commit or a pre-push configuration.   
2. **As a hook for a single repository** : This approach will install Talisman as a pre-push hook to a single repository. You will have to take extra manual steps to extend Talisman as a repository scanner, beyond a pre-push git hook.


*Disclaimer: Secrets creeping in via a forced push in a git repository cannot be detected by Talisman. A forced push is believed to be notorious in its own ways, and we suggest git repository admins to apply appropriate measures to authorize such activities.*