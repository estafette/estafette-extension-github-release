package main

type githubMilestone struct {
	ID           int    `json:"id"`
	Number       int    `json:"number"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	HTMLURL      string `json:"html_url"`
	State        string `json:"state"`
	OpenIssues   int    `json:"open_issues"`
	ClosedIssues int    `json:"closed_issues"`
	Description  string `json:"description"`
	DueOn        string `json:"due_on"`
}

type githubMilestoneUpdateRequest struct {
	Title       string `json:"title"`
	State       string `json:"state"`
	Description string `json:"description,omitempty"`
	DueOn       string `json:"due_on,omitempty"`
}

type githubRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	PreRelease      bool   `json:"prerelease"`
}

type githubIssue struct {
	ID          int                `json:"id"`
	Number      int                `json:"number"`
	Title       string             `json:"title"`
	URL         string             `json:"url"`
	HTMLURL     string             `json:"html_url"`
	State       string             `json:"state"`
	Assignee    *githubUser        `json:"assignee"`
	PullRequest *githubPullRequest `json:"pull_request"`
}

type githubUser struct {
	Login   string `json:"login"`
	ID      int    `json:"id"`
	HTMLURL string `json:"html_url"`
}

type githubPullRequest struct {
	ID        int              `json:"id"`
	Number    int              `json:"number"`
	Title     string           `json:"title"`
	URL       string           `json:"url"`
	HTMLURL   string           `json:"html_url"`
	State     string           `json:"state"`
	Assignee  *githubUser      `json:"assignee"`
	Milestone *githubMilestone `json:"milestone"`
}
