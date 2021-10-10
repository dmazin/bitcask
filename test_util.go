package bitcask

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func CopyFile(srcFilepath, dstFilepath string) error {
	data, err := os.ReadFile(srcFilepath)
	if err != nil {
		return err
	}

	err = os.WriteFile(dstFilepath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func SuppressLogs(tb testing.TB) {
	flags := log.Flags()
	output := log.Writer()

	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	tb.Cleanup(func() {
		log.SetFlags(flags)
		log.SetOutput(output)
	})
}
