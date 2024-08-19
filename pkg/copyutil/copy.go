package copyutil

import (
	"encoding/json"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

func Clone(src interface{}) (dst interface{}) {
	if src == nil {
		return
	}

	srcVal := reflect.ValueOf(src)
	if !srcVal.IsValid() {
		return
	}

	var dstValue reflect.Value
	kd := srcVal.Kind()

	isMap := false
	isStruct := false

	switch kd {
	case reflect.Ptr, reflect.Uintptr:
		elem := srcVal.Elem()
		if !elem.IsValid() {
			return
		}
		dstValue = reflect.New(elem.Type())
		dst = dstValue.Interface()
	case reflect.Map, reflect.Array, reflect.Slice:
		dstValue = reflect.New(srcVal.Type())
		dst = dstValue.Interface()
		isMap = true
	case reflect.Struct:
		dstValue = reflect.New(srcVal.Type())
		dst = dstValue.Interface()
		isStruct = true
	default:
		dst = src
		return
	}

	switch src.(type) {
	case proto.Marshaler:
		pbSrc := src.(proto.Marshaler)
		buf, _ := pbSrc.Marshal()
		_ = dst.(proto.Unmarshaler).Unmarshal(buf)
	default:
		buf, _ := json.Marshal(src)
		_ = json.Unmarshal(buf, dst)
	}

	if isMap || isStruct {
		dst = dstValue.Elem().Interface()
	}

	return
}
