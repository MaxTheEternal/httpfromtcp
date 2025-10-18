package request

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"
)

var SEPERATOR = "\r\n"

type parserState int

const (
	StateInit parserState = 0
	StateDone parserState = 1
)

const BUFFERSIZE = 1028

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 0, BUFFERSIZE)
	// localBufferSize := BUFFERSIZE
	bufLen := 0
	req := NewRequest()
	for req.state != StateDone {
		// if len(buf) == localBufferSize {
		// 	localBufferSize = 2 * localBufferSize
		// 	newBuf := make([]byte, len(buf), localBufferSize)
		// 	copy(nwBuf, buf)
		// 	buf = newBuf
		// }
		n, err := reader.Read(buf[bufLen:])
		if err == io.EOF {
			req.state = StateDone
			break // End of file, break the loop
		}
		if err != nil {
			log.Fatal(err)
		}
		bufLen += n

		parsed, err := req.parse(buf)
		if err != nil {
			log.Fatal(err)
		}

		copy(buf, buf[parsed:bufLen])
		bufLen -= parsed

	}
	return req, nil
}

func NewRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == StateDone {
		return 0, errors.New("trying to parse read stuff")
	}
	req, read, err := parseRequestLine(string(data))
	if err != nil {
		return 0, err
	}

	if read == 0 {
		return 0, nil
	}

	r.state = StateDone
	r.RequestLine = *req
	return read, nil
}

func parseRequestLine(req string) (*RequestLine, int, error) {
	idx := strings.Index(req, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	line := req[:idx]

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return &RequestLine{}, 0, errors.New("Number of Arguments in Requenst Line not 3")
	}

	// Method Part
	for _, c := range parts[0] {
		if unicode.IsLower(c) || !unicode.IsLetter(c) {
			return &RequestLine{}, 0, errors.New("Method Wrong")
		}
	}

	// Version Part
	version, found := strings.CutPrefix(parts[2], "HTTP/")
	if !found || version != "1.1" {
		return nil, idx, errors.New("Version Wrong")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, idx, nil
}
