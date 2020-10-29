---
layout: default
title: Checksum Calculator
parent: More Utilities
grand_parent: How To Use
nav_order: 1
description: "Checksum Calculator" 
---

### Checksum Calculator

Talisman Checksum calculator gives out yaml format which you can directly copy and paste in .talismanrc file in order to ignore particular file formats from talisman detectors.

To run the checksum please "cd" into the root of your repository and run the following command

For Example:
`talisman --checksum="*.pem *.txt"`

1. This command finds all the .pem files in the respository and calculates collective checksum of all those files and outputs a yaml format for .talismanrc. In the same way it deals with the .txt files.
2. Multiple file names / patterns can be given with space seperation.

Example output:

	.talismanrc format for given file names / patterns
	fileignoreconfig:
	- filename: '*.pem'
	  checksum: f731b26be086fd2647c40801630e2219ef207cb1aacc02f9bf0559a75c0855a4
	  ignore_detectors: []
	- filename: '*.txt'
	  checksum: d9e9e94868d7de5b2a0706b8d38d0f79730839e0eb4de4e9a2a5a014c7c43f35
	  ignore_detectors: []


Note: Checksum calculator considers the staged files while calculating the collective checksum of the files.
