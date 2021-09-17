package bitcask

import (
	"fmt"
	"os"
)

func set(key string, value string) error {
	f, err := os.OpenFile("database", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(fmt.Sprintf("%s,%s\n", key, value)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		return err
	}

	return f.Close()
}

func get(key string) (string, error) {
	f, err := os.OpenFile("database", os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
}
