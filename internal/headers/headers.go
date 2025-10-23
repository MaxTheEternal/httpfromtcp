package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	Headers map[string]string
}

var COULDNT_PARSE_HEADERS = fmt.Errorf("couldnt parse header line")

func NewHeaders() Headers {
	return Headers{
		Headers: map[string]string{},
	}
}

func (h Headers) Get(key string) string {
	return h.Headers[strings.ToLower(key)]
}

func (h Headers) Set(key string, value string) {
	lowerKey := strings.ToLower(key)

	if v, ok := h.Headers[lowerKey]; !ok {
		h.Headers[lowerKey] = value
	} else {
		h.Headers[lowerKey] = v + ", " + value
	}
	fmt.Printf("Added value: %s\n", value)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	readBytes := 0
	done := false

	for {
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

		valid := validKeyToken(name)
		if valid != nil {
			return readBytes, false, valid
		}

		h.Set(string(name), string(value))
		readBytes += idx + 2
		data = data[idx+2:]

	}
	return readBytes, done, nil
}

func validKeyToken(key []byte) error {
	if len(key) < 1 {
		return fmt.Errorf("FieldName cant be smaller than 1 Char")
	}
	for _, ch := range key {
		found := false
		if ch >= 'a' && ch <= 'z' {
			found = true
		}
		if ch >= 'A' && ch <= 'Z' {
			found = true
		}
		if ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if !found {
			return fmt.Errorf("invalid Token in Key")
		}
	}
	return nil
}
