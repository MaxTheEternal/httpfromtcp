package request

import (
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/MaxTheEternal/httpfromtcp/internal/headers"
)

var SEPERATOR = "\r\n"

type parserState int

const (
	StateInit parserState = iota
	StateDone
	StateParsingHeaders
)

const BUFFERSIZE = 1028

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func NewRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFFERSIZE)
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
			return nil, err
		}
		bufLen += n

		parsed, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[parsed:bufLen])
		bufLen -= parsed
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateDone:
			break outer
		case StateInit:
			rl, n, err := parseRequestLine(string(data[read:]))
			if err != nil {
				r.state = StateDone
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n + 2
			r.state = StateParsingHeaders
		case StateParsingHeaders:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			read += n
			if done {
				r.state = StateDone
			}

		}
	}
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
