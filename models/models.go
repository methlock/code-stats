package models

import (
	"bufio"
	"os"
	"strings"

	"github.com/methlock/code-stats/utils"
)

// FileStats holds some statistics for file.
type FileStats struct {
	Name         string `colprint:"Name,1"`
	Extension    string
	TotalLines   int    `colprint:"Total lines,2"`
	CodeLines    int    `colprint:"Code lines,3"`
	CommentLines int    `colprint:"Comment lines,4"`
	Path         string `colprint:"Path,5"`
}

// GetExtension resolves file extension.
func (fs *FileStats) GetExtension() {
	if strings.Contains(fs.Name, ".") {
		fs.Extension = strings.Split(fs.Name, ".")[1]
	} else {
		fs.Extension = "none"
	}
}

// CountLines counts lines (code and comments) in file.
func (fs *FileStats) CountLines() {
	file, _ := os.Open(fs.Path)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		row := fileScanner.Text()
		row = strings.TrimSpace(row) // remove whitespaces
		if utils.IsCommentLine(row) {
			fs.CommentLines++
		} else {
			fs.CodeLines++
		}
		fs.TotalLines++
	}

}

// ExtensionStats holds statistics by extension.
type ExtensionStats struct {
	Extension    string `colprint:"Extension,1"`
	Files        int    `colprint:"Files,2"`
	Lines        int    `colprint:"Total lines,3"`
	CodeLines    int    `colprint:"Code lines,4"`
	CommentLines int    `colprint:"Comment lines,5"`
}
