package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	req := strings.Split(string(data), "\r\n")
	reqLine, err := parseRequestLine(req[0])
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return &RequestLine{}, errors.New("Number of Arguments in Requenst Line not 3")
	}

	// Method Part
	for _, c := range parts[0] {
		if unicode.IsLower(c) || !unicode.IsLetter(c) {
			return &RequestLine{}, errors.New("Method Wrong")
		}
	}

	// Version Part
	version, found := strings.CutPrefix(parts[2], "HTTP/")
	if !found || version != "1.1" {
		return &RequestLine{}, errors.New("Version Wrong")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, nil
}
