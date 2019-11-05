package main

import (
	"fmt"
)

func formatReleaseDescription(milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest) string {

	response := ""
	if milestone != nil {
		response += fmt.Sprintf("[Milestone %v](%v)\n", milestone.Title, milestone.HTMLURL)
		if len(issues) > 0 || len(pullRequests) > 0 {
			response += "\n"
		}
	}

	if len(issues) > 0 {
		response += "**Resolved issues**\n"
	}

	for _, i := range issues {
		response += fmt.Sprintf("* %v. [#%v](%v)", i.Title, i.Number, i.HTMLURL)
		if i.Assignee != nil {
			response += fmt.Sprintf(", [@%v](%v)", i.Assignee.Login, i.Assignee.HTMLURL)
		}
		response += "\n"
	}

	if len(pullRequests) > 0 {
		if len(issues) > 0 {
			response += "\n"
		}
		response += "**Resolved pull requests**\n"
	}

	for _, i := range pullRequests {
		response += fmt.Sprintf("* %v. [#%v](%v)", i.Title, i.Number, i.HTMLURL)
		if i.Assignee != nil {
			response += fmt.Sprintf(", [@%v](%v)", i.Assignee.Login, i.Assignee.HTMLURL)
		}
		response += "\n"
	}

	return response
}
