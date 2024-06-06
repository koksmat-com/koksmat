package repos

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/koksmat-com/koksmat/kitchen"
)

type Release struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []any     `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}

/*
DownloadRepo downloads the latest release of a GitHub repository.

This is designed to be interactively with user output to the console.
*/
func DownloadRepo(repoOwner string, repoName string) (*string, error) {

	latestRelease, err := getLatestRelease(repoOwner, repoName)
	if err != nil {
		fmt.Println("Error getting latest release:", err)
		return nil, err
	}

	// if len(latestRelease.Assets) == 0 {
	// 	fmt.Println("No assets found for the latest release")
	// 	return
	// }
	fmt.Println("Latest release:", latestRelease.TagName)
	fmt.Println("Release notes:", latestRelease.Body)
	fmt.Println("Released by:", latestRelease.Author.Login)
	downloadURL := latestRelease.ZipballURL

	// Create a temporary directory
	tempDir, err := kitchen.SetupSessionPath("koksmat", kitchen.GenerateSessionId())
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return nil, err
	}

	extractDir := filepath.Join(tempDir, "extracted")

	if err := os.Mkdir(extractDir, 0755); err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return nil, err
	}

	// Download the zip file
	zipFilePath := filepath.Join(tempDir, "latest_release.zip")
	err = downloadFile(zipFilePath, downloadURL)
	if err != nil {
		fmt.Println("Error downloading zip file:", err)
		return nil, err
	}

	defer os.Remove(zipFilePath)

	// Extract the zip file
	err = Unzip(zipFilePath, extractDir, "latest")
	if err != nil {
		fmt.Println("Error extracting zip file:", err)
		return nil, err
	}

	fmt.Println("Zip file extracted to:", extractDir)
	fmt.Println("Zip file downloaded to:", zipFilePath)
	returnPath := filepath.Join(extractDir, "latest")
	return &returnPath, nil
	// You can now proceed to extract the zip file if needed
}

// getLatestRelease fetches the latest release information from the GitHub API.
func getLatestRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

// downloadFile downloads a file from the specified URL and saves it to the specified path.
func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Unzip(zipfile string, dst string, rootfoldername string) error {

	archive, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		//		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", filePath)

		}
		if f.FileInfo().IsDir() {
			//			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	destPath := path.Join(dst, rootfoldername)
	err = os.RemoveAll(destPath)
	if err != nil {
		return err
	}
	err = os.Rename(path.Join(dst, archive.File[0].Name), destPath)
	if err != nil {
		return err
	}

	return nil
}
