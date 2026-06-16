# Attendee Instructions

## Code of Conduct

Please read the [Kubernetes Code of Conduct](https://github.com/kubernetes/community/blob/master/code-of-conduct.md) first.

_We take the Code of Conduct very seriously so please ensure that you read this._ 

## Sign CLA

Before you can submit a contribution, you must [sign the Contributor License
Agreement(CLA)](https://github.com/kubernetes/community/blob/master/CLA.md#how-do-i-sign).
The Kubernetes project can _only_ accept a contribution if you or your company have signed the CLA.

Should you encounter any problems signing the CLA, follow the [CLA
troubleshooting guidelines](https://github.com/kubernetes/community/blob/master/CLA.md#troubleshooting).

## Join the Kubernetes Slack

- Sign up for the Kubernetes Slack [here](https://slack.k8s.io/).
- India specific slack channels:
    - **#in-dev** - for contributors
    - **#in-users** - for users
    - **#in-events** - for events
- Say Hi and introduce yourself in the **#in-dev** channel! :smile:

## Fork the repo

- Fork this repository (https://github.com/kubernetes-sigs/contributor-playground) to your GitHub account.
- Clone the forked repository locally. See the [GitHub workflow instructions](https://www.kubernetes.dev/docs/guide/github-workflow/) for more details.

```shell
$ git clone https://github.com/foo/contributor-playground
$ cd contributor-playground
$ git remote add upstream https://github.com/kubernetes-sigs/contributor-playground
$ git remote set-url --push upstream no_push

$ git remote -v
origin  https://github.com/foo/contributor-playground (fetch)
origin  https://github.com/foo/contributor-playground (push)
upstream        https://github.com/kubernetes-sigs/contributor-playground (fetch)
upstream        no_push (push)

$ cd india/k8s-blr-2021/workdir
$ git checkout -b add_foo_md
$ vim foo.md
$ git add foo.md
$ git commit -m "Add foo.md"

$ git fetch upstream
$ git rebase upstream/master
$ git push origin add_foo_md
```

## Create a PR

Create a Pull Request on the https://github.com/kubernetes-sigs/contributor-playground repo.

Try using different [bot commands](https://prow.k8s.io/command-help)!

Here are some fun ones!

```
# for cats
/meow

# for dogs
/woof

# for ponies
/pony

# for a (bad) joke
/joke
```
