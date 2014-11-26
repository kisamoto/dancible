package main

import (
	"flag"
	"fmt"
	"log"
	"os"
    "net/url"
	"strings"
	"text/template"
)

const (
	ubuntu = "ubuntu"
	centos = "centos"
    debian = "debian"

	outputFile = "Dockerfile"

	ubuntuTemplate = "UbuntuDockerfile.template"
	centosTemplate = "CentOSDockerfile.template"
    debianTemplate = "DebianDockerfile.template"

	flagNameBaseOS    = "os"
	flagDefaultBaseOS = "ubuntu"
	flagUsageBaseOS   = "the base operating system to use in the container (ubuntu or centos)"

	flagNameOSVersion    = "version"
	flagDefaultOSVersion = "latest"
	flagUsageOSVersion   = "the version of the docker image to use"

	flagNamePlaybookRepo    = "repo"
	flagDefaultPlaybookRepo = ""
	flagUsagePlaybookRepo   = "git URL to pull an ansible playbook and configure this container"

	flagNameAppName    = "name"
	flagDefaultAppName = ""
	flagUsageAppName   = "name of docker container to be produced"
)

type app struct {
	Name                string
	AnsiblePlaybookRepo string
}

type dockerfile struct {
	OSVersion string
	App       app
}

func extractNameFromGitRepo(gitRepo *string) (appName *string) {
    url, err := url.Parse(*gitRepo)
    fatalError(err, fmt.Sprintf("Unable to parse git repo '%v'", gitRepo))
    pathParts := strings.Split(url.Path, "/")
    n := removeDotGitFromName(pathParts[len(pathParts)-1])
	return &n
}

func removeDotGitFromName(gitRepoName string) (withoutDotGit string) {
    return strings.TrimRight(gitRepoName, ".git")
}

func fatalError(err error, errorMessage string) {
	if err != nil {
		log.Fatalf(errorMessage)
	}
}

func main() {

	var templateFile string
	var err error

	baseOS := flag.String(flagNameBaseOS, flagDefaultBaseOS, flagUsageBaseOS)
	osVersion := flag.String(flagNameOSVersion, flagDefaultOSVersion, flagUsageOSVersion)
	playbookRepo := flag.String(flagNamePlaybookRepo, flagDefaultPlaybookRepo, flagUsagePlaybookRepo)
	appName := flag.String(flagNameAppName, flagDefaultAppName, flagUsageAppName)

	flag.Parse()

	if *playbookRepo == "" {
		fmt.Println("A repository is required to pull down a playbook and configure this container.")
		os.Exit(1)
	}

	if *appName == "" {
		appName = extractNameFromGitRepo(playbookRepo)
	}

	lbos := strings.ToLower(*baseOS)

	if lbos == ubuntu {
		templateFile = ubuntuTemplate
	} else if lbos == centos {
		templateFile = centosTemplate
	}

	appVars := app{Name: *appName, AnsiblePlaybookRepo: *playbookRepo}
	fileVars := dockerfile{OSVersion: *osVersion, App: appVars}

	outputTemplate, err := template.ParseFiles(templateFile)
	fatalError(err, fmt.Sprintf("Unable to parse file: %v", templateFile))

	outputWriter, err := os.Create(outputFile)
	fatalError(err, fmt.Sprintf("Unable to create output file '%v'"))

	err = outputTemplate.Execute(outputWriter, fileVars)
	fatalError(err, fmt.Sprintf("Unable to write to '%v'", outputFile))
}
