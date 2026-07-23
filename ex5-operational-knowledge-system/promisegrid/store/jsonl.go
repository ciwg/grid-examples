package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

const MaxJSONLLineBytes = 1 << 20

func OpenAppendOnlyLog(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
}

// Intent: Keep append-only replay mechanics substrate-owned so runtimes and
// relays can share the same durable JSONL read path without ex5-specific
// workflow code owning the scanner semantics. Source: DI-lemor
func ReadJSONL[T any](file *os.File) (records []T, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	// Intent: Keep replay within the large-request envelope ex5 already relied
	// on, so durable large revisions stay readable after restart while the
	// scanner logic moves into reusable substrate. Source: DI-busor; DI-lemor
	scanner.Buffer(make([]byte, 64*1024), MaxJSONLLineBytes)
	records = []T{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record T
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

// Intent: Keep append-only write and fsync behavior substrate-owned so local
// runtimes and relays share one durable JSONL append contract. Source: DI-lemor
func AppendJSONL(file *os.File, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(body, '\n')); err != nil {
		return err
	}
	return file.Sync()
}
