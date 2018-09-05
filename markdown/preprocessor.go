package markdown

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

var includePattern = regexp.MustCompile(`<!--\s*include\s+(.+?)\s*-->`)

func Preprocess(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		if includePattern.Match(bytes) {
			filename := string(includePattern.FindSubmatch(bytes)[1])
			file, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error including %s: %s\n", filename, err.Error())
				continue
			}
			buf.Write(file)
		} else {
			buf.Write(bytes)
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}