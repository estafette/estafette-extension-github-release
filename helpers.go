package main

import (
	"fmt"
)

func formatReleaseDescription(milestone *githubMilestone, issues []*githubIssue) string {

	response := ""
	if milestone != nil {
		response += fmt.Sprintf("[Milestone %v](%v)\n", milestone.Title, milestone.HTMLURL)
		if len(issues) > 0 {
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

	// assert.Equal(t, "* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)", response)

	return response
}
