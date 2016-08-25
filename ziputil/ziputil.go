package ziputil

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Walkable represents a walkable filesystem-like type.
type Walkable interface {
	http.FileSystem
	Walk(path string, walkFn filepath.WalkFunc) error
}

type localfs struct{}

func (localfs) Walk(path string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(path, walkFn)
}

func (localfs) Open(path string) (http.File, error) {
	return os.Open(path)
}

// ZipPaths writes the provided paths to the given Writer.
//
// If walkable is nil the local filesystem will be used.
func ZipPaths(out io.Writer, paths []string, walkable Walkable) error {
	if walkable == nil {
		walkable = &localfs{}
	}
	w := zip.NewWriter(out)
	defer w.Close()
	for _, path := range paths {
		if err := addToZip(w, path, walkable); err != nil {
			return err
		}
	}
	return nil
}

func addToZip(w *zip.Writer, path string, walkable Walkable) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return walkable.Walk(path, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fh.Name = fp
		if fi.IsDir() {
			if path == fp {
				return nil
			}
			fh.Name += "/"
		} else {
			fh.Method = zip.Deflate
		}
		wr, err := w.CreateHeader(fh)
		if err != nil {
			return err
		}
		if fh.Mode().IsDir() {
			return nil
		}
		if fh.Mode().IsRegular() {
			f, err := walkable.Open(fp)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(wr, f); err != nil {
				return err
			}
		}
		return nil
	})
}
