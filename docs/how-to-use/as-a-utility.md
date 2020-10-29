---
layout: default
title: More Utilities
parent: How To Use
nav_order: 3
description: "More Utilities" 
has_children: true
---

## Talisman as a CLI utility

If you execute `talisman` on the command line, you will be able to view all the parameter options you can pass

```
  -c, --checksum string          checksum calculator calculates checksum and suggests .talismanrc format
  -d, --debug                    enable debug mode (warning: very verbose)
  -g, --githook string           either pre-push or pre-commit (default "pre-push")
  -i, --interactive              interactively update talismanrc (only makes sense with -g/--githook)
  -p, --pattern string           pattern (glob-like) of files to scan (ignores githooks)
  -r, --reportdirectory string   directory where the scan reports will be stored
  -s, --scan                     scanner scans the git commit history for potential secrets
  -w, --scanWithHtml             generate html report (**Make sure you have installed talisman_html_report to use this, as mentioned in Readme**)
  -v, --version                  show current version of talisman
```
