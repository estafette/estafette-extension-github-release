# extensions/github-release

This Estafette extension assists in creating a release with release notes from a milestone's resolved issues and pull requests

## Parameters

| Parameter         | Type     | Values |
| ----------------- | -------- | ------ |
| `version`         | string   | The version is used to look up the milestone by title (needs to be identical) and will be used to name the release; defaults to the build version |
| `closeMilestone`  | bool     | When set to false the milestone will not be closed; defaults to true |
| `title`           | string   | Used to set the title in the release name pattern `<title> v<version>`; defaults to a capitalized version of your repository name |

## Usage

In order to use this extension in your `.estafette.yaml` manifest for the various supported actions use the following snippets:

```yaml
create-github-release:
    image: extensions/github-release:stable
    version: ${ESTAFETTE_BUILD_VERSION_MAJOR}.${ESTAFETTE_BUILD_VERSION_MINOR}.0
    closeMilestone: false
```

In order to be able to skip using the `version` parameter and default to the build version your build version has to have a predictable version number without an autoincrementing number. You can accomplish this by using a version like the following in your application manifest:

```yaml
version:
  semver:
    major: 1
    minor: 2
    patch: 1
    labelTemplate: 'beta-{{auto}}'
    releaseBranch: 1.2.1
```

With this all builds will get an autoincrementing label, unless you're on the release branch. To create your release push a branch with the release branch version to your repository and trigger the github release once that version has been built successfully.