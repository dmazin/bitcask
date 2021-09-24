package naivedb

import "os"

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
