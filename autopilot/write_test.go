package autopilot

import (
	"fmt"
	"testing"
)

func TestWrite(t *testing.T) {
	req := WriteRequest{
		Files: []FileData{
			{Directory: "dir1", FileName: "file1.txt", Data: "Hello, World!"},
			{Directory: "dir2", FileName: "file2.txt", Data: "Go is awesome!"},
		},
		Category:   "example",
		Force:      false,
		RootFolder: "",
	}

	statuses, _ := handleFileWriteRequest(req)

	if statuses != nil {
		for _, status := range statuses {
			fmt.Printf("File: %s, Status: %s, Error: %v\n", status.FileName, status.Status, status.Error)
		}

	} else {
		fmt.Println("All files written successfully")
	}

}
