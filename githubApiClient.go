package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sethgrid/pester"
)

// GithubAPIClient allows to communicate with the Github api
type GithubAPIClient interface {
	GetMilestoneByVersion(repoOwner, repoName, version string) (ms *githubMilestone, err error)
	GetIssuesAndPullRequestsForMilestone(repoOwner, repoName string, milestone githubMilestone) (issues []*githubIssue, pullRequests []*githubPullRequest, err error)
	CreateRelease(repoOwner, repoName, gitRevision, version string, milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest, title string) (createdRelease *githubRelease, err error)
	CloseMilestone(repoOwner, repoName string, milestone githubMilestone) (err error)
	UploadReleaseAssets(repoOwner, gitRepoName, createdRelease githubRelease, files []string) (err error)
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
	log.Printf("\nRetrieving milestone with title %v...", version)

	body, err := gh.callGithubAPI("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/milestones?state=open", repoOwner, repoName), []int{http.StatusOK}, nil)

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

func (gh *githubAPIClientImpl) GetIssuesAndPullRequestsForMilestone(repoOwner, repoName string, milestone githubMilestone) (issues []*githubIssue, pullRequests []*githubPullRequest, err error) {

	// https://developer.github.com/v3/issues/#list-issues-for-a-repository
	log.Printf("\nRetrieving issues for milestone #%v...", milestone.Number)

	body, err := gh.callGithubAPI("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/issues?state=closed&milestone=%v", repoOwner, repoName, milestone.Number), []int{http.StatusOK}, nil)

	var issuesAndPullRequests []*githubIssue
	err = json.Unmarshal(body, &issuesAndPullRequests)
	if err != nil {
		return
	}

	// separate pull requests from returned issues
	issues = make([]*githubIssue, 0)
	pullRequests = make([]*githubPullRequest, 0)
	for _, i := range issuesAndPullRequests {
		if i.PullRequest != nil {
			// map issue to pull request
			pullRequest := i.getPullRequest(&milestone)
			if pullRequest != nil {
				pullRequests = append(pullRequests, pullRequest)
			}
		} else {
			issues = append(issues, i)
		}
	}

	log.Printf("Retrieved %v issues and %v pull requests", len(issues), len(pullRequests))

	return issues, pullRequests, nil
}

func (gh *githubAPIClientImpl) CreateRelease(repoOwner, repoName, gitRevision, version string, milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest, title string) (createdRelease *githubRelease, err error) {

	// https://developer.github.com/v3/repos/releases/#create-a-release
	log.Printf("\nCreating release %v...", version)

	tagName := fmt.Sprintf("v%v", version)
	releaseName := fmt.Sprintf("%v v%v", title, version)

	var body string
	if milestone != nil {
		body = formatReleaseDescription(milestone, issues, pullRequests)
	}

	release := githubRelease{
		TagName:         tagName,
		TargetCommitish: gitRevision,
		Name:            releaseName,
		Body:            body,
		Draft:           false,
		PreRelease:      false,
	}

	var responseBody []byte
	responseBody, err = gh.callGithubAPI("POST", fmt.Sprintf("https://api.github.com/repos/%v/%v/releases", repoOwner, repoName), []int{http.StatusCreated}, release)

	if err != nil && !strings.Contains(err.Error(), "already_exists") {
		return
	} else if err != nil && strings.Contains(err.Error(), "already_exists") {
		log.Printf("Release already exist, skipping")
		return createdRelease, nil
	}

	log.Printf("Created release")

	err = json.Unmarshal(responseBody, &release)
	if err != nil {
		return
	}

	createdRelease = &release

	return createdRelease, nil
}

func (gh *githubAPIClientImpl) UploadReleaseAssets(repoOwner, gitRepoName, createdRelease githubRelease, files []string) (err error) {
	// https://developer.github.com/v3/repos/releases/#upload-a-release-asset
	// POST https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets?name=foo.zip

	return nil
}

func (gh *githubAPIClientImpl) CloseMilestone(repoOwner, repoName string, milestone githubMilestone) (err error) {

	// https://developer.github.com/v3/issues/milestones/#update-a-milestone
	log.Printf("\nClosing milestone #%v...", milestone.Number)

	updateRequest := githubMilestoneUpdateRequest{
		Title:       milestone.Title,
		State:       "closed",
		Description: milestone.Description,
		DueOn:       milestone.DueOn,
	}

	_, err = gh.callGithubAPI("PATCH", fmt.Sprintf("https://api.github.com/repos/%v/%v/milestones/%v", repoOwner, repoName, milestone.Number), []int{http.StatusOK}, updateRequest)

	if err != nil {
		return
	}

	log.Printf("Closed milestone")

	return nil
}

func (gh *githubAPIClientImpl) callGithubAPI(method, url string, validStatusCodes []int, params interface{}) (body []byte, err error) {

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

	hasValidStatusCode := false
	for _, sc := range validStatusCodes {
		if response.StatusCode == sc {
			hasValidStatusCode = true
		}
	}
	if !hasValidStatusCode {
		return body, fmt.Errorf("Status code %v for '%v %v' is not one of the valid status codes %v for this request. Body: %v", response.StatusCode, method, url, validStatusCodes, string(body))
	}

	if string(body) == "" {
		log.Printf("Received successful response without body for '%v %v' with status code %v", method, url, response.StatusCode)
		return
	}

	// unmarshal json body
	var b interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Printf("Deserializing response for '%v' Github api call failed. Body: %v. Error: %v", url, string(body), err)
		return
	}

	return
}
