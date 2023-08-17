package parser

import (
	"bytes"
	"errors"
)

func Parse(data []byte) (map[string]string, error) {
	lines := bytes.Split(data, []byte("\n"))
	result := make(map[string]string)

	for _, line := range lines {
		trimmedLine := bytes.TrimSpace(line)

		// Skip empty lines and lines starting with #
		if len(trimmedLine) == 0 || trimmedLine[0] == '#' {
			continue
		}

		parts := bytes.SplitN(trimmedLine, []byte("="), 2)
		if len(parts) != 2 {
			return nil, errors.New("error parsing line: " + string(trimmedLine))
		}

		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))
		result[key] = value
	}

	return result, nil
}
