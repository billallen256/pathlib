package pathlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Path type alias
type Path string

// Exists returns true if the Path exists.
func (p Path) Exists() bool {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return false
	}

	_, err = os.Stat(absPath)

	if err != nil {
		return false
	}

	return true
}

// ReadBytes reads all the bytes from a file Path.
func (p Path) ReadBytes() ([]byte, error) {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(absPath)

	if err != nil {
		return nil, err
	}

	return contents, nil
}

// IsDir returns true if the Path is a directory. Note that false is returned if the Path does not exist.
func (p Path) IsDir() bool {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return false
	}

	stat, err := os.Stat(absPath)

	if err != nil {
		return false
	}

	mode := stat.Mode()

	if mode.IsDir() {
		return true
	}

	return false
}

// IsFile returns true if the Path is a file. Note that false is returned if the Path does not exist.
func (p Path) IsFile() bool {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return false
	}

	stat, err := os.Stat(absPath)

	if err != nil {
		return false
	}

	mode := stat.Mode()

	if mode.IsRegular() {
		return true
	}

	return false
}

// Permissions returns the Path's permissions as from os.Stat().
func (p Path) Permissions() (os.FileMode, error) {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return 0, err
	}

	stat, err := os.Stat(absPath)

	if err != nil {
		return 0, err
	}

	return stat.Mode().Perm(), nil
}

// Glob returns a list of Paths that match the pattern within the directory.
func (p Path) Glob(pattern string) ([]Path, error) {
	if !p.IsDir() {
		return nil, fmt.Errorf("Glob only works on directories: %s", p)
	}

	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return nil, err
	}

	absPattern := filepath.Join(absPath, pattern)
	matches, err := filepath.Glob(absPattern)

	if err != nil {
		return nil, err
	}

	matchPaths := make([]Path, 0, len(matches))

	for _, match := range matches {
		matchPaths = append(matchPaths, Path(match))
	}

	return matchPaths, nil
}

// Resolve returns the absolute form of the Path, if it exists.
func (p Path) Resolve() (Path, error) {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return p, err
	}

	resolvedPath := Path(absPath)

	if !resolvedPath.Exists() {
		return p, fmt.Errorf("Cannot resolve path that does not exist: %s", resolvedPath)
	}

	return resolvedPath, nil
}

// WithSuffix returns a new Path with the specified suffix (file extension). If the Path has no existing extension, the new extension will be added. If the Path has an extension, it will be replaced.
func (p Path) WithSuffix(suffix string) Path {
	pStr := string(p)
	suffix = strings.TrimSpace(suffix)
	oldSuffix := filepath.Ext(pStr)

	if len(oldSuffix) == 0 && len(suffix) == 0 {
		return p
	}

	if len(oldSuffix) == 0 && len(suffix) > 0 {
		return Path(pStr + "." + suffix)
	}

	oldSuffixLen := len(oldSuffix)

	if len(suffix) > 0 {
		oldSuffixLen -= 1 // don't include the dot for removal
	}

	withoutOldSuffix := pStr[0 : len(pStr)-oldSuffixLen]
	return Path(withoutOldSuffix + suffix)
}

// Touch creates a file at the Path if it does not already exist.
func (p Path) Touch() error {
	if p.Exists() {
		return nil
	}

	f, err := os.Create(string(p))

	if err != nil {
		return err
	}

	f.Close()
	return nil
}

// Age returns the last modification time of the Path, if it exists.
func (p Path) Age(now time.Time) (time.Duration, error) {
	if !p.Exists() {
		return time.Duration(0), fmt.Errorf("%s does not exist", p)
	}

	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return time.Duration(0), err
	}

	stat, err := os.Stat(absPath)

	if err != nil {
		return time.Duration(0), err
	}

	return now.Sub(stat.ModTime()), nil
}

// JoinPath returns any number of Paths joined by the OS specific path separator (eg. / or \).
func (p Path) JoinPath(paths ...Path) Path {
	ret := string(p)

	for _, path := range paths {
		ret = filepath.Join(ret, string(path))
	}

	return Path(ret)
}

// Name returns only the last portion of the Path as a string.
func (p Path) Name() string {
	return filepath.Base(string(p))
}

// Parent returns the last directory in the Path. For a file, it returns the directory that the file is in.  For a directory, it just returns the directory, not the directory above it.
func (p Path) Parent() Path {
	return Path(filepath.Dir(string(p)))
}

// Mkdir creates the directory Path, including any parent directories that
// need to be created along the way.
func (p Path) Mkdir() error {
	if p.Exists() {
		return fmt.Errorf("Cannot make directory %s because it already exists", p)
	}

	return os.MkdirAll(string(p), 0755)  // note umask will be applied
}

// WriteBytes writes the bytes to the Path.
func (p Path) WriteBytes(data []byte) error {
	outfile, err := os.Create(string(p))

	if err != nil {
		return err
	}

	defer outfile.Close()

	_, err = outfile.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// Unlink removes a file Path, but will return an error if the Path is a directory (see Rmdir).
func (p Path) Unlink() error {
	if p.IsDir() {
		return fmt.Errorf("%s is a directory.  Use Rmdir() instead.", p)
	}

	return os.Remove(string(p))
}

// Rmdir removes a directory, but will return an error if there are items within that directory (see RmdirRecursive).
func (p Path) Rmdir() error {
	if !p.IsDir() {
		return fmt.Errorf("%s is not a directory.  Use Unlink() instead.", p)
	}

	return os.Remove(string(p))
}

// RmdirRecursive removes a directory and all items within it.
func (p Path) RmdirRecursive() error {
	if !p.IsDir() {
		return fmt.Errorf("%s is not a directory.  Use Unlink() instead.", p)
	}

	return os.RemoveAll(string(p))
}

// Rename changes the name of the file to the target Path (essentially a move).
func (p Path) Rename(target Path) error {
	return os.Rename(string(p), string(target))
}

// OpenWithPermissions opens the Path with the specified mode and permissions.  If the Path does not exist, it creates it.
func (p Path) OpenWithPermissions(mode string, perms os.FileMode) (*os.File, error) {
	if p.IsDir() {
		return nil, fmt.Errorf("Cannot open %s because it is a directory.", p)
	}

	flag := os.O_RDONLY // default to read mode

	if strings.Contains(mode, "r") && strings.Contains(mode, "w") {
		flag = os.O_RDWR
	} else if strings.Contains(mode, "r") {
		flag = os.O_RDONLY
	} else if strings.Contains(mode, "w") {
		flag = os.O_WRONLY
	}

	if strings.Contains(mode, "+") {
		flag |= os.O_APPEND
	}

	if !p.Exists() {
		flag |= os.O_CREATE
	}

	return os.OpenFile(string(p), flag, perms)
}

// Open opens the Path with the specified mode with 0755 permissions.  If the Path does not exist, it creates it.
func (p Path) Open(mode string) (*os.File, error) {
	return p.OpenWithPermissions(mode, 0755)
}
