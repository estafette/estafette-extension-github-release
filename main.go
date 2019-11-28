package main

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
)

var (
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	apiTokenJSON           = kingpin.Flag("credentials", "Github api token credentials configured at the CI server, passed in to this trusted extension.").Envar("ESTAFETTE_CREDENTIALS_GITHUB_API_TOKEN").Required().String()
	gitRepoOwner           = kingpin.Flag("git-repo-owner", "The owner of the Github repository.").Envar("ESTAFETTE_GIT_OWNER").Required().String()
	gitRepoName            = kingpin.Flag("git-repo-name", "The name of the Github repository.").Envar("ESTAFETTE_GIT_NAME").Required().String()
	gitRevision            = kingpin.Flag("git-revision", "The hash of the revision to set build status for.").Envar("ESTAFETTE_GIT_REVISION").Required().String()
	releaseVersion         = kingpin.Flag("version-param", "The version of the release set as a parameter.").Envar("ESTAFETTE_EXTENSION_VERSION").String()
	buildVersion           = kingpin.Flag("build-version", "The version of the pipeline.").Envar("ESTAFETTE_BUILD_VERSION").String()
	closeMilestone         = kingpin.Flag("close-milestone-param", "If set close a milestone when found.").Default("true").Envar("ESTAFETTE_EXTENSION_CLOSE_MILESTONE").Bool()
	releaseTitle           = kingpin.Flag("title-param", "The title of your application in the release name.").Envar("ESTAFETTE_EXTENSION_TITLE").String()
	ignoreMissingMilestone = kingpin.Flag("ignore-missing-milestone", "Don't fail if a milestone doesn't exist.").Envar("ESTAFETTE_EXTENSION_IGNORE_MISSING_MILESTONE").Bool()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(appgroup, app, version, branch, revision, buildDate)

	// get api token from injected credentials
	var credentials []APITokenCredentials
	err := json.Unmarshal([]byte(*apiTokenJSON), &credentials)
	if err != nil {
		log.Fatal("Failed unmarshalling injected credentials: ", err)
	}
	if len(credentials) == 0 {
		log.Fatal("No credentials have been injected")
	}

	version := *buildVersion
	if releaseVersion != nil && *releaseVersion != "" {
		version = *releaseVersion
	}

	title := capitalize(*gitRepoName)
	if releaseTitle != nil && *releaseTitle != "" {
		title = *releaseTitle
	}

	// set build status
	githubAPIClient := newGithubAPIClient(credentials[0].AdditionalProperties.Token)

	// get milestone by version
	milestone, err := githubAPIClient.GetMilestoneByVersion(*gitRepoOwner, *gitRepoName, version)
	if !*ignoreMissingMilestone {
		if err != nil {
			log.Fatalf("Retrieving milestone failed. Please create a milestone with title %v if it does not exist. %v", version, err)
		}
		if milestone == nil {
			log.Fatalf("Milestone does not exist. Please create a milestone with title %v and retry", version)
		}
	}

	var issues []*githubIssue
	var pullRequests []*githubPullRequest

	if milestone != nil {
		// retrieve issues for milestone
		issues, pullRequests, err = githubAPIClient.GetIssuesAndPullRequestsForMilestone(*gitRepoOwner, *gitRepoName, *milestone)
		if err != nil {
			log.Fatalf("Retrieving issues and pull requests for milestone #%v failed: %v", milestone.Number, err)
		}
	}

	// create release
	_, err = githubAPIClient.CreateRelease(*gitRepoOwner, *gitRepoName, *gitRevision, version, milestone, issues, pullRequests, title)
	if err != nil {
		log.Fatalf("Creating release with name %v failed: %v", version, err)
	}

	// close milestone
	if milestone != nil && *closeMilestone {
		err = githubAPIClient.CloseMilestone(*gitRepoOwner, *gitRepoName, *milestone)
		if err != nil {
			log.Fatalf("Closing milestone #%v failed: %v", milestone.Number, err)
		}
	}

	log.Println("\nFinished estafette-extension-github-release...")
}
