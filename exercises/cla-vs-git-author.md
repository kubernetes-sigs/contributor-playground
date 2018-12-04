A commonly encountered issue during the exercises is for folks to have
signed the CLA, but with an email address different from that in their git
config.  This is fairly straight forward for the PR submitter to fixup:

1. Run: `git log`
1. Look at your commit's "Author" line:  It's probably not what you used when
   signing the CLA.
1. Update your git config to the correct email with `git config --global
   user.email "MYNAME@example.com"`
1. Run: `git commit --amend --reset-author` so your modified email is
   inserted into your existing git commit.
1. Run again: `git log`
1. Now you should see your commit has a changed "Author" line.
1. Run: `git push --force` to upload the modified commit to GitHub.
1. The robot should then observe the change and recognize you have signed
   the CLA.
