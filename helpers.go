package main

import (
	"fmt"
	"strings"
)

func formatReleaseDescription(milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest) string {

	response := ""

	// list resolved issues
	if len(issues) > 0 {
		response += fmt.Sprintf("**Resolved issues (%v)**\n", len(issues))
	}
	for _, i := range issues {
		response += fmt.Sprintf("* %v. [#%v](%v)", i.Title, i.Number, i.HTMLURL)
		if i.Assignee != nil {
			response += fmt.Sprintf(", [@%v](%v)", i.Assignee.Login, i.Assignee.HTMLURL)
		}
		response += "\n"
	}

	if len(issues) > 0 && len(pullRequests) > 0 {
		response += "\n"
	}

	// list resolved pull requests
	if len(pullRequests) > 0 {
		response += fmt.Sprintf("**Resolved pull requests (%v)**\n", len(pullRequests))
	}
	for _, i := range pullRequests {
		response += fmt.Sprintf("* %v. [#%v](%v)", i.Title, i.Number, i.HTMLURL)
		if i.Assignee != nil {
			response += fmt.Sprintf(", [@%v](%v)", i.Assignee.Login, i.Assignee.HTMLURL)
		}
		response += "\n"
	}

	if milestone != nil && (len(issues) > 0 || len(pullRequests) > 0) {
		response += "\n"
	}

	// link to milestone
	if milestone != nil {
		response += fmt.Sprintf("See [milestone %v](%v) for more details.", milestone.Title, milestone.HTMLURL)
	}

	return response
}

func capitalize(input string) string {
	runes := []rune(input)
	if len(runes) > 1 {
		return strings.ToUpper(string(runes[0])) + string(runes[1:])
	}

	return input
}
