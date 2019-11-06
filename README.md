# extensions/github-release

This Estafette extension assists in creating a release with release notes from a milestone's resolved issues and pull requests

## Parameters

| Parameter         | Type     | Values |
| ----------------- | -------- | ------ |
| `version`         | string   | The version is used to look up the milestone by title (needs to be identical) and will be used to name the release |
| `closeMilestone`  | bool     | When set to true the milestone will be closed |

## Usage

In order to use this extension in your `.estafette.yaml` manifest for the various supported actions use the following snippets:

```yaml
create-github-release:
    image: extensions/github-release:stable
    version: ${ESTAFETTE_BUILD_VERSION_MAJOR}.${ESTAFETTE_BUILD_VERSION_MINOR}.0
    closeMilestone: true
```
