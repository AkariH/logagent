package tailfile

import (
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
