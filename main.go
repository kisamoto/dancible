package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"text/template"
)

const (
	ubuntu = "ubuntu"
	centos = "centos"
	debian = "debian"

	outputFile = "Dockerfile"

	ubuntuUpdate          = "apt-get -qq update"
	ubuntuUpgrade         = "apt-get upgrade -y"
	ubuntuInstallSoftware = "apt-get install -y python-pip git"

	centosUpdate          = "yum -y update"
	centosUpgrade         = ""
	centosInstallSoftware = "yum -y install python-pip git"

	debianUpdate          = "aptitude -qq update"
	debianUpgrade         = "aptitude full-upgrade -y"
	debianInstallSoftware = "aptitude install -y python-pip git"

	flagNameBaseOS    = "os"
	flagDefaultBaseOS = "ubuntu"
	flagUsageBaseOS   = "the base operating system to use in the container (ubuntu, centos, debian)"

	flagNameOSVersion    = "version"
	flagDefaultOSVersion = "latest"
	flagUsageOSVersion   = "the version of the docker image to use"

	flagNamePlaybookRepo    = "repo"
	flagDefaultPlaybookRepo = ""
	flagUsagePlaybookRepo   = "git URL to pull an ansible playbook and configure this container"

	flagNameAppName    = "name"
	flagDefaultAppName = ""
	flagUsageAppName   = "name of docker container to be produced"

	flagNameBranch    = "branch"
	flagDefaultBranch = ""
	flagUsageBranch   = "a branch to be checked out from playbook repo containing 'site.yml'"

	templateFile = "Dockerfile.template"
)

var (
	osVars                                               map[string]*osvars
	baseOS, osVersion, playbookRepo, branchName, appName *string
)

type osvars struct {
	Name    string
	Version string
	Update  string
	Upgrade string
	Install string
}

func (self *osvars) setVersion(version string) {
	self.Version = version
}

type app struct {
	Name                string
	AnsiblePlaybookRepo string
	Branch              string
}

type dockerfile struct {
	App app
	OS  osvars
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
        log.Println(err)
		log.Fatalf(errorMessage)
	}
}

func verify() {
	if *playbookRepo == "" {
		fmt.Println("A repository is required to pull down a playbook and configure this container.")
		os.Exit(1)
	}

	if *appName == "" {
		appName = extractNameFromGitRepo(playbookRepo)
	}
}

func init() {
	osVars = make(map[string]*osvars)
	osVars[ubuntu] = &osvars{Name: ubuntu, Update: ubuntuUpdate, Upgrade: ubuntuUpgrade, Install: ubuntuInstallSoftware}
	osVars[centos] = &osvars{Name: centos, Update: centosUpdate, Upgrade: centosUpgrade, Install: centosInstallSoftware}
	osVars[debian] = &osvars{Name: debian, Update: debianUpdate, Upgrade: debianUpgrade, Install: debianInstallSoftware}

	baseOS = flag.String(flagNameBaseOS, flagDefaultBaseOS, flagUsageBaseOS)
	osVersion = flag.String(flagNameOSVersion, flagDefaultOSVersion, flagUsageOSVersion)
	playbookRepo = flag.String(flagNamePlaybookRepo, flagDefaultPlaybookRepo, flagUsagePlaybookRepo)
	appName = flag.String(flagNameAppName, flagDefaultAppName, flagUsageAppName)
	branchName = flag.String(flagNameBranch, flagDefaultBranch, flagUsageBranch)
}

func main() {

	var err error

	flag.Parse()
	verify()

	lbos := strings.ToLower(*baseOS)
	osVal, ok := osVars[lbos]
	if !ok {
		fmt.Printf("No configuration for os '%v' available.\n", lbos)
		os.Exit(1)
	}

	osVal.setVersion(*osVersion)
	appVars := app{Name: *appName, AnsiblePlaybookRepo: *playbookRepo, Branch: *branchName}
	fileVars := dockerfile{App: appVars, OS: *osVal}

	outputTemplate, err := template.ParseFiles(templateFile)
	fatalError(err, fmt.Sprintf("Unable to parse file: %v", templateFile))

	outputWriter, err := os.Create(outputFile)
	fatalError(err, fmt.Sprintf("Unable to create output file '%v'"))

	err = outputTemplate.Execute(outputWriter, fileVars)
	fatalError(err, fmt.Sprintf("Unable to write to '%v'", outputFile))
}
