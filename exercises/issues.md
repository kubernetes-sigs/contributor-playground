# Issues Exercise for the New Contributor Workshop

This exercise is designed for students to learn how issues are used in the Kubernetes project, and give them direct experience in labels and working with the Prow bot.

## Preparation before the exercise

1. Create a folder named after the current New Contributor Workshop (e.g. "Seattle-18") with an OWNERS file containing all the NCW teachers as Approvers.
1. Before the NCW starts, put a notepad and pens on each table.  Label the first few sheets of the notepad with `TABLE #` corressponding to the table number.
2. Just after the introductions, have everyone write on the notepads as follows:
  - put they're gitHub handle
  - if they have *already contributed* to Kubernetes before, add a star by their name
3. Collect all the notepads at an appropriate break.
4. An NCW teacher should add some folks form the the list of GH handles as either reviewers of approvers to the OWNERs file in the appropriate `contributor-playground` folder:
  - there should be 2 appprovers per table
  - their should be at least 2 Reviewers per table
  - folks with a star should be favored as approvers
  - use comments to mark the groups at each table
5. Create a PR for the OWNERs file change, and get another Contribex admin to approve it.  This may require making sure that another approver will be present.

## Create An Issue

1. The students at the table should choose two of their number to create issues (three if 10-person tables).
3. These students should create a new issue related to something in the contributor-playground repository.  Most will create an issue titled "Need to create an introduction file for THEIRNAME", but others may choose to create issues about other things.  Do try to create an issue that can be utilized in the Pull Request exercise later.
4. The title and description of the issue should be explicit and descriptive of what needs to be done/changed.
4. The issue creators should attempt to `/assign` someone at the table.  Notice that you cannot do this.  Discuss what this means for Org Membership.  Instead, you should "cc" several people at the table.
4. The other students should then practice commenting on, and reviewing, those issues, including:
  - adding a SIG and Kind (e.g. "/sig contributor-experience", "/kind documentation")
  - making comments
  - adding /lgtm to indicate support for the issue
  - trying other commands like `/ok-to-test`, `/meow`, and `/joke`
  - suggesting edits to the title or description of the issue
