package kitchen

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/koksmat-com/koksmat/connectors"
)

// /Users/nielsgregersjohansen/kitchens/sharepoint-branding/install/20 apply-sitetemplate.ps1
func extractEnvVarsFromPowerShellFile(fileContent string) map[string]string {
	envVars := make(map[string]string)

	// find all matches for  $env: in the fileContent string and return the value after the :
	// the $env should be made case insensitive

	re := regexp.MustCompile(`(?i)\$env:([a-zA-Z0-9_]+)`)
	matches := re.FindAllStringSubmatch(fileContent, -1)
	for _, match := range matches {
		envVar := match[1]
		envVars[envVar] = os.Getenv(envVar)
	}

	return envVars
}

func ParsePowerShellFile(fileContent string) (string, []string, error) {
	markdown := ""
	powershell := ""
	commentBlock := strings.Split(string(fileContent), "<#")

	for _, block := range commentBlock {
		if strings.Index(block, "#>") == -1 {
			if block != "" {
				powershell += block
				markdown += "\n```powershell\n"
				markdown += block
				markdown += "\n```\n"
			}
		} else {
			comments := strings.Split(block, "#>")
			markdown += comments[0]
			if len(comments) > 1 {
				powershell += comments[1]
				markdown += "\n```powershell\n"
				markdown += comments[1]
				markdown += "\n```\n"

			}
		}

	}

	envs := extractEnvVarsFromPowerShellFile(powershell)

	keys := make([]string, 0, len(envs))
	for k := range envs {
		keys = append(keys, k)
	}

	return markdown, keys, nil
}

func ParseGoFile(fileContent string) (string, error) {
	markdown := ""
	powershell := ""
	commentBlock := strings.Split(string(fileContent), "/*")

	for _, block := range commentBlock {

		block = strings.TrimPrefix(block, "\n")

		if strings.Index(block, "*/") == -1 {
			if block != "" {
				powershell += block
				markdown += "\n```go\n"
				markdown += block
				markdown += "\n```\n"
			}
		} else {
			comments := strings.Split(block, "*/")
			markdown += comments[0]
			if len(comments) > 1 {
				powershell += comments[1]
				markdown += "\n```go\n"
				markdown += comments[1]
				markdown += "\n```\n"

			}
		}

	}

	return markdown, nil
}
func ReadMarkdownFromPowerShell(filepath string) (string, []string, error) {

	if !FileExists(filepath) {
		return "", nil, nil
	}
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return "", nil, err
	}
	return ParsePowerShellFile(string(fileContent))
}

func ReadMarkdownFromGo(filepath string) (string, error) {

	if !FileExists(filepath) {
		return "", fmt.Errorf("File %s not found", filepath)
	}
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return ParseGoFile(string(fileContent))
}

type PowershellParameter struct {
	IsMndatory bool   `json:"mandatory"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Help       string `json:"help"`
}

func GetPowerShellFileParameters(filepath string) ([]PowershellParameter, error) {
	cmd := fmt.Sprintf(`
	
		$x = (Get-Command "%s").Parameters
		$parameters = @()
		foreach ($key in $x.Keys) {
			if ($x[$key].Attributes[0].Position -gt -1){
			  
				$p = @{
					Name = $key
					Position = $x[$key].Attributes[0].Position
					HelpMessage = $x[$key].Attributes[0].HelpMessage
					ParameterType = $x[$key].ParameterType.FullName
					Mandatory = $x[$key].Attributes[0].Mandatory
				}
				$parameters += $p
			  
			}
		   
		}
		$parameters | ConvertTo-Json
		
	`, filepath)

	bytes, err := connectors.Execute("pwsh", *&connectors.Options{}, "-nologo", "-noprofile", "-Command", cmd)
	if err != nil {
		return nil, err
	}

	parms := make([]PowershellParameter, 0)

	s := string(bytes)
	if strings.Index(s, "[") == -1 {
		s = fmt.Sprintf("[%s]", s)
	}

	err = json.Unmarshal([]byte(s), &parms)
	if err != nil {
		return nil, err
	}
	return parms, nil

}
