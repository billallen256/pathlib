# pathlib

Simple `Path` type and helper methods for Go in the spirit of Python's [pathlib](https://docs.python.org/3/library/pathlib.html).

In addition to providing many convenience methods, using the Path type provides more type safety than simply passing strings around.

```go
myPath := Path("/foo/bar/baz.boo")
```
