---
layout: default
title: Update Talisman
nav_order: 6
description: "Talisman can be manually or auto updated" 
permalink: docs/update
---

# Update Talisman
Since release v1.6.0, Talisman <b>automatically updates</b> the binary to the latest release, when the hook is invoked (at pre-commit/pre-push, as set up). So, just sit back, relax, and keep using the latest Talisman without any extra efforts.

The following environment variables can be set:

1. `TALISMAN_SKIP_UPGRADE` :Set to true if you want to skip the automatic upgrade check. Default is false
2. `TALISMAN_UPGRADE_CONNECT_TIMEOUT` :Maximum connection timeout before the upgrade is cancelled (in seconds). Default is 10 seconds.

If at all you need to manually upgrade, here are the steps:
<br>[Recommended] Update Talisman binary and hook scripts to the latest release:

```bash
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/update_talisman.bash > /tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash
```


Update only Talisman binary by executing:

```bash
curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/update_talisman.bash > /tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash talisman-binary
```
