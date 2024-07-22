package autopilot

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileData represents the structure of each file
type FileData struct {
	Directory string `json:"directory"`
	FileName  string `json:"file_name"`
	Data      string `json:"data"`
}

// Request represents the request object
type WriteRequest struct {
	Files      []FileData `json:"files"`
	Category   string     `json:"category"`
	Force      bool       `json:"force"`
	RootFolder string     `json:"root_folder"`
}

// FileStatus represents the status of a file operation
type FileStatus struct {
	FileName string `json:"file_name"`
	Status   string `json:"status"`
	Error    error  `json:"error"`
}

// ensureDirectories creates all necessary directories relative to rootFolder
func ensureDirectories(rootFolder string, files []FileData) error {
	for _, file := range files {
		dirPath := filepath.Join(rootFolder, file.Directory)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// writeFile writes data to a file
func writeFile(path string, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}

// handleFileWriteRequest handles the file writing request
func handleFileWriteRequest(req WriteRequest) ([]FileStatus, error) {
	// Ensure all directories are created relative to rootFolder
	if err := ensureDirectories(req.RootFolder, req.Files); err != nil {
		return nil, err
	}

	var statuses []FileStatus
	var existingFiles []FileStatus

	for _, file := range req.Files {
		filePath := filepath.Join(req.RootFolder, file.Directory, file.FileName)
		if fileExists(filePath) {
			if !req.Force {
				existingFiles = append(existingFiles, FileStatus{FileName: filePath, Status: "existing", Error: nil})
			}
		}
	}

	if len(existingFiles) > 0 && !req.Force {
		return existingFiles, nil
	}

	for _, file := range req.Files {
		filePath := filepath.Join(req.RootFolder, file.Directory, file.FileName)
		if err := writeFile(filePath, file.Data); err != nil {
			statuses = append(statuses, FileStatus{FileName: filePath, Status: "failed", Error: err})
		} else {
			statuses = append(statuses, FileStatus{FileName: filePath, Status: "written", Error: nil})
		}
	}

	if len(statuses) > 0 {
		for _, status := range statuses {
			if status.Status == "failed" {
				return statuses, errors.New("one or more files could not be written")
			}
		}
	}

	return nil, nil
}
