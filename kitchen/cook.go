package kitchen

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/koksmat-com/koksmat/connectors"

	"github.com/spf13/viper"
)

func clearString(str string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func GenerateSessionId() string {
	return clearString(time.Now().Format(time.RFC3339Nano))
}
func SetupSessionPath(kitchenPath string, sessionId string) (string, error) {
	sessionPath := path.Join(kitchenPath, ".koksmat", "sessions", sessionId)
	err := os.MkdirAll(sessionPath, 0755)
	if err != nil {
		return "", err
	}

	return sessionPath, nil
}
func getEnvironmentFilePath(startPath string, tenantName string) (string, error) {

	if startPath == "" {
		return "", fmt.Errorf("Environment file %s not found", tenantName)
	}

	tenantenvPath := path.Join(startPath, fmt.Sprintf(".env-%s", tenantName))
	if FileExists(tenantenvPath) {
		return tenantenvPath, nil
	}

	envPath := path.Join(startPath, ".env")
	if FileExists(envPath) {
		return envPath, nil
	}
	startPath = path.Dir(startPath)
	return getEnvironmentFilePath(startPath, tenantName)

}

type environmentStack struct {
	environmentPath string
	tenantname      string
	environment     []string
}

func getEnvironmentStack(endPath string, index int, stack []environmentStack) ([]environmentStack, error) {

	if endPath == "" {
		return nil, fmt.Errorf("You cannot start from the root")
	}

	elems := strings.Split(endPath, "/")
	startPath := "/"
	if index == len(elems) {
		startPath = endPath
	} else {
		if index > 0 {
			startPath = fmt.Sprintf("%s", strings.Join(elems[:index+1], "/"))
		}
	}
	envPath := path.Join(startPath, ".env")
	if FileExists(envPath) {
		env, err := ReadEnvironmentVariables(envPath)
		if err != nil {
			return nil, err
		}
		stack = append(stack, environmentStack{environmentPath: envPath, environment: env})

	}

	if startPath == endPath {
		return stack, nil
	}
	return getEnvironmentStack(endPath, index+1, stack)

}

func ReadEnvironmentVariables(filepath string) ([]string, error) {

	if !FileExists(filepath) {
		return nil, nil
	}
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(fileContent), "\n"), nil
}

func PowerShellEnvironmentVariables(filepath string) string {

	// environment, err := ReadEnvironmentVariables(filepath)
	powershellEnv := fmt.Sprintf(`
# variables read from %s
`, filepath)
	stack, err := getEnvironmentStack(filepath, 0, []environmentStack{})
	if err != nil {
		return fmt.Sprintf("# Error reading environment %s", err)
	}
	for _, env := range stack {
		environment := env.environment
		powershellEnv += "#--------------------------------------\n"
		powershellEnv += fmt.Sprintf("# %s\n", env.environmentPath)
		powershellEnv += "#--------------------------------------\n"
		for _, env := range environment {
			if (env == "") || (strings.HasPrefix(env, "#")) {
				continue
			}
			pairs := strings.Split(env, "=")
			name, values := pairs[0], pairs[1:]

			val := strings.Join(values, "=")
			powershellEnv += fmt.Sprintf(`$env:%s="%s"
`, name, val)
		}
	}
	return powershellEnv
}

func getConnectionsScript(tenant string, connectionString, sessionPath string) (string, error) {
	connectScript := ""
	kitchenRoot := viper.GetString("KITCHENROOT")

	if connectionString == "" {
		connectionString = "sharepoint"
	}
	if connectionString != "" {
		connections := strings.Split(connectionString, ",")
		for _, connection := range connections {
			connection := strings.TrimSpace(connection)
			connectionPath := path.Join(kitchenRoot, ".koksmat", "tenants", tenant, connection)
			if !FileExists(path.Join(connectionPath, "connect.ps1")) {
				return "", fmt.Errorf("Connection %s not found for tenant %s", connection, tenant)
			}
			CreateIfNotExists(path.Join(sessionPath, connection), 0755)
			err := CopyDirectory(connectionPath, path.Join(sessionPath, connection))
			if err != nil {
				return "", err
			}
			connectScript += fmt.Sprintf(`. $PSScriptRoot/%s/connect.ps1 
			`, connection)
		}
	}

	return connectScript, nil
}

func Cook(isDebug bool, tenantName string, kitchenName string, stationName string, journeyId string, filename string, environment []string, args ...string) (string, error) {
	// 	Run: run,
	// }
	// return cmd

	//context := connectors.Context{}

	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, kitchenName)
	scriptPath := path.Join(root, kitchenName, stationName, filename)
	//if (journeyId != "") && (journeyId != "null") {
	journeyId = "latest"
	//}
	sessionId := GenerateSessionId()

	sessionPath, err := SetupSessionPath(kitchenPath, sessionId)
	//hostConnections := []HostConnection{}

	if err != nil {
		return "", err

	}

	//envPath, err := getEnvironmentFilePath(kitchenPath, tenantName)
	psEnv := ""
	if err == nil {
		psEnv = PowerShellEnvironmentVariables(kitchenPath)

	}
	markdown, _, err := ReadMarkdownFromPowerShell(scriptPath)
	if err != nil {
		return "", err
	}
	_, metadata, err := ParseMarkdown(false, "", markdown)
	if err != nil {
		return "", err
	}
	connectScript, err := getConnectionsScript(tenantName, GetMetadataProperty(metadata, "connection", ""), sessionPath)
	if err != nil {
		return "", err
	}

	err = Copy(scriptPath, path.Join(sessionPath, "script.ps1"))
	if err != nil {
		return "", err
	}
	workDir := path.Join(kitchenPath, ".koksmat", "workdir")
	CreateIfNotExists(workDir, 0755)
	fullScript := fmt.Sprintf(`
# --------------------------------------
#  Generated by koksmat
# --------------------------------------
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"
$VerbosePreference = "SilentlyContinue"
$DebugPreference = "SilentlyContinue"

%s
$ENV:WORKDIR="%s"
$ENV:MILLERSESSIONID="%s"
$ENV:MILLERJOURNEY="%s"
$ENV:KITCHEN="%s"
Start-Transcript -Path "$PSScriptRoot/transcript.txt" -Append
$result=""
write-host "Running script"
%s
. $PSScriptRoot/script.ps1 ##ARGS##  
Out-File -InputObject $result  -FilePath $PSScriptRoot/output.txt -Encoding:utf8NoBOM
Stop-Transcript
return # the remaining code is not executed in the current session
$libraryName = "Miller Sessions"
Connect-PnPOnline -Url $ENV:SITEURL  -ClientId $PNPAPPID -Tenant $PNPTENANTID -CertificatePath "$PNPCERTIFICATEPATH"

New-PnPList -Title $libraryName -Template DocumentLibrary -ErrorAction SilentlyContinue
Add-PnPFolder -Name $env:KITCHEN -Folder $libraryName -ErrorAction SilentlyContinue
Add-PnPFolder -Name $env:MILLERJOURNEY -Folder "$($libraryName)/$($env:KITCHEN)" -ErrorAction SilentlyContinue

Add-PnPFolder -Name $env:MILLERSESSIONID -Folder "$($libraryName)/$($env:KITCHEN)/$($env:MILLERJOURNEY)" -ErrorAction SilentlyContinue

$files = Get-ChildItem -Path "$env:WORKDIR" -Filter *.json 
foreach ($file in $files) {
    Add-PnPFile -Path $file.FullName -Folder "$($libraryName)/$($env:KITCHEN)/$($env:MILLERJOURNEY)" 
    
}
Add-PnPFile -Path "$PSScriptRoot/output.txt" -Folder "$($libraryName)/$($env:KITCHEN)/$($env:MILLERJOURNEY)/$($env:MILLERSESSIONID)" 
Add-PnPFile -Path "$PSScriptRoot/script.ps1" -Folder "$($libraryName)/$($env:KITCHEN)/$($env:MILLERJOURNEY)/$($env:MILLERSESSIONID)" 
Add-PnPFile -Path "$PSScriptRoot/transcript.txt" -Folder "$($libraryName)/$($env:KITCHEN)/$($env:MILLERJOURNEY)/$($env:MILLERSESSIONID)" 
    


`, psEnv, workDir, sessionId, journeyId, kitchenName, connectScript)

	os.WriteFile(path.Join(sessionPath, "run.ps1"), []byte(fullScript), 0755)
	if isDebug {
		_, err = connectors.Execute("code", *&connectors.Options{Dir: sessionPath}, "run.ps1")

		if err != nil {
			return "", err
		}

		return fmt.Sprintf(`
A new sessions has been created for you and the file run.ps1 has been opened in Visual Studio Code.
	
	Session path: %s`, sessionPath), nil
	} else {
		return sessionPath, nil
	}
}

func KitchenCommand(kitchenName string, cmd string, args ...string) (string, error) {
	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, kitchenName)

	output, err := connectors.Execute(cmd, *&connectors.Options{Dir: kitchenPath}, args...)

	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "Done", nil
	}
	return string(output), nil
}

func Build(kitchenName string, args ...string) (string, error) {
	return KitchenCommand(kitchenName, "go", "install")
}

func Open(kitchenName string, args ...string) (string, error) {
	return KitchenCommand(kitchenName, "code", ".")
}

func Launch(kitchenName string, args ...string) (string, error) {
	return KitchenCommand(kitchenName, kitchenName, "-h")
}

func Download(url string, destFilePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	fmt.Println("status", resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to download file %s", url)
	}

	// Create the file
	out, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
