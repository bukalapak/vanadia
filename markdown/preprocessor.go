package markdown

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

var includePattern = regexp.MustCompile(`<!--\s*include\s+([^;]+)(?:;\s*(.*)?)?\s+-->`)

func Preprocess(r io.Reader, dir string) ([]byte, error) {
	buf := new(bytes.Buffer)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Bytes()
		if includePattern.Match(line) {
			buf.Write(resolveInclude(dir, line))
		} else {
			buf.Write(line)
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}

func resolveInclude(dir string, line []byte) []byte {
	submatch := includePattern.FindSubmatch(line)
	filename := string(submatch[1])
	file, err := ioutil.ReadFile(path.Join(dir, filename))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error including %s: %s\n", filename, err.Error())
		return nil
	}
	substitutions := bytes.Split(submatch[2], []byte(";"))
	for _, sub := range substitutions {
		parts := bytes.Split(sub, []byte("="))
		if len(parts) < 2 {
			continue
		}
		file = bytes.Replace(file, parts[0], parts[1], -1)
	}
	return file
}
