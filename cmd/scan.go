package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/methlock/code-stats/models"
	"github.com/methlock/code-stats/utils"
	"github.com/peteabre/colprint"
	"github.com/spf13/cobra"
)

var (
	allFiles       = []string{}
	filesStats     = []models.FileStats{}
	extensionStats = map[string]*models.ExtensionStats{}

	// string flags
	excludedPathsInput     string
	excludedPaths          []string
	allowedExtensionsInput string
	allowedExtensions      []string

	// bool flags
	searchHidden    = false
	processBinaries = false
	printExts       = false
	printFiles      = false
)

// checkPaths checks if given paths exists
func checkPaths(paths []string) {

	for _, path := range paths {

		if path == "." {
			fmt.Println("Relative paths, like '.', are not supported!")
			os.Exit(1)
		}

		_, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Path '%v' does not exists!\n", path)
			os.Exit(1)
		}
	}
}

// processFile for some statistics.
func processFile(filePath string) {

	splittedPath := strings.Split(filePath, "/")
	fileName := splittedPath[len(splittedPath)-1] // last item
	fs := models.FileStats{Name: fileName, Path: filePath}
	fs.GetExtension()

	// check if extensions filter enabled
	if len(allowedExtensions) != 0 {
		if !utils.IsThisIn(fs.Extension, allowedExtensions) {
			return
		}
	}

	// skip binaries - files with no extension
	if fs.Extension == "none" && !processBinaries == true {
		return
	}

	fs.CountLines()

	filesStats = append(filesStats, fs)

	// stats for extensions
	ext, ok := extensionStats[fs.Extension]
	if !ok {
		extStats := models.ExtensionStats{
			Extension:    fs.Extension,
			Files:        1,
			Lines:        fs.TotalLines,
			CodeLines:    fs.CodeLines,
			CommentLines: fs.CommentLines,
		}
		extensionStats[fs.Extension] = &extStats
	} else {
		ext.Files++
		ext.Lines += fs.TotalLines
		ext.CodeLines += fs.CodeLines
		ext.CommentLines += fs.CommentLines
	}
}

// collectFiles walks given path recursively - calls self in case of dir.
func collectFiles(root string) {

	objectsInPath, err := ioutil.ReadDir(root)
	if err != nil {
		// handling file in root path
		if err.Error() == "readdirent: not a directory" {
			processFile(root)
			return
		}
		log.Fatal(err)
	}

	// walk throug all objects
	for _, object := range objectsInPath {
		fullPath := utils.JoinPaths(root, object.Name())

		// skip hidden
		if !searchHidden == true && strings.HasPrefix(object.Name(), ".") {
			continue
		}

		// excluded paths
		if len(excludedPaths) > 0 && utils.IsThisIn(fullPath, excludedPaths) {
			continue
		}

		// enter directories or append file
		switch object.IsDir() {
		case true:
			collectFiles(fullPath)
		case false:
			allFiles = append(allFiles, fullPath)
		}
	}
}

// printFileStats prints collected stats
func printResults() {

	// collect stats
	totalLines := 0
	totalCode := 0
	totalComments := 0
	for _, file := range filesStats {
		totalLines = totalLines + file.TotalLines
		totalCode = totalCode + file.CodeLines
		totalComments = totalComments + file.CommentLines
	}
	fmt.Println("Overall stats:")
	fmt.Printf(" Number of files: %v\n", len(filesStats))
	fmt.Printf(" Total lines: %v\n", totalLines)
	fmt.Printf(" Code lines: %v\n", totalCode)
	fmt.Printf(" Comment lines: %v\n", totalComments)

	// statistics for extensions
	if printExts == true {
		// NOTE: Sadly colprint accepts only slice, not map
		//  thus conversion is needed.
		s := []*models.ExtensionStats{}
		for _, e := range extensionStats {
			s = append(s, e)
		}
		fmt.Println("\nStats by extensions")
		colprint.Print(s)
		fmt.Println()
	}

	// statistics for individual files
	if printFiles == true {
		fmt.Println("\nStats by files:")
		colprint.Print(filesStats)
	}
	fmt.Println()
}

func parseFlags() {
	if len(allowedExtensionsInput) > 0 {
		allowedExtensions = strings.Split(allowedExtensionsInput, " ")
	}
	if len(excludedPathsInput) > 0 {
		excludedPaths = strings.Split(excludedPathsInput, " ")
	}
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [PATH or PATHS]",
	Short: "Scan path or paths for files and process them",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, paths []string) {

		// flags setting
		searchHidden, _ = cmd.Flags().GetBool("hidden")
		processBinaries, _ = cmd.Flags().GetBool("process-no-ext")
		printExts, _ = cmd.Flags().GetBool("print-exts")
		printFiles, _ = cmd.Flags().GetBool("print-files")
		parseFlags()

		checkPaths(paths)
		for _, path := range paths {
			collectFiles(path)
		}
		for _, file := range allFiles {
			processFile(file)
		}

		printResults()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&excludedPathsInput, "exclude-path", "e", "", "Paths to exclude. Usage: '-e \"path1 path2\"'")
	scanCmd.Flags().StringVarP(&allowedExtensionsInput, "allowed-exts", "x", "", "Specify the extensions to process. Pass them without dot, like '-x \"go py\"'")
	scanCmd.Flags().BoolP("process-no-ext", "n", false, "Process also files with no extensions, like binaries")
	scanCmd.Flags().BoolP("hidden", "s", false, "Search and process also hidden folders and files")
	scanCmd.Flags().BoolP("print-exts", "p", false, "Prints also information by extensions")
	scanCmd.Flags().BoolP("print-files", "f", false, "Prints stats for each file")
}
