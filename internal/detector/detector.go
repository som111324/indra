package detector

import (
	"fmt"
	"net/http"
	"strings"
)

type ProjectType struct {
	Type           string
	InstallCommand string
	StartCommand   string
}

var nodejs = ProjectType{
	Type:           "nodejs",
	InstallCommand: "npm install",
	StartCommand:   "npm start",
}

var python = ProjectType{
	Type:           "python",
	InstallCommand: "pip install -r requirements.txt",
	StartCommand:   "python main.py",
}

func DetectProjectType(repoUrl string) (ProjectType, error) {
	rawurl := converTorawURL(repoUrl)

	if fileExists(rawurl + "package.json") {
		return nodejs, nil
	} else if fileExists(rawurl + "requirements.txt") {
		return python, nil
	} else if fileExists(rawurl + "app.py") {
		return python, nil
	} else if fileExists(rawurl + "main.py") {
		return python, nil
	}

	return ProjectType{}, fmt.Errorf("could not detect project type from repository")

}

func converTorawURL(repoUrl string) string {
	repoUrl = strings.TrimSuffix(repoUrl, ".git")
	rawurl := strings.Replace(repoUrl, "github.com", "raw.githubusercontent.com", 1)
	return rawurl + "/main/"
}

func fileExists(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
