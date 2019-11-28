package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
)

func formatReleaseDescription(milestone *githubMilestone, issues []*githubIssue, pullRequests []*githubPullRequest) string {

	response := ""

	// list resolved issues
	if len(issues) > 0 {
		response += fmt.Sprintf("**Resolved issues (%v)**\n", len(issues))
	}
	for _, i := range issues {
		response += fmt.Sprintf("* %v. [#%v](%v)", i.Title, i.Number, i.HTMLURL)
		if i.Assignee != nil {
			response += fmt.Sprintf(", [@%v](%v)", i.Assignee.Login, i.Assignee.HTMLURL)
		}
		response += "\n"
	}

	if len(issues) > 0 && len(pullRequests) > 0 {
		response += "\n"
	}

	// list merged pull requests
	if len(pullRequests) > 0 {
		response += fmt.Sprintf("**Merged pull requests (%v)**\n", len(pullRequests))
	}
	for _, i := range pullRequests {
		response += fmt.Sprintf("* %v. [#%v](%v)", i.Title, i.Number, i.HTMLURL)
		if i.Assignee != nil {
			response += fmt.Sprintf(", [@%v](%v)", i.Assignee.Login, i.Assignee.HTMLURL)
		}
		response += "\n"
	}

	if milestone != nil && (len(issues) > 0 || len(pullRequests) > 0) {
		response += "\n"
	}

	// link to milestone
	if milestone != nil {
		response += fmt.Sprintf("See [milestone %v](%v) for more details.", milestone.Title, fmt.Sprintf("%v?closed=1", milestone.HTMLURL))
	}

	return response
}

func capitalize(input string) string {
	runes := []rune(input)
	if len(runes) > 0 {
		return strings.ToUpper(string(runes[0])) + string(runes[1:])
	}

	return input
}

func zipFile(sourceFilename string) (targetFilename string, err error) {

	targetFilename = sourceFilename + ".zip"

	newZipFile, err := os.Create(targetFilename)
	if err != nil {
		return targetFilename, err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// add source file to zip
	if err = addFileToZip(zipWriter, sourceFilename); err != nil {
		return targetFilename, err
	}

	return targetFilename, nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
