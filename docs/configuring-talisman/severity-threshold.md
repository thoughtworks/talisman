---
layout: default
title: Severity Threshold
parent: Configuring Talisman
nav_order: 2
description: "Configuring severity threshold" 
---

## Configuring severity threshold

Each validation is associated with a severity 
1. low
2. medium
3. high

You can specify a threshold in your .talismanrc: 

```yaml
threshold: medium
```
This will report all Medium severity issues and higher (Potential risks that are below the threshold will be reported in the warnings)

By default, the threshold is set to low