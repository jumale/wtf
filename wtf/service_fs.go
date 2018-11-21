package wtf

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type FileSystem struct {
	OpenFileUtil string
}

// OpenFile opens the file defined in `path` via the operating system
func (fs FileSystem) OpenFile(path string) {
	openFile := fs.OpenFileUtil
	if openFile == "" {
		openFile = "open"
	}

	filePath, _ := fs.ExpandHomeDir(path)
	cmd := exec.Command(openFile, filePath)

	ExecuteCommand(cmd)
}

// Dir returns the home directory for the executing user.
// An error is returned if a home directory cannot be detected.
func (fs FileSystem) Home() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func (fs FileSystem) ExpandHomeDir(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := fs.Home()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

func (fs FileSystem) ReadFileBytes(filePath string) ([]byte, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	return fileData, nil
}

// CreateFile creates the file, if it does not already exist.
// If the file exists it does not recreate it.
// If successful, returns the absolute path to the file
// If unsuccessful, returns an error
func (fs FileSystem) CreateFile(path string) (string, error) {
	// Check if the file already exists; if it does not, create it
	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}

	// if this is another error, but not "not-exist"
	if os.IsExist(err) {
		return "", nil
	}

	// create dir recursively if not exists
	dirPath := filepath.Dir(path)
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModePerm)
	}
	if err != nil {
		return "", errors.WithStack(err)
	}

	// create file
	_, err = os.Create(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return path, nil
}
