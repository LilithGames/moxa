package utils

import (
	"google.golang.org/protobuf/proto"
)

func Clone[T proto.Message](src T) T {
	return proto.Clone(src).(T)
}
