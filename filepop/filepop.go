package filepop

import (
	"bytes"
	//"fmt"
	"io"
	"os"
)

// Function to delete one line from file.
func popLine(f *os.File) ([]byte, error) {

	// Stating file for file metadata.
	fi, err := f.Stat()
	check(err)

	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, io.SeekStart)
	check(err)

	_, err = io.Copy(buf, f)
	check(err)

	line, err := buf.ReadBytes('\n')
	check(err)

	_, err = f.Seek(0, io.SeekStart)
	check(err)

	nw, err := io.Copy(f, buf)
	check(err)

	err = f.Truncate(nw)
	check(err)

	err = f.Sync()
	check(err)

	_, err = f.Seek(0, io.SeekStart)
	check(err)

	return line, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Pop function deletes one line from the beginning of the file.
func Pop(fname string) {
	f, err := os.OpenFile(fname, os.O_RDWR, 0666)
	check(err)
	defer f.Close()
	_, err = popLine(f)
	check(err)
}
