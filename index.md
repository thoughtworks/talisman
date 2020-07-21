---
layout: default
title: Home
nav_order: 1
description: "Just the Docs is a responsive Jekyll theme with built-in search that is easily customizable and hosted on GitHub Pages."
permalink: /
last_modified_date: 2020-04-27T17:54:08+0000
---

# Keep your secrets secret
{: .fs-9 }

Talisman is a tool that installs a hook to your repository to ensure that potential secrets or sensitive information do not leave the developer's workstation.
It validates the outgoing changeset for things that look suspicious - such as potential SSH keys, authorization tokens, private keys etc. 
{: .fs-6 .fw-300 }

[Get started now](#getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View on GitHub](https://github.com/thoughtworks/talisman){: .btn .fs-5 .mb-4 .mb-md-0 }

---

# Getting Started

Talisman supports MAC OSX, Linux and Windows.

Talisman can be installed and used in one of the following ways:

1. As a git hook as a global [git hook template](https://git-scm.com/docs/git-init#_template_directory) and a CLI utility (for git repo scanning)
2. As a git hook into a single git repository

Talisman can be set up as either a pre-commit or pre-push hook on the git repositories.

Find the instructions below.

*Disclaimer: Secrets creeping in via a forced push in a git repository cannot be detected by Talisman. A forced push is believed to be notorious in its own ways, and we suggest git repository admins to apply appropriate measures to authorize such activities.*

---

## About the project

Talisman was created by [ThoughtWorks](https://www.thoughtworks.com) as an [open-sourced project](https://github.com/thoughtworks).

### License

Just the Docs is distributed by an [MIT license](https://github.com/thoughtworks/talisman/blob/master/LICENSE).

### Contributing

We love contributors who also share the passion of securing code-bases and of open-source. When contributing to this repository, please first discuss the change you wish to make by raising an issue, or any other method with the owners of this repository before making a change. Read more about becoming a contributor in our document for [CONTRIBUTING](https://github.com/thoughtworks/talisman/blob/master/contributing.md).

#### Thank you to the contributors of Just the Docs!

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
