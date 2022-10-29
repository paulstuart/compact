package compact

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz/lzma"
)

type LineStats struct {
	SampleSize int
	Min        int
	Max        int
	Avg        int
}

var ErrNotFound = errors.New("not found")

func FileSize(filename string) (int64, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}

func Exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func AbsPath(filename string) string {
	abs, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalln("abs path error:", err)
	}
	return abs
}

type FWriter struct {
	f *os.File
	b *bufio.Writer
	g *gzip.Writer
	w io.Writer
}

func (fw *FWriter) Close() error {
	if fw.g != nil {
		err := fw.g.Close()
		if err != nil {
			fw.b.Flush()
			fw.f.Close()
			return err
		}
	}
	fw.b.Flush()
	return fw.f.Close()
}

func (fw *FWriter) Write(b []byte) (int, error) {
	return fw.w.Write(b)
}

type FileBuf struct {
	w        *bufio.Writer
	f        *os.File
	filename string
}

func NewFileBuf(filename string, size int) (*FileBuf, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("can't create %q -- %w", filename, err)
	}
	w := bufio.NewWriterSize(f, size)
	buff := &FileBuf{
		w:        w,
		f:        f,
		filename: filename,
	}
	return buff, nil
}

func (fb *FileBuf) Write(p []byte) (int, error) {
	return fb.w.Write(p)
}

func (fb *FileBuf) Close() error {
	if err := fb.w.Flush(); err != nil {
		if err2 := fb.f.Close(); err2 != nil {
			log.Printf("also failed to close file: %v", err2)
		}
		return fmt.Errorf("error flushing file %q -- %w", fb.filename, err)
	}
	if err := fb.f.Close(); err != nil {
		return fmt.Errorf("error closing file %q -- %w", fb.filename, err)
	}
	return nil
}

type devNull struct{}

func (d devNull) Write(b []byte) (int, error) {
	return len(b), nil
}
func (d devNull) Close() error {
	return nil
}

type stdOut struct{}

func (d stdOut) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}
func (d stdOut) Close() error {
	return nil
}

// FileWriter creates in io.WriteCloser that buffers
// and optionally compresses the output
func FileWriter(filename string, bufSize int) (io.WriteCloser, error) {
	switch filename {
	case "/dev/null":
		return devNull{}, nil
	case "/dev/stdout", "-":
		return stdOut{}, nil
	}

	out, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	b := bufio.NewWriterSize(out, bufSize)

	fw := &FWriter{
		f: out,
		b: b,
		w: b,
	}
	if filepath.Ext(filename) == ".gz" {
		fw.g = gzip.NewWriter(fw.b)
		fw.w = fw.g
	}
	return fw, nil
}

type FileReader struct {
	f *os.File
	b *bufio.Reader
	g *gzip.Reader
	r io.Reader
}

func (br FileReader) Close() error {
	if br.g != nil {
		br.g.Close()
	}
	return br.f.Close()
}

func (br FileReader) Read(buf []byte) (int, error) {
	return br.r.Read(buf)
}

func NewFileReader(filename string) (io.ReadCloser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	b := bufio.NewReader(f)
	br := &FileReader{
		f: f,
		b: b,
		r: b,
	}

	switch filepath.Ext(filename) {
	case ".gz":
		br.g, err = gzip.NewReader(br.b)
		if err != nil {
			log.Printf("ERR (%T): %v", err, err)
			return nil, err
		}
		br.r = br.g
	case ".zz", ".lzma":
		br.r, err = lzma.NewReader(br.b)
		if err != nil {
			return nil, err
		}
	}
	return br, nil
}
