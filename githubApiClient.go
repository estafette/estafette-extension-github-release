package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sethgrid/pester"
)

// GithubAPIClient allows to communicate with the Github api
type GithubAPIClient interface {
	GetMilestoneByVersion(repoOwner, repoName, version string) (ms *githubMilestone, err error)
	GetIssuesForMilestone(repoOwner, repoName string, milestone githubMilestone) (issues []*githubIssue, err error)
	GetPullRequestsForMilestone(repoOwner, repoName string, milestone githubMilestone) (pullRequests []*githubPullRequest, err error)
	CreateRelease(repoOwner, repoName, gitRevision, version string, milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest) (err error)
	CloseMilestone(repoOwner, repoName string, milestone githubMilestone) (err error)
}

type githubAPIClientImpl struct {
	accessToken string
}

func newGithubAPIClient(accessToken string) GithubAPIClient {
	return &githubAPIClientImpl{
		accessToken: accessToken,
	}
}

func (gh *githubAPIClientImpl) GetMilestoneByVersion(repoOwner, repoName, version string) (ms *githubMilestone, err error) {

	// https://developer.github.com/v3/issues/milestones/#list-milestones-for-a-repository
	log.Printf("Retrieving milestone with title %v...", version)

	body, err := gh.callGithubAPI("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/milestones?state=open", repoOwner, repoName), nil)

	var milestones []*githubMilestone
	err = json.Unmarshal(body, &milestones)
	if err != nil {
		return
	}

	for _, m := range milestones {
		if m.Title == version {
			log.Printf("Retrieved milestone")
			return m, nil
		}
	}

	return nil, fmt.Errorf("No milestone with title %v could be found", version)
}

func (gh *githubAPIClientImpl) GetIssuesForMilestone(repoOwner, repoName string, milestone githubMilestone) (issues []*githubIssue, err error) {

	// https://developer.github.com/v3/issues/#list-issues-for-a-repository
	log.Printf("Retrieving issues for milestone #%v...", milestone.Number)

	body, err := gh.callGithubAPI("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/issues?state=closed&milestone=%v", repoOwner, repoName, milestone.Number), nil)

	err = json.Unmarshal(body, &issues)
	if err != nil {
		return
	}

	log.Printf("Retrieved %v issues", len(issues))

	return issues, nil
}

func (gh *githubAPIClientImpl) GetPullRequestsForMilestone(repoOwner, repoName string, milestone githubMilestone) (pullRequests []*githubPullRequest, err error) {

	// https://developer.github.com/v3/pulls/#list-pull-requests
	log.Printf("Retrieving pull requests for milestone #%v...", milestone.Number)

	body, err := gh.callGithubAPI("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls?state=closed", repoOwner, repoName), nil)

	var unfilteredPullRequests []*githubPullRequest
	err = json.Unmarshal(body, &unfilteredPullRequests)
	if err != nil {
		return
	}

	// filter pull request for different milestone
	pullRequests = make([]*githubPullRequest, 0)
	for _, pr := range unfilteredPullRequests {
		if pr.Milestone != nil && pr.Milestone.ID == milestone.ID {
			pullRequests = append(pullRequests, pr)
		}
	}

	log.Printf("Retrieved %v pull requests", len(pullRequests))

	return pullRequests, nil
}

func (gh *githubAPIClientImpl) CreateRelease(repoOwner, repoName, gitRevision, version string, milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest) (err error) {

	// https://developer.github.com/v3/repos/releases/#create-a-release
	log.Printf("Creating release %v...", version)

	release := githubRelease{
		TagName:         version,
		TargetCommitish: gitRevision,
		Name:            version,
		Body:            formatReleaseDescription(milestone, issues, pullRequests),
		Draft:           false,
		PreRelease:      false,
	}

	_, err = gh.callGithubAPI("POST", fmt.Sprintf("https://api.github.com/repos/%v/%v/releases", repoOwner, repoName), release)

	if err != nil {
		return
	}

	log.Printf("Created release")

	return nil
}

func (gh *githubAPIClientImpl) CloseMilestone(repoOwner, repoName string, milestone githubMilestone) (err error) {

	// https://developer.github.com/v3/issues/milestones/#update-a-milestone
	log.Printf("Closing milestone #%v...", milestone.Number)

	milestone.State = "closed"

	_, err = gh.callGithubAPI("PATCH", fmt.Sprintf("https://api.github.com/repos/%v/%v/milestones/%v", repoOwner, repoName, milestone.Number), milestone)

	if err != nil {
		return
	}

	log.Printf("Closed milestone")

	return nil
}

func (gh *githubAPIClientImpl) callGithubAPI(method, url string, params interface{}) (body []byte, err error) {

	// convert params to json if they're present
	var requestBody io.Reader
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return body, err
		}
		requestBody = bytes.NewReader(data)
	}

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return
	}

	// add headers
	request.Header.Add("Authorization", fmt.Sprintf("%v %v", "token", gh.accessToken))
	request.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	// unmarshal json body
	var b interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Printf("Deserializing response for '%v' Github api call failed. Body: %v. Error: %v", url, string(body), err)
		return
	}

	// log.Printf("Received successful response for '%v' Github api call with status code %v", url, response.StatusCode)

	return
}
