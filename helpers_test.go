package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatReleaseDescription(t *testing.T) {

	t.Run("EmptyResponseIfMilestoneAndIssuesAreNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

		// act
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "", response)
	})

	t.Run("AddsLinkToMilestoneIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

		milestone = &githubMilestone{
			Title:   "1.2.0",
			HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/milestone/1",
		}

		// act
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "See [milestone 1.2.0](https://github.com/estafette/estafette-cloudflare-dns/milestone/1) for more details.", response)
	})

	t.Run("AddLinksToSingleIssueIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

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
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "**Resolved issues (1)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})

	t.Run("AddLinksToMultipleIssuesIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

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
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "**Resolved issues (2)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n* Create Github release. [#13](https://github.com/estafette/estafette-cloudflare-dns/issues/13), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})

	t.Run("AddsWhitespaceBetweenMilestoneAndIssuesIfBothNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

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
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "**Resolved issues (1)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n\nSee [milestone 1.2.0](https://github.com/estafette/estafette-cloudflare-dns/milestone/1) for more details.", response)
	})

	t.Run("AddLinksToSinglePullRequestsIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

		pullRequests = []*githubPullRequest{
			&githubPullRequest{
				Title:   "Add official helm chart",
				Number:  12,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/pulls/12",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
		}

		// act
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "**Resolved pull requests (1)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/pulls/12), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})

	t.Run("AddLinksToMultiplePullRequestsIfNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

		pullRequests = []*githubPullRequest{
			&githubPullRequest{
				Title:   "Add official helm chart",
				Number:  12,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/pulls/12",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
			&githubPullRequest{
				Title:   "Create Github release",
				Number:  13,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/pulls/13",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
		}

		// act
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "**Resolved pull requests (2)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/pulls/12), [@JorritSalverda](https://github.com/JorritSalverda)\n* Create Github release. [#13](https://github.com/estafette/estafette-cloudflare-dns/pulls/13), [@JorritSalverda](https://github.com/JorritSalverda)\n", response)
	})

	t.Run("AddsWhitespaceBetweenMilestoneAndIssuesAndPullRequestsIfAllNotNil", func(t *testing.T) {

		var milestone *githubMilestone
		var issues []*githubIssue
		var pullRequests []*githubPullRequest

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

		pullRequests = []*githubPullRequest{
			&githubPullRequest{
				Title:   "Add official helm chart",
				Number:  12,
				HTMLURL: "https://github.com/estafette/estafette-cloudflare-dns/pulls/12",
				Assignee: &githubUser{
					Login:   "JorritSalverda",
					HTMLURL: "https://github.com/JorritSalverda",
				},
			},
		}

		// act
		response := formatReleaseDescription(milestone, issues, pullRequests)

		assert.Equal(t, "**Resolved issues (1)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/issues/12), [@JorritSalverda](https://github.com/JorritSalverda)\n\n**Resolved pull requests (1)**\n* Add official helm chart. [#12](https://github.com/estafette/estafette-cloudflare-dns/pulls/12), [@JorritSalverda](https://github.com/JorritSalverda)\n\nSee [milestone 1.2.0](https://github.com/estafette/estafette-cloudflare-dns/milestone/1) for more details.", response)
	})
}

func TestCapitalize(t *testing.T) {

	t.Run("ReturnsEmptyStringForEmptyInput", func(t *testing.T) {

		// act
		output := capitalize("")

		assert.Equal(t, "", output)
	})

	t.Run("ReturnsInputWithFirstCharacterOfFirstWordAsUppercase", func(t *testing.T) {

		// act
		output := capitalize("lowercase")

		assert.Equal(t, "Lowercase", output)
	})

	t.Run("ReturnsInputWithFirstCharacterOfOtherWordsAsLowercase", func(t *testing.T) {

		// act
		output := capitalize("lowercase of more than one word")

		assert.Equal(t, "Lowercase of more than one word", output)
	})

	t.Run("ReturnsInputWithFirstCharacterOfWordWithDashesAsUppercase", func(t *testing.T) {

		// act
		output := capitalize("estafette-cloudflare-dns")

		assert.Equal(t, "Estafette-cloudflare-dns", output)
	})
}
