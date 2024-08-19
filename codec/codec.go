package codec

import "fmt"

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type Codec interface {
	Marshaler
	Unmarshaler
}

func Marshal(v any) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}

	if m, ok := v.(Marshaler); ok {
		return m.Marshal()
	}

	return nil, fmt.Errorf("unkown data")
}
