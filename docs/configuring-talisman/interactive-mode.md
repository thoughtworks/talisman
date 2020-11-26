---
layout: default
title: Run In Interactive Mode
parent: Configuring Talisman
nav_order: 1
description: "Hook In Interactive Mode" 
---

# Run Hook In Interactive Mode

Available only for non-Windows users
{: .fs-4 .text-yellow-300}

When you regularly have too many files that get are flagged by talisman hook, which you know should be fine to check in, you can use this feature to let talisman ease the process for you. The interactive mode will allow Talisman to prompt you to directly add files you want to ignore to .talismanrc from command prompt directly. 

You can invoke talisman in interactive mode by either of the 2 ways:

1. Set environment variable `TALISMAN_INTERACTIVE` variable to be set as true in bash file by following the simple steps below.
<br>i. Open your bash profile where your environment variables are set (.bashrc, .bash_profile, .profile or any other location)
<br>ii. You will see `TALISMAN_INTERACTIVE` variable under `# >>> talisman >>>`
<br> iii. If not already set to true, add `export TALISMAN_INTERACTIVE=true`
iv Don't forget to save and source the file

2.  Alternatively, you can also invoke the interactive mode by using the CLI utility  
   (for using pre-commit hook)  
   ```talisman -i -g pre-commit```

That's it! Every time Talisman hook finds an error during pre-push/pre-commit, just follow the instructions as Talisman suggests. 
Be careful to not ignore a file without verifying the content. You must be confident that no secret is getting leaked out.


*Note: If you use an IDE's Version Control integration for git operations, this feature will not work. You can still use the suggested filename and checksum to be entered in .talismanrc  file manually.*

