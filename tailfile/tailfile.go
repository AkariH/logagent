package tailfile

import (
	"fmt"
	"github.com/nxadm/tail"
	"logagent/utils/logger"
)

var (
	Tailfile *tail.Tail
	err      error
)

func Init(filename string) error {

	cfg := tail.Config{
		Follow:    true,
		ReOpen:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	Tailfile, err = tail.TailFile(filename, cfg)
	logger.Infof("tail file %s success", filename)

	if err != nil {
		return err
	}
	return nil
}

func GetTailFile() *tail.Tail {
	return Tailfile
}

func ReadLine() (*tail.Line, error) {
	line, ok := <-Tailfile.Lines
	if !ok {
		// Log the warning about the tail file being closed
		logger.Warnf("Tail file closed, filename: %s", Tailfile.Filename)

		return nil, fmt.Errorf("tail file closed, filename: %s", Tailfile.Filename)
	}
	logger.Infof("line: %s", line.Text)
	return line, nil
}
