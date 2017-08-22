package main

import (
	"bytes"
	"io"
	"strings"
)

// Parses a json string. Stops at handlebars. Everything read is written. The
// result is a runtime templater.
type Parser struct {
	source io.ByteReader
	sink   WriterByter
}

type WriterByter interface {
	io.ByteWriter
	io.Writer
}

func NewParser(source io.ByteReader, sink WriterByter) Parser {
	return Parser{source, sink}
}

// Moves to the next handlebar node and returns the text.
func (p *Parser) Next() (string, bool, error) {

find:
	err := p.findLeftCurlies()
	if err != nil {
		if err == io.EOF {
			return "", false, nil
		}
		return "", false, err
	}

	cut, err := p.cutTillRightCurlies()
	if err != nil {
		// if hit EOF before right curlies
		if err == io.EOF {
			_, err = p.sink.Write([]byte("{{" + cut))
			if err != nil {
				return "", false, err
			}
		}
		return "", false, nil
	}

	if strings.Contains(cut, "\n") {
		_, err = p.sink.Write([]byte("{{" + cut + "}}"))
		if err != nil {
			return "", false, err
		}
		goto find
	}

	return cut, true, nil
}

func (p *Parser) findLeftCurlies() error {

	err := rwUntil(p.source, p.sink, '{')
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) cutTillRightCurlies() (string, error) {

	var buf bytes.Buffer
	err := rwUntil(p.source, &buf, '}')
	if err != nil {
		return string(buf.Bytes()), err
	}

	return string(buf.Bytes()), nil
}

func rwUntil(source io.ByteReader, sink io.ByteWriter, delimeter byte) error {
	for {

		var b byte
		var err error

		b, err = source.ReadByte()
		if err != nil {
			return err
		}

		if b != delimeter {
			err = sink.WriteByte(b)
			if err != nil {
				return err
			}
			continue
		}

		b, err = source.ReadByte()
		if err != nil {
			return err
		}

		if b != delimeter {
			err = sink.WriteByte(delimeter)
			if err != nil {
				return err
			}
			err = sink.WriteByte(b)
			if err != nil {
				return err
			}
			continue
		}

		break
	}
	return nil
}
