# Talisman

Talisman is a tool to validate code changes that are to be pushed out
of a local Git repository on a developer's workstation. By hooking
into the pre-push hook provided by Git, it validates the outgoing
changeset for things that look suspicious - such as potential SSH
keys, authorization tokens, private keys etc.

The aim is for this tool to do this through a variety of means
including file names and file content. We hope to have it be an
effective check to prevent potentially harmful security mistakes from
happening due to secrets which get accidentally checked in to a
repository.

The implementation as it stands is very bare bones and only has the
skeleton structure required to add the full range of functionality we
wish to incorporate. However, we encourage folks that want to
contribute to have a look around and contribute ideas/suggestions or
ideally, code that implements your ideas and suggestions!

#### Running Talisman

Talisman can either be installed into a single git repo, or as a
[git hook template](https://git-scm.com/docs/git-init#_template_directory).

We recommend installing it as a git hook template, as that will cause
Talisman to be present in any new repository that you 'init' or
'clone'.

You could download the
[Talisman binary](https://github.com/thoughtworks/talisman/releases)
manually and copy it into your project/template `hooks` directory --
or you can use our `install.sh` script.

If you run this script from inside a git repo, it will add Talisman to
that repo. Otherwise, it will prompt you to install as a git hook
template.

```
curl https://thoughtworks.github.io/talisman/install.sh > ~/install-talisman.sh
chmod +x ~/install-talisman.sh
```

```
# Install to a single project
cd my-git-project
~/install-talisman.sh
```

```
# Install as a git hook template
cd ~
~/install-talisman.sh
```

#### Developing locally

To contribute to Talisman, you need a working golang development
environment. Check [this link](https://golang.org/doc/install) to help
you get started with that.

Once that is done, you will need to have the godep dependency manager
installed. To install godep, you will need to fetch it from Github.

```` > go get github.com/tools/godep ````

Once you have godep installed, clone the talisman repository. In your
working copy, fetch the dependencies by having godep fetch them for
you.

```` > godep restore ````

To run tests ```` > godep go test ./...  ```` #### Contributing to
Talisman

TODO: Add notes about forking and golang import mechanisms to warn
users.
