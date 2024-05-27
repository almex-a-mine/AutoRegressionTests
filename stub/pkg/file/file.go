package file

import (
	"os"
	"path/filepath"
)

func GetEnv(envname string) string {
	result := os.Getenv(envname)
	return result
}
func GetCurrentDir() string {
	result := ""
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	result = path
	return result
}

func FileExists(filename string) bool {
	if f, err := os.Stat(filename); os.IsNotExist(err) || f.IsDir() {
		return false
	} else {
		return true
	}
}
func DirExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	} else {
		return true
	}
}
func AdjustFileName(filename string) string {
	result := filename
	filelen := len(filename)
	if filelen != 0 {
		c := filename[filelen-1:]
		if c != "\\" {
			result = result + "\\"
		}
	}
	return result
}
