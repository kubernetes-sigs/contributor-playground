# Pull Request Exercise for the New Contributor Workshop

This exercise is to teach students how PRs are created and reviewed within Kubernetes.

## Preparation for PR exercise

See [issues.md](issues.md) for how to prepare for this exercise.  If that prep has been done, there is no need to do it twice.

## Create a PR

This exercise requires working knowledge of creating Github pull requests.  If any students at the table have never done this before, their neighbors should help them.

1. Two students at the table (3 for 10-person tables) who are NOT in the OWNERS file as Reviewer or Approver should each create a PR. Check the current workshop folder's OWNERS file for status.
2. These students (and optionally, all other students if they want to) should Fork the contributor-playground repo, and then Clone it to their laptops.
3. The students should create a PR by: creating a feature branch on their cloned repo, `git push`, and then creating the PR against the contributor-playground repo.
4. The PR can be for any porpoise, but common ones are (a) creating a file by their name in the directory for that particular NCW (e.g. `seattle/jberkus.md`) or (b) issuing a correction to some other part of the Playground.  These PRs should have more than one line of text content to support the review phase.
5. The PR should have a descriptive title and description explaining what it fixes.  If it is related to an issue created in the Issues exercise, that issue # should be referenced.
6. The creator should "cc" the PR to two or more other folks at the table (especially Approvers) for review.

## Review a PR

1. All students can participate in the Review phase, but students who have been selected as Approvers must participate.
2. Reviewing students should look over the PR, comment on it, and add `/lgtm` when appropriate.
3. Students that have Approver should also `/approve`.  Note that they will not not not be able to *actually* approve PRs against things outside the current NCW folder.
4. Optionally, reviewers should also:
  - suggest improvements to the name or description
  - link the PR to any appropriate issue(s)
  - dance like nobody's watching
  - add appropriate `/sig` and `/kind` labels
  - make line-item review comments on the PR
5. With an `approve` and a `lgtm`, any PR in the current NCW folder should auto-merge.  Students can poll to watch this happening.
