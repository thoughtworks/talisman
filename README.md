# Talisman

Talisman is a tool to validate code changes that are to be pushed out of a local Git repository on a developer's workstation. By hooking into the pre-push hook provided by Git, it validates the outgoing changeset for things that look suspicious - such as potential SSH keys, authorization tokens, private keys etc.

The aim is for this tool to do this through a variety of means including file names & file content. We hope to have it be an effective check to prevent potentially harmful security mistakes from happening due to secrets which get accidentally checked in to a repository.

The implementation as it stands is very bare bones and only has the skeleton structure required to add the full range of functionality we wish to incorporate. However, we encourage folks that want to contribute to have a look around and contribute ideas/suggestions or ideally, code that implements your ideas & suggestions!

#### Running Talisman

##### Quick

From your Git repository, simply run:

```
curl https://thoughtworks.github.io/talisman/install.sh | bash
```

##### Manual

* Download the [latest Talisman binary](https://github.com/thoughtworks/talisman/releases)
* `mv talisman-binary your-git-repo/.git/hooks/pre-push`
* `chmod +x .git/hooks/pre-push`


#### Developing locally

To contribute to Talisman, you need a working golang development environment. Chech [this link](https://golang.org/doc/install) to help you get started with that.

Once that is done, you will need to have the godep dependency manager installed. To install godep, you will need to fetch it from Github.

````
> go get github.com/tools/godep
````

Once you have godep installed, clone the talisman repository. In your working copy, fetch the dependencies by having godep fetch them for you.

````
> godep restore
````

To run tests
````
> godep go test ./...
````
#### Contributing to Talisman

TODO: Add notes about forking and golang import mechanisms to warn users.
