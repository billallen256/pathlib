package pathlib

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func randomString(length int) string {
	var builder strings.Builder
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		c := string(chars[rand.Intn(len(chars))])
		builder.WriteString(c)
	}

	return builder.String()
}

func TestResolve(t *testing.T) {
	p := Path("/etc/../etc/passwd")
	resolved, err := p.Resolve()

	if err != nil {
		t.Errorf(err.Error())
	}

	if resolved != Path("/etc/passwd") {
		t.Errorf("Path did not resolve correctly: %s", resolved)
	}
}

func TestNoResolve(t *testing.T) {
	p := Path("fjdksafdsljakfldsjf")
	_, err := p.Resolve()

	if err == nil {
		t.Errorf("Resolve should fail for non-existent paths")
	}
}

func TestExists(t *testing.T) {
	p := Path("/usr")

	if !p.Exists() {
		t.Errorf("Path should exist")
	}
}

func TestNotExists(t *testing.T) {
	p := Path("/foo/bar/baz/fjkdsalfjaklrejakfdsa")

	if p.Exists() {
		t.Errorf("Path should not exist")
	}
}

func TestReadBytes(t *testing.T) {
	p := Path("path.go")
	content, err := p.ReadBytes()

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(content) == 0 {
		t.Errorf("Received zero bytes")
	}
}

func TestIsDir(t *testing.T) {
	if !Path("/usr").IsDir() {
		t.Errorf("Path should be a directory")
	}
}

func TestIsNotDir(t *testing.T) {
	if Path("/etc/passwd").IsDir() {
		t.Errorf("Path should not be a directory")
	}
}

func TestIsFile(t *testing.T) {
	if !Path("/etc/passwd").IsFile() {
		t.Errorf("Path should be a file")
	}
}

func TestIsNotFile(t *testing.T) {
	if Path("/usr").IsFile() {
		t.Errorf("Path should not be a file")
	}
}

func TestGlob(t *testing.T) {
	paths, err := Path("/etc").Glob("*.conf")

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(paths) == 0 {
		t.Errorf("/etc/*.conf returned no results")
	}

	for _, path := range paths {
		absPath, err := path.Resolve()

		if err != nil {
			t.Errorf(err.Error())
		}

		if path != absPath {
			t.Errorf("Glob should only return absolute paths")
		}
	}
}

func suffixTest(orig, newSuffix, target string, t *testing.T) {
	changed := Path(orig).WithSuffix(newSuffix)

	if changed != Path(target) {
		t.Errorf("WithSuffix failed: %s != %s", changed, target)
	}
}

func TestWithSuffix(t *testing.T) {
	suffixTest("foo.bar", "baz", "foo.baz", t)
	suffixTest("foo", "baz", "foo.baz", t)
	suffixTest("/foo/bar/baz.a", "baz", "/foo/bar/baz.baz", t)
	suffixTest("foo/bar/baz.a", "baz", "foo/bar/baz.baz", t)
	suffixTest("/foo/bar.a/baz.zip", "", "/foo/bar.a/baz", t)
	suffixTest("foo/bar.a/baz.zip", "", "foo/bar.a/baz", t)
}

func TestTouch(t *testing.T) {
	p := Path("/tmp/pathlib-" + randomString(20))
	p.Touch()

	if !p.Exists() {
		t.Errorf("Touch failed for path %s", p)
	}
}

func TestAge(t *testing.T) {
	p := Path("/tmp/pathlib-" + randomString(20))
	p.Touch()
	time.Sleep(time.Duration(2) * time.Second)
	age, err := p.Age(time.Now())

	if err != nil {
		t.Errorf(err.Error())
	}

	if age < time.Duration(2)*time.Second || age > time.Duration(2500)*time.Millisecond {
		t.Errorf("Received incorrect age of %s", age)
	}
}

func TestJoinPath(t *testing.T) {
	tests := map[Path]Path{
		Path("/tmp/").JoinPath(Path("foo"), Path("bar")): Path("/tmp/foo/bar"),
		Path("/tmp").JoinPath(Path("foo/bar")):           Path("/tmp/foo/bar"),
		Path("foo").JoinPath(Path("bar")):                Path("foo/bar"),
		Path("foo/bar").JoinPath(Path("baz")):            Path("foo/bar/baz"),
	}

	for test, target := range tests {
		if test != target {
			t.Errorf("JoinPath failed: %s != %s", test, target)
		}
	}
}

func TestParent(t *testing.T) {
	p := Path("/var/log")
	parent := p.Parent()
	target := Path("/var")

	if parent != target {
		t.Errorf("Parent failed: %s != %s", parent, target)
	}
}

func TestName(t *testing.T) {
	tests := map[string]string{
		Path("/var/log/messages").Name(): "messages",
		Path("foo").Name():               "foo",
		Path("foo/bar.baz").Name():       "bar.baz",
	}

	for test, target := range tests {
		if test != target {
			t.Errorf("Name failed: %s != %s", test, target)
		}
	}
}
