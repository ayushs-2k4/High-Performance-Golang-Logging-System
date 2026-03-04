package main

import (
	"reflect"
	"strconv"
	"sync"
)

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
	MessageKey       = "message"
)

func (j *JSONEncoder) Encode(rec Record) ([]byte, error) {
	j.b = append(j.b, '{')
	j.addCharacter(NewLineCharacter)
	j.addCharacter(TabCharacter)
	j.addKeyValue(MessageKey, Value{
		val:     rec.Message,
		valType: reflect.String,
	})

	for _, kv := range rec.KVs {
		j.b = append(j.b, ',')
		key := kv.Key
		val := kv.Value

		j.addCharacter(NewLineCharacter)
		j.addCharacter(TabCharacter)

		j.addKeyValue(key, *val)

	}

	j.addCharacter(NewLineCharacter)
	j.b = append(j.b, '}')

	res := j.b
	j.b = j.b[:0]

	return res, nil
}

func (j *JSONEncoder) addCharacter(c rune) {
	j.b = append(j.b, byte(c))
}

func (j *JSONEncoder) addKeyValue(key string, value Value) {
	j.addString(key)
	j.b = append(j.b, ':')
	j.b = append(j.b, ' ')

	switch value.valType {
	case reflect.String:
		j.addString(value.val.(string))
	case reflect.Int:
		j.addInt(value.val.(int))
	}

}

func (j *JSONEncoder) addString(str string) {
	j.b = append(j.b, '"')
	j.b = append(j.b, str...)
	j.b = append(j.b, '"')
}

func (j *JSONEncoder) addInt(val int) {
	j.b = strconv.AppendInt(j.b, int64(val), 10)
}
