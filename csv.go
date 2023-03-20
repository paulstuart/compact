package compact

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CSVWriter interface {
	WriteCSV(io.Writer) error
}

type Unknown struct {
	Columns []string
	Rows    [][]string
	Index   map[string]int
}

func LoadUnknown(filename string, strip ...string) (Unknown, error) {
	var header []string
	var rows [][]string

	fn := func(ss []string) error {
		if len(header) == 0 {
			header = ss
			for i, head := range header {
				for _, prefix := range strip {
					head = strings.TrimPrefix(head, prefix)
				}
				header[i] = strings.Replace(head, ".", "_", -1)
			}
		} else {
			rows = append(rows, ss)
		}
		return nil
	}
	err := LoadCSV(filename, fn)
	if err != nil {
		return Unknown{}, err
	}

	m := make(map[string]int)
	for i, col := range header {
		m[col] = i
	}
	return Unknown{header, rows, m}, nil
}

func (u Unknown) Show(columns ...string) {
	idx := make([]int, len(columns))
	for i, col := range columns {
		idx[i] = u.Index[col]
	}

	for i, col := range idx {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Print(u.Columns[col])
	}
	fmt.Println()

	for _, row := range u.Rows {
		for i, col := range idx {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Print(row[col])
		}
		fmt.Println()
	}
}

func LineReader(r io.Reader, fn func(string) error) error {
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Text()
		if err := fn(line); err != nil {
			return err
		}
	}
	return nil
}

/*
func Restore(src, dest string) error {
	in, err := NewFileReader(src)
	if err != nil {
		return err
	}
	out, err := FileWriter(dest, 1<<16)
	if err != nil {
		in.Close()
		return err
	}
	fn := func(text string) error {
		var fp FirePersist
		if err := (&fp).Import(text); err != nil {
			return err
		}
		return fp.Save(out)
	}
	if err = LineReader(in, fn); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
*/

func LineCount(filename string) (int, error) {
	var count int
	return count, LoadLines(filename, func(_ string) error {
		count++
		return nil
	})
}

func LoadLines(filename string, fn func(string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: experiment with buffer sizes
	b := bufio.NewReader(f)
	r := io.Reader(b)

	// TODO: add lzma and others
	gzipped := filepath.Ext(filename) == ".gz"
	if gzipped {
		gzr, err := gzip.NewReader(b)
		if err != nil {
			return err
		}
		// TODO: does this need to get closed before f?
		defer gzr.Close()
		r = gzr
	}

	if gzipped {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("gunzip fail for %q -- %w", filename, err)
		}
		r = gzr
	}

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Text()
		if err = fn(line); err != nil {
			return err
		}
	}
	return nil
}

var errSampleComplete = errors.New("sample is complete")

// EstimateRecordCount samples the file to estimate how many lines
func EstimateRecordCount(filename string) (int, error) {
	fileSize, err := FileSize(filename)
	if err != nil {
		return 0, err
	}

	r, err := NewFileReader(filename)
	if err != nil {
		return 0, err
	}
	const sampleSize = 1000
	sum := 0
	counter := 0
	fn := func(s string) error {
		//log.Println("linelen:", len(s))
		sum += len(s)
		if counter++; counter >= sampleSize {
			return errSampleComplete
		}
		return nil
	}
	err = LineReader(r, fn)
	if err != nil && err != errSampleComplete {
		return 0, err
	}
	lineSize := sum / counter
	//log.Println("average line length:", lineSize)
	log.Printf("%d / %d = average line length: %d\n", sum, sampleSize, lineSize)
	// TODO: do test compression on collected data to understand compression ratio
	if filepath.Ext(filename) == ".gz" {
		fileSize *= 10 // rough guesstimate of compression ratio
	}
	return int(fileSize) / lineSize, nil
}

// LineLen will give stats on CSV record size for estimating storage needs
func LineLen(r io.Reader, count int, headers int) (LineStats, error) {
	var stats LineStats
	scan := bufio.NewScanner(r)
	var sum int64
	var idx int
	for scan.Scan() {
		idx++
		if idx < headers {
			continue // skip header
		}
		line := scan.Text()
		llen := len(line)
		if llen == 0 {
			continue
		}
		if stats.SampleSize >= count {
			break
		}
		stats.SampleSize++
		if llen > stats.Max {
			stats.Max = llen
		} else if stats.Min == 0 || stats.Min > llen {
			stats.Min = llen
		}
		sum += int64(llen)
	}
	stats.Avg = int(int64(stats.SampleSize) / sum)
	return stats, nil
}

func LoadCSV(filename string, fn func([]string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: experiment with buffer sizes
	b := bufio.NewReader(f)
	r := io.Reader(b)

	// TODO: add lzma and others
	gzipped := filepath.Ext(filename) == ".gz"
	if gzipped {
		gzr, err := gzip.NewReader(b)
		if err != nil {
			return fmt.Errorf("GZIP ERR (%T): %w", err, err)
			return err
		}
		// TODO: does this need to get closed before f?
		defer gzr.Close()
		r = gzr
	}
	/*
		if gzipped {
			gzr, err := gzip.NewReader(f)
			if err != nil {
				return fmt.Errorf("gunzip fail for %q -- %w", filename, err)
			}
			r = gzr
		}
	*/
	cr := csv.NewReader(r)
	for {
		records, err := cr.Read()
		if len(records) > 0 {
			if err := fn(records); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadCSVHeaders(filename string, skip int, fn func([]string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: experiment with buffer sizes
	b := bufio.NewReader(f)
	r := io.Reader(b)

	// TODO: add lzma and others
	gzipped := filepath.Ext(filename) == ".gz"
	if gzipped {
		gzr, err := gzip.NewReader(b)
		if err != nil {
			return fmt.Errorf("GZIP ERR (%T): %w", err, err)
			return err
		}
		// TODO: does this need to get closed before f?
		defer gzr.Close()
		r = gzr
	}
	/*
		if gzipped {
			gzr, err := gzip.NewReader(f)
			if err != nil {
				return fmt.Errorf("gunzip fail for %q -- %w", filename, err)
			}
			r = gzr
		}
	*/
	cr := csv.NewReader(r)
	count := 0
	for {
		records, err := cr.Read()
		count++
		if count <= skip {
			continue
		}
		if len(records) > 0 {
			if err := fn(records); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type FileInfo struct {
	Source    string // source of the data
	Bin       string // binary encoded version of the file
	Headers   int    // count of header lines in source
	Columns   []string
	SampleRow []string
	Stats     LineStats
}

func FileScan(filename string) (FileInfo, error) {
	info := FileInfo{
		Source: filename,
	}
	f, err := os.Open(filename)
	if err != nil {
		return info, fmt.Errorf("failed to open %q -- %w", filename, err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	var sum int64
	var idx int
	headCheck := true
	const count = 100 // let's get percent
	for scan.Scan() {
		idx++
		line := scan.Text()
		cols := strings.Fields(line)
		if strings.HasPrefix(line, "#") {
			info.Columns = cols
			info.Headers++
			continue // skip header
		}
		llen := len(line)
		if llen == 0 {
			continue
		}
		if _, err := strconv.ParseFloat(cols[0], 64); err != nil {
			if headCheck {
				info.Headers++
				continue
			}
		}
		if info.Stats.SampleSize >= count {
			break
		}
		headCheck = false
		info.Stats.SampleSize++
		if llen > info.Stats.Max {
			info.Stats.Max = llen
		} else if info.Stats.Min == 0 || info.Stats.Min > llen {
			info.Stats.Min = llen
		}
		// for i := range info.
		// sum += int64(llen)
	}
	info.Stats.Avg = int(int64(info.Stats.SampleSize) / sum)
	return info, nil
}

/*
Scan for header count, reset file
*/
type Recorder interface {
	GetRec() Record
}

// func NewRig()
func Makers(in string) Record {
	_, err := strconv.Atoi(in)
	if err == nil {
		var v NX32
		return &v
	}
	_, err = strconv.ParseFloat(in, 32)
	if err == nil {
		var v F32
		return &v
	}
	var v F32
	if err = v.Input(in); err != nil {
		log.Printf("input fail: %v", err)
	}
	return &v
}

// EvaluateFields establishs what type the columns represent
func EvaluateFields(filename string) (int, error) {
	fileSize, err := FileSize(filename)
	if err != nil {
		return 0, err
	}

	r, err := NewFileReader(filename)
	if err != nil {
		return 0, err
	}
	const sampleSize = 1000
	sum := 0
	counter := 0
	fn := func(s string) error {
		//log.Println("linelen:", len(s))
		sum += len(s)
		if counter++; counter >= sampleSize {
			return errSampleComplete
		}
		return nil
	}
	err = LineReader(r, fn)
	if err != nil && err != errSampleComplete {
		return 0, err
	}
	lineSize := sum / counter
	//log.Println("average line length:", lineSize)
	log.Printf("%d / %d = average line length: %d\n", sum, sampleSize, lineSize)
	// TODO: do test compression on collected data to understand compression ratio
	if filepath.Ext(filename) == ".gz" {
		fileSize *= 10 // rough guesstimate of compression ratio
	}
	return int(fileSize) / lineSize, nil
}
