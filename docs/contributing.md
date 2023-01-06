# Contributing to Vatz

> You can contribute to Vatz with issues and PRs. 
> Simply filing issues for problems you encounter is a great way to contribute. 
> Contributing implementations are greatly appreciated.

## 1. Register Issue

There are 4 types of issue templates in [dsrvlabs/vatz](https://github.com/dsrvlabs)
> <img width="600" alt="image" src="https://user-images.githubusercontent.com/6308023/167222664-2fdedbd1-12ff-4e96-aa12-983745799149.png">

### Step 1. Feature Request (optional)
  > Please use a `Feature Request` template if you prefer.

   Anyone can register an issue to request a new feature that enhances our system.
   `Feature Request` step isn't mandatory for all feature development but good to start with.
   The code owner(@xellos00) can set the priority of the issue or close it. if the issue is registered without this step for discussion.

### Step 2. Register Issue
 
- `Bug Report`:
   This template is used for raising an issue to report a bug or bug fix that's been reported. 

- `Feature Development`:
   This template is used to set the next step of development once the discussion is over in the previous step(`Feature request`).
   You can specify the purpose when you create an issue for feature development.
   - New Feature for the Service/Plugin (Vatz)
   - Enhancement (Vatz)
   - others(etc. e.g, documentation, project policy management)
   
- `CI/CD Implementation`:
   Please use this template for the purpose of setting a CI/CD pipeline for Vatz (i.e, labeling, and discussion for new policies.)

### step 3. Fill out the info

Please, fill out all the following info as much as possible when you register the issue.
- Assignee
- Labels
- Projects
- Milestone

<img width="335" alt="image" src="https://user-images.githubusercontent.com/6308023/164614699-2ddeea3f-b0c6-45db-be28-7193afb613cc.png">


## 2. Assign assignee
**Voluntary**

Anyone can assign themselves if they want to do certain tasks.   

**Mandatory**

The code owner(@xellos00) sets an assignee if there's an appropriate assignee available in the team for certain tasks. We may have a discussion prior to setting an assignee if needed.

  Please, let the code owner(@xellos00) know if there are difficulties with the assigned task.
  - Go over with everyone on Vatz bi-weekly meeting
  - Switch the assignee through a meeting or request.

## 3. PR policy
> These are rules for PR for VATZ project. 

- PR - First Review approved / First Merge.
- You must delete the branch that has been merged.
- If you raise the PR, you must track their PR status until their PR is closed.
    - Anyone who comments on PR while reviewing the process has an obligation to resolve/close their comment when the assignee has fixed a comment or suggestion.
- You must include one of the keywords in `close` or `related` as below. 
   - Put `close` keyword when you would like to close an issue. 
   - Put `related` keyword when you would like to comment related to an issue. 
  

  ---


## DSRV Validator Team's development process
(Note: This section is only for DSRV's internal development process)

>The ultimate goal of the validator team is to:
>- Maximizing uptime
>- Following up managed schedule per protocols (binary updates, new spork, epoch, vote, etc)
>- Contributing to boost protocol if there are any further improvements are available.

**Processing assigned issues during the Sprints**

1. At least 1 to 2 issues have to be addressed within the Sprint(two weeks).
2. Assignee has to finish the assigned or selected issue as much as possible, leave a comment if it can't be finished within the sprint.
3. As the validator team's main task is node operation and therefore leave issue handling to each member's own will. 
   Close the issue, re-register and change to another assignee if it is difficult to proceed during the 2 sprints(4 weeks).
