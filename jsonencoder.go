package main

import "sync"

var _jsonPOOL = sync.Pool{New: func() any {
	return NewJSONEncoder()
}}

type JSONEncoder struct {
	b []byte
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		b: make([]byte, 0, 1024),
	}
}

const (
	NewLineCharacter = '\n'
	TabCharacter     = '\t'
)

func (j *JSONEncoder) Encode(rec Record) ([]byte, error) {
	j.b = append(j.b, '{')
	j.b = append(j.b, NewLineCharacter)
	j.b = append(j.b, TabCharacter)
	j.addKeyValue("message", rec.Message)
	j.b = append(j.b, NewLineCharacter)
	j.b = append(j.b, '}')

	res := j.b
	j.b = j.b[:0]

	return res, nil
}

func (j *JSONEncoder) addKeyValue(key string, value string) {

	j.addString(key)
	j.b = append(j.b, ':')
	j.b = append(j.b, ' ')
	j.addString(value)

}

func (j *JSONEncoder) addString(str string) {
	j.b = append(j.b, '"')
	j.b = append(j.b, []byte(str)...)
	j.b = append(j.b, '"')
}
