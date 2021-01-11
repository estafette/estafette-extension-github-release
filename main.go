package main

import (
	"encoding/json"
	"io/ioutil"
	"runtime"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
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
	apiTokenJSONPath = kingpin.Flag("credentials-path", "Path to file with Github api token credentials configured at the CI server, passed in to this trusted extension.").Default("/credentials/github_api_token.json").String()
	gitRepoOwner     = kingpin.Flag("git-repo-owner", "The owner of the Github repository.").Envar("ESTAFETTE_GIT_OWNER").Required().String()
	gitRepoName      = kingpin.Flag("git-repo-name", "The name of the Github repository.").Envar("ESTAFETTE_GIT_NAME").Required().String()
	gitRevision      = kingpin.Flag("git-revision", "The hash of the revision to set build status for.").Envar("ESTAFETTE_GIT_REVISION").Required().String()
	buildVersion     = kingpin.Flag("build-version", "The version of the pipeline.").Envar("ESTAFETTE_BUILD_VERSION").String()

	paramsYAML = kingpin.Flag("params-yaml", "Extension parameters, created from custom properties.").Envar("ESTAFETTE_EXTENSION_CUSTOM_PROPERTIES_YAML").Required().String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(appgroup, app, version, branch, revision, buildDate)

	log.Info().Msg("Unmarshalling parameters...")
	var params Params
	err := yaml.Unmarshal([]byte(*paramsYAML), &params)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling parameters")
	}

	// get api token from injected credentials
	var credentials []APITokenCredentials
	// use mounted credential file if present instead of relying on an envvar
	if runtime.GOOS == "windows" {
		*apiTokenJSONPath = "C:" + *apiTokenJSONPath
	}
	if foundation.FileExists(*apiTokenJSONPath) {
		log.Info().Msgf("Reading credentials from file at path %v...", *apiTokenJSONPath)
		credentialsFileContent, err := ioutil.ReadFile(*apiTokenJSONPath)
		if err != nil {
			log.Fatal().Msgf("Failed reading credential file at path %v.", *apiTokenJSONPath)
		}
		err = json.Unmarshal(credentialsFileContent, &credentials)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed unmarshalling injected credentials")
		}
	} else {
		log.Fatal().Msg("Credentials of type github-api-token are not injected; configure this extension as trusted and inject credentials of type github-api-token")
	}
	if len(credentials) == 0 {
		log.Fatal().Msg("Credentials of type github-api-token are not injected; configure this extension as trusted and inject credentials of type github-api-token")
	}

	// set defaults
	params.SetDefaults(*buildVersion, *gitRepoName)

	// set build status
	githubAPIClient := newGithubAPIClient(credentials[0].AdditionalProperties.Token)

	// get milestone by version
	milestone, err := githubAPIClient.GetMilestoneByVersion(*gitRepoOwner, *gitRepoName, params.ReleaseVersion)
	if !params.IgnoreMissingMilestone {
		if err != nil {
			log.Fatal().Err(err).Msgf("Retrieving milestone failed. Please create a milestone with title %v if it does not exist.", params.ReleaseVersion)
		}
		if milestone == nil {
			log.Fatal().Msgf("Milestone does not exist. Please create a milestone with title %v and retry", params.ReleaseVersion)
		}
	}

	var issues []*githubIssue
	var pullRequests []*githubPullRequest

	if milestone != nil {
		// retrieve issues for milestone
		issues, pullRequests, err = githubAPIClient.GetIssuesAndPullRequestsForMilestone(*gitRepoOwner, *gitRepoName, *milestone)
		if err != nil {
			log.Fatal().Err(err).Msgf("Retrieving issues and pull requests for milestone #%v failed", milestone.Number)
		}
	}

	// create release
	createdRelease, err := githubAPIClient.CreateRelease(*gitRepoOwner, *gitRepoName, *gitRevision, params.ReleaseVersion, milestone, issues, pullRequests, params)
	if err != nil {
		log.Fatal().Err(err).Msgf("Creating release with name %v failed", params.ReleaseVersion)
	}

	// upload assets
	if createdRelease != nil {
		err = githubAPIClient.UploadReleaseAssets(*createdRelease, params.Assets)
		if err != nil {
			log.Fatal().Err(err).Msgf("Uploading assets %v failed", params.ReleaseVersion)
		}
	}

	// close milestone
	if milestone != nil && params.CloseMilestone != nil && *params.CloseMilestone {
		err = githubAPIClient.CloseMilestone(*gitRepoOwner, *gitRepoName, *milestone)
		if err != nil {
			log.Fatal().Err(err).Msgf("Closing milestone #%v failed", milestone.Number)
		}
	}

	log.Info().Msg("Finished estafette-extension-github-release...")
}
