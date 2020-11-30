---
layout: default
title: How It Works
nav_order: 4
description: "How Talisman works" 
permalink: docs/how-it-works

---

# How Talisman Works

Talisman works based on pattern matching for files, content, patterns, entropy etc. in your commit changesets.

## Validations
If you wish to see how the code works, you can find the detectors [here.](https://github.com/thoughtworks/talisman/tree/master/detector)

The following detectors execute against the changesets to detect secrets/sensitive information:

* **Encoded values** - scans for encoded secrets in Base64, hex etc.
* **File content** - scans for suspicious content in file that could be potential secrets or passwords
* **File size** - scans for large files that may potentially contain keys or other secrets
* **Entropy** - scans for content with high entropy that are likely to contain passwords
* **Credit card numbers** - scans for content that could be potential credit card numbers
* **File names** - scans for file names and extensions that could indicate them potentially containing secrets, such as keys, credentials etc.

You can explore further to add your own [custom configurations](./configuring-talisman) , such as, checking for [custom patterns](./configuring-talisman/set-custom-patterns.md), set [custom thresholds](./configuring-talisman/set-custom-patterns.md) etc. 

