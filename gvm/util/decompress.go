package util

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func closeHelper(c io.Closer) { _ = c.Close() }

func destHelper(src, suffix string) string {
	var destName = filepath.Base(src)
	destName = destName[:len(destName)-len(suffix)]
	dest := filepath.Join(filepath.Dir(src), destName)
	return dest
}

func ExtractTarGz(src string) (dest string, err error) {
	var file *os.File
	file, err = os.Open(src)
	if err != nil {
		return
	}
	defer closeHelper(file)

	var gzipReader *gzip.Reader
	gzipReader, err = gzip.NewReader(file)
	if err != nil {
		return
	}
	defer closeHelper(gzipReader)

	dest = destHelper(src, ".tar.gz")
	tarReader := tar.NewReader(gzipReader)

	for {
		header, e := tarReader.Next()
		if e == io.EOF {
			break
		}
		fullPath := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if e := os.MkdirAll(fullPath, os.FileMode(header.Mode)); e != nil {
				return
			}
		case tar.TypeReg:
			writer, e := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if e != nil {
				return
			}
			if _, e = io.Copy(writer, tarReader); e != nil {
				return
			}
			_ = writer.Close()
		}
	}
	return
}

func ExtractZip(src string) (dest string, err error) {
	dest = destHelper(src, ".zip")

	var writeFile = func(destination string, f *zip.File) error {
		rc, e := f.Open()
		if e != nil {
			return e
		}
		defer closeHelper(rc)

		destinationPath := filepath.Join(destination, f.Name)
		if !strings.HasPrefix(destinationPath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", destinationPath)
		}

		if f.FileInfo().IsDir() {
			e = os.MkdirAll(destinationPath, f.Mode())
			if e != nil {
				return e
			}
		} else {
			e = os.MkdirAll(filepath.Dir(destinationPath), 0755)
			if e != nil {
				return e
			}

			var file *os.File
			file, e = os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if e != nil {
				return e
			}
			defer closeHelper(file)

			_, e = io.Copy(file, rc)
			if e != nil {
				return e
			}
		}

		return nil
	}

	var r *zip.ReadCloser
	r, err = zip.OpenReader(src)
	if err != nil {
		return
	}

	defer closeHelper(r)

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return
	}

	for _, f := range r.File {
		err = writeFile(dest, f)
		if err != nil {
			return
		}
	}

	return
}
