package blueprint

//go:generate esc -private -o drafter.go -pkg blueprint ext/drafter/bin

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/bukalapak/snowboard/api"
	snowboard "github.com/bukalapak/snowboard/parser"
)

type parser struct{}

func GetStructure(apibFile []byte) (*api.API, error) {
	engine := parser{}

	bp, err := snowboard.Parse(bytes.NewReader(apibFile), engine)
	if err != nil {
		return nil, err
	}

	return bp, nil
}

func (p parser) Parse(r io.Reader) ([]byte, error) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	err := p.exec(r, &stdOut, &stdErr, "--format", "json")
	if err != nil {
		if stdErr.String() != "" {
			return nil, errors.New(stdErr.String())
		}
		return nil, err
	}

	return stdOut.Bytes(), nil
}

func (p parser) Validate(r io.Reader) ([]byte, error) {
	return nil, errors.New("This is a parse-only engine")
}

func (p parser) Version() string {
	return ""
}

func (p parser) exec(r io.Reader, stdOut io.Writer, stdErr io.Writer, args ...string) error {
	exe, err := tmpCommand()
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(exe)
	}()

	cmd := exec.Command(exe, args...)
	cmd.Stdin = r
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	return cmd.Run()
}

func tmpCommand() (string, error) {
	b := embeddedDrafter()
	tmp, err := ioutil.TempFile(os.TempDir(), "drafter")
	if err != nil {
		return "", err
	}

	tmp.Write(b)
	tmp.Chmod(0700)
	tmp.Close()

	return tmp.Name(), nil
}

func embeddedDrafter() []byte {
	return _escFSMustByte(false, "/ext/drafter/bin/drafter")
}
