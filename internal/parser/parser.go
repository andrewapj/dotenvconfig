package parser

import (
	"bytes"
	"errors"
)

func Parse(data []byte) (map[string]string, error) {
	lines := bytes.Split(data, []byte("\n"))
	result := make(map[string]string)

	for _, line := range lines {
		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) != 2 {
			return nil, errors.New("error parsing line: " + string(line))
		}
		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))
		result[key] = value
	}
	return result, nil
}
