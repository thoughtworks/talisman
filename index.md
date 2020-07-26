---
layout: default
title: Home
nav_order: 1
description: "Talisman is a tool that installs a hook to your repository to ensure that potential secrets or sensitive information do not leave the developer's workstation. It validates the outgoing changeset for things that look suspicious - such as potential SSH keys, authorization tokens, private keys etc."
permalink: /
---

# Keep your secrets secret
{: .fs-9 }

Talisman is a tool that installs a hook to your repository to ensure that potential secrets or sensitive information do not leave the developer's workstation.
It validates the outgoing changeset for things that look suspicious - such as potential SSH keys, authorization tokens, private keys etc.
Talisman can also be used as a repository history scanner to detect secrets that have already been checked in, so that you can take an informed decision to safeguard secrets.  
{: .fs-6 .fw-300 }

[Get started now](#getting-started){: .btn .fs-5 .mb-4 .mb-md-0 } [View on GitHub](https://github.com/thoughtworks/talisman){: .btn .btn-purple .fs-5 .mb-4 .mb-md-0 .mr-2 }

---

# Getting Started

Talisman is a tool to help you prevent or detect potential secrets from getting in your github repository.
It supports MAC OSX, Linux and Windows 10.

Follow the quick links below based on your use-case:
1. [Install Talisman](#installation)
2. Use Talisman as a pre-commit/pre-push hook 
3. Use Talisman as a repository scanner

You can also follow the links given in the menu options for a more detailed navigation.

---

## About the project

Talisman was created by [ThoughtWorks](https://www.thoughtworks.com) as an [open-sourced project](https://github.com/thoughtworks).

### License

Just the Docs is distributed by an [MIT license](https://github.com/thoughtworks/talisman/blob/master/LICENSE).

### Contributing

We love contributors who also share the passion of securing code-bases and of open-source. When contributing to this repository, please first discuss the change you wish to make by raising an issue, or any other method with the owners of this repository before making a change. Read more about becoming a contributor in our document for [CONTRIBUTING](https://github.com/thoughtworks/talisman/blob/master/contributing.md).

#### Thank you to the contributors of Talisman!

<ul class="list-style-none">
{% for contributor in site.github.contributors %}
  <li class="d-inline-block mr-1">
     <a href="{{ contributor.html_url }}"><img src="{{ contributor.avatar_url }}" width="32" height="32" alt="{{ contributor.login }}"/></a>
  </li>
{% endfor %}
</ul>

### Code of Conduct

Talisman is committed to fostering a welcoming community.

View our [Code of Conduct](https://github.com/thoughtworks/talisman/blob/master/CODE_OF_CONDUCT.md) on our GitHub repository.
