package pathlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Path string

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

func (p Path) JoinPath(paths ...Path) Path {
	ret := string(p)

	for _, path := range paths {
		ret = filepath.Join(ret, string(path))
	}

	return Path(ret)
}

func (p Path) Name() string {
	return filepath.Base(string(p))
}

func (p Path) Parent() Path {
	return Path(filepath.Dir(string(p)))
}

func (p Path) Mkdir() error {
	if p.Exists() {
		return fmt.Errorf("Cannot make directory %s because it already exists", p)
	}

	parent := p.Parent()
	perms, err := parent.Permissions()

	if err != nil {
		return err
	}

	// Copy source permissions, making sure that the destination is
	// at least readable, writable, and executable
	perms = perms | 0700

	return os.Mkdir(string(p), perms)
}

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

func (p Path) Unlink() error {
	if p.IsDir() {
		return fmt.Errorf("%s is a directory.  Use Rmdir() instead.", p)
	}

	return os.Remove(string(p))
}

func (p Path) Rmdir() error {
	if !p.IsDir() {
		return fmt.Errorf("%s is not a directory.  Use Unlink() instead.", p)
	}

	return os.Remove(string(p))
}

func (p Path) RmdirRecursive() error {
	if !p.IsDir() {
		return fmt.Errorf("%s is not a directory.  Use Unlink() instead.", p)
	}

	return os.RemoveAll(string(p))
}

func (p Path) Rename(target Path) error {
	return os.Rename(string(p), string(target))
}
