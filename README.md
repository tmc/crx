

# crx
`import "github.com/tmc/crx"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package crx contains tools for producing Chrome crx files.


## <a name="pkg-index">Index</a>
* [func CRXFromPath(path string, rsaKey io.Reader, walkable ziputil.Walkable) (io.Reader, error)](#CRXFromPath)
* [func WriteCRXFromZip(w io.Writer, zipContents io.Reader, rsaKey io.Reader) error](#WriteCRXFromZip)


## <a name="CRXFromPath">func</a> [CRXFromPath](/src/target/crx.go?s=420:513#L15)
``` go
func CRXFromPath(path string, rsaKey io.Reader, walkable ziputil.Walkable) (io.Reader, error)
```
CRXFromPath returns a reader that contains a CRX file.

If walkable is nil it will fall back to the local filesystem.



## <a name="WriteCRXFromZip">func</a> [WriteCRXFromZip](/src/target/crx.go?s=837:917#L25)
``` go
func WriteCRXFromZip(w io.Writer, zipContents io.Reader, rsaKey io.Reader) error
```
WriteCRXFromZip writes to the given writer a crx file that is described by zipContents.
