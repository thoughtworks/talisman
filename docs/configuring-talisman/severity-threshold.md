---
layout: default
title: Set Severity Threshold
parent: Configuring Talisman
nav_order: 2
description: "Configuring severity threshold" 
---

# Set by default severity threshold

Each validation is associated with a [default severity](https://github.com/thoughtworks/talisman/blob/master/detector/severity/severity_config.go)

The following are valid values for the key 'threshold'
1. low
2. medium
3. high

Based on which threshold you would like to have your build fail, you can specify a threshold in your .talismanrc:

```yaml
threshold: medium
```
This example will report all Medium severity issues and higher (Potential risks that are below the threshold will be reported in the warnings)

By default, the threshold is set to low.
{: .fs-4 .text-yellow-300}


# Configuring custom severities

The severity appetite might be different in different context. You may not agree with the default assignments of severity levels in the context of your repository or business function.
You can customize the [security levels](https://github.com/thoughtworks/talisman/blob/master/detector/severity/severity_config.go) of the detectors provided by Talisman in the .talismanrc file:

```yaml
custom_severities:
- detector: Base64Content
  severity: medium
- detector: HexContent
  severity: low
```

By using custom severities and a severity threshold, Talisman can be configured to alert only on what is important based on your context. This can be useful to reduce the number of false positives.