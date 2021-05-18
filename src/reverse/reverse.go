package reverse

import (
	"os"
)

type ReverseReader struct {
	bytesRead []byte
}

func NewReverseReader(file *os.File) (*ReverseReader, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := fileInfo.Size()

	bytesToRead := 5
	bytesRead := make([]byte, bytesToRead)

	file.ReadAt(bytesRead, size - int64(bytesToRead))

	return &ReverseReader{
		bytesRead: bytesRead,
	}, nil
}
