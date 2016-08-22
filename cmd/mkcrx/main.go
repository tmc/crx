// Program mkcrx packages a directory as a Chrome crx extension.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tmc/crx"
	"github.com/tmc/crx/ziputil"
)

var flagKey = flag.String("key", "key.pem", "path to private key in PEM format")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [flags] <directory-name>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}
func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := mkcrx(*flagKey, flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mkcrx(keyPath, dirPath string) error {
	baseName := filepath.Base(dirPath)
	keyFile, err := os.Open(keyPath)
	if err != nil {
		return err
	}
	crxFile, err := os.Create(baseName + ".crx")
	if err != nil {
		return err
	}
	defer crxFile.Close()

	crxContents, err := crx.FromPath(dirPath, keyFile, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(crxFile, crxContents)
	return err
}

func mkzip(zipPath, dirPath string) error {
	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()
	return ziputil.ZipPaths(zf, []string{dirPath}, nil)
}
