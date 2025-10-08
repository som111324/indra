package services

import (
	"cloud/internal/detector"
	"fmt"
	"os"
	"strings"
)

type ScriptService struct{}

func NewScriptService() *ScriptService {
	return &ScriptService{}
}

func (s *ScriptService) generatestarupscript(repoURL, projectType, startcommand string) (string, error) {
	var templatepath string
	var installcommand string

	switch projectType {
	case "nodejs":
		templatepath = "internal/templates/nodejs.sh"
		installcommand = "npm install"
	case "python":
		templatepath = "internal/templates/python.sh"
		installcommand = "pip install -r requirements.txt"
	default:
		return "", fmt.Errorf("unsupported project type: %s", projectType)
	}

	templateBytes, err := os.ReadFile(templatepath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	script := string(templateBytes)

	script = strings.ReplaceAll(script, "{{REPO_URL}}", repoURL)
	script = strings.ReplaceAll(script, "{{INSTALL_COMMAND}}", installcommand)
	script = strings.ReplaceAll(script, "{{START_COMMAND}}", startcommand) //it was startcommand

	return script, nil

}

func (s *ScriptService) DetectandGenerateScript(repoUrl, startCommand string) (string, error) {

	projectType, err := detector.DetectProjectType(repoUrl)
	if err != nil {
		return "", fmt.Errorf("failed to detect project type: %v", err)
	}

	if startCommand == "" {
		startCommand = projectType.StartCommand //staypu command
	}

	return s.generatestarupscript(repoUrl, projectType.Type, startCommand)
}
