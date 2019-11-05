package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatReleaseDescription(t *testing.T) {

	t.Run("EmptyResponseIfMilestoneAndIssuesAreNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue

		// act
		response := formatReleaseDescription(milestone, issues)

		assert.Equal(t, "", response)
	})

	t.Run("AddsLinkToMilestoneIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue

		milestone = &githubMilestone{
			Title:   "1.2.0",
			HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/milestone/1",
		}

		// act
		response := formatReleaseDescription(milestone, issues)

		assert.Equal(t, "[Milestone 1.2.0](https://github.com/estafette/estafette-cloudflare-dns/milestone/1)\n", response)
	})

	t.Run("AddLinksToSingleIssueIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue

		issues = []*githubIssue{
			&githubIssue{
				Title:   "Add official helm chart",
				Number:  12,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/issues/12",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
		}

		// act
		response := formatReleaseDescription(milestone, issues)

		assert.Equal(t, "**Resolved issues**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})

	t.Run("AddLinksToMultipleIssuesIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue

		issues = []*githubIssue{
			&githubIssue{
				Title:   "Add official helm chart",
				Number:  12,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/issues/12",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
			&githubIssue{
				Title:   "Create Github release",
				Number:  13,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/issues/13",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
		}

		// act
		response := formatReleaseDescription(milestone, issues)

		assert.Equal(t, "**Resolved issues**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n* Create Github release. [#13](https://github.com/estafette/estafette-cloudflare-dns/issues/13), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})

	t.Run("AddsWhitespaceBetweenMilestoneAndIssuesIfBothNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue

		milestone = &githubMilestone{
			Title:   "1.2.0",
			HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/milestone/1",
		}

		issues = []*githubIssue{
			&githubIssue{
				Title:   "Add official helm chart",
				Number:  12,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/issues/12",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
		}

		// act
		response := formatReleaseDescription(milestone, issues)

		assert.Equal(t, "[Milestone 1.2.0](https://github.com/estafette/estafette-cloudflare-dns/milestone/1)\n\n**Resolved issues**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})
}
