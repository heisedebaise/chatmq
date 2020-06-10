package chatmq

import (
	"bytes"
	"encoding/gob"
)

func encode(e interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(e); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

//Decode decode.
func Decode(b []byte, e interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(e)
}
