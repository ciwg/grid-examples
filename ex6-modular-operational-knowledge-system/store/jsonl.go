package store

import (
	"bufio"
	"errors"
	"os"
)

const maxJSONLLineBytes = 1 << 20

func openAppendOnlyLog(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
}

func readLines(file *os.File) (lines [][]byte, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxJSONLLineBytes)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		copied := make([]byte, len(line))
		copy(copied, line)
		lines = append(lines, copied)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func appendLine(file *os.File, line []byte) error {
	if _, err := file.Write(append(append([]byte{}, line...), '\n')); err != nil {
		return err
	}
	return file.Sync()
}
