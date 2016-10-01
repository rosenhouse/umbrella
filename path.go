package umbrella

import (
	"flag"
	"path/filepath"
)

func getInProcessFilePath() string {
	filename := flag.Lookup("test.coverprofile").Value.String()
	if filename == "" {
		return ""
	}
	dir := flag.Lookup("test.outputdir").Value.String()
	if dir == "" {
		return filename
	}
	return filepath.Join(dir, filename)
}

func getProfilePath() (string, error) {
	orig := getInProcessFilePath()
	if orig == "" {
		return "", nil
	}

	return filepath.Abs(orig + ".external.coverprofile")
}
