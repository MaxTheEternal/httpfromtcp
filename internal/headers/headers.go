package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var COULDNT_PARSE_HEADERS = fmt.Errorf("couldnt parse header line")

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	readBytes := 0
	done := false

	for {
		data = data[readBytes:]
		idx := bytes.Index(data, []byte("\r\n"))
		if idx == -1 {
			break
		}
		if idx == 0 {
			readBytes += 2
			done = true
			break
		}

		headerLine := data[:idx]
		name, value, found := bytes.Cut(headerLine, []byte(":"))
		if !found {
			return readBytes, false, COULDNT_PARSE_HEADERS
		}

		if name[len(name)-1] == byte(' ') {
			return readBytes, false, COULDNT_PARSE_HEADERS
		}
		name = bytes.TrimSpace(name)
		value = bytes.TrimSpace(value)

		h[string(name)] = string(value)
		readBytes += idx + 2

	}
	return readBytes, done, nil
}
