package kitchen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileDiff represents the differences between two folders

type FileDiffPair struct {
	Master  string
	Replica string
}

type FileInfo struct {
	FullPath     string
	RelativePath string
}
type FileDiff struct {
	Root               string
	FilesOnlyInMaster  []FileInfo
	FilesOnlyInReplica []FileInfo

	DifferentFiles []FileDiffPair
}

// Walks through the folder and populates the files map
func walkFolder(folderPath string, recurse bool) map[string]bool {
	files := make(map[string]bool)
	if !Exists(folderPath) {
		return files
	}
	fileWalker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", path, err)
			return nil
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(folderPath, path)
			files[relPath] = true
		}
		return nil
	}

	// Walk through the folder and populate files map
	if recurse {
		filepath.Walk(folderPath, fileWalker)
	} else {
		fileInfos, _ := os.ReadDir(folderPath)
		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				relPath := fileInfo.Name()
				files[relPath] = true
			}
		}
	}

	return files
}

// Reads .gitignore file and returns a list of patterns to ignore
func readGitIgnore(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	var patterns []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns, nil
}

// Checks if a file matches any of the patterns in the ignore list
func fileMatchesIgnoreList(filePath string, ignoreList []string) bool {
	for _, pattern := range ignoreList {
		match, _ := filepath.Match(pattern, filePath)
		if match {
			return true
		}
		// Check for patterns with directory wildcards
		if strings.HasSuffix(pattern, "/**/*") {
			// Trim  "/*" suffix
			prefix := strings.TrimSuffix(pattern, "/**/*")
			if strings.HasPrefix(filePath, prefix) {
				return true
			}
		}
	}
	return false
}

// Compares two folders and returns the differences
func compareFolders(masterPath, replicaPath string, recurse bool) FileDiff {
	masterFiles := walkFolder(masterPath, recurse)

	replicaFiles := walkFolder(replicaPath, recurse)

	diff := FileDiff{}
	// Read .gitignore file in master folder and create ignore list
	masterIgnoreList := []string{}
	masterGitIgnorePath := filepath.Join(masterPath, ".koksmatignore")
	if _, err := os.Stat(masterGitIgnorePath); err == nil {
		masterIgnoreList, err = readGitIgnore(masterGitIgnorePath)
		if err != nil {
			return FileDiff{}
		}
	}
	// Compare master and replica files
	for file := range masterFiles {
		if !replicaFiles[file] {
			if !fileMatchesIgnoreList(file, masterIgnoreList) {
				diff.FilesOnlyInMaster = append(diff.FilesOnlyInMaster, FileInfo{FullPath: filepath.Join(masterPath, file), RelativePath: file})
			}
		} else {
			masterContent, _ := os.ReadFile(filepath.Join(masterPath, file))
			replicaContent, _ := os.ReadFile(filepath.Join(replicaPath, file))
			if string(masterContent) != string(replicaContent) {
				pair := FileDiffPair{Master: filepath.Join(masterPath, file), Replica: filepath.Join(replicaPath, file)}
				diff.DifferentFiles = append(diff.DifferentFiles, pair)
			}
		}
	}

	// Files present only in replica folder
	for file := range replicaFiles {
		if !masterFiles[file] {
			diff.FilesOnlyInReplica = append(diff.FilesOnlyInReplica, FileInfo{FullPath: filepath.Join(replicaPath, file), RelativePath: file})
		}
	}

	return diff
}

type CompareOptions struct {
	PrintResults   bool
	PrintMergeLink bool
	CopyFunction   func(src string, dest string) error
	MergeFunction  func(src string, dest string) error
}

func Merge(srcFile, dstFile string) error {
	return nil
	// get parent directory from dstFile
	parentDir := filepath.Dir(dstFile)
	CreateIfNotExists(parentDir, 0755)
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Compare(masterRoot string, replicaRoot string, subfolders []string, recurse bool, action CompareOptions) ([]*FileDiff, error) {

	// Define subfolders to compare

	// Map to store comparison results
	comparisonResults := make(map[string]FileDiff)

	// setup an array to store the results
	result := []*FileDiff{}

	// Iterate over each subfolder
	for _, subfolder := range subfolders {
		// Construct full paths for master and replica subfolders
		masterSubfolder := filepath.Join(masterRoot, subfolder)
		replicaSubfolder := filepath.Join(replicaRoot, subfolder)

		// Perform folder comparison
		diff := compareFolders(masterSubfolder, replicaSubfolder, recurse)

		// Store the comparison result in the map
		fd := &FileDiff{Root: subfolder, FilesOnlyInMaster: diff.FilesOnlyInMaster, FilesOnlyInReplica: diff.FilesOnlyInReplica, DifferentFiles: diff.DifferentFiles}
		result = append(result, fd)

	}

	if action.PrintResults {
		// Print results
		for subfolder, diff := range comparisonResults {
			fmt.Printf("Comparison result for subfolder '%s':\n", subfolder)
			fmt.Println("Files only in master folder:")
			for _, file := range diff.FilesOnlyInMaster {
				fmt.Println(file)
			}
			fmt.Println("Files only in replica folder:")
			for _, file := range diff.FilesOnlyInReplica {
				fmt.Println(file)
			}
			fmt.Println("Different files:")
			for _, file := range diff.DifferentFiles {
				fmt.Println(file)
			}
			fmt.Println()
		}
	}

	for _, diff := range result {
		if action.CopyFunction != nil {
			for _, file := range diff.FilesOnlyInMaster {
				// Copy file from master to replica
				masterFile := file.FullPath
				replicaFile := filepath.Join(replicaRoot, file.RelativePath)
				fmt.Printf("Copying file '%s' to '%s'\n", masterFile, replicaFile)
				// Read master file content
				action.CopyFunction(masterFile, replicaFile)

			}
		}
		if action.PrintMergeLink {
			for _, file := range diff.DifferentFiles {
				// Copy file from master to replica

				fmt.Printf("code --diff '%s' '%s'\n", file.Master, file.Replica)

			}
		}
		if action.MergeFunction != nil {
			for _, file := range diff.DifferentFiles {
				// Copy file from master to replica

				// Read master file content
				action.MergeFunction(file.Master, file.Replica)

			}
		}
		fmt.Println()
	}

	return result, nil
}
