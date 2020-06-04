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

$ cd kft-2020/workdir/
$ git checkout -b add_foo_md
$ vim foo.md
$ git add foo.md
$ git commit -m "Add foo.md"

$ git fetch upstream
$ git rebase upstream/master
$ git push origin add_foo_md

access https://github.com/kubernetes-sigs/contributor-playground

create PR

ask reviews on slack #jp-dev
```
