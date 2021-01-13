package util

import (
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// WrapUint64 function returns nil when an error occurred otherwise returns value in wrappers.UInt64Value format.
func WrapUint64(value string) *wrappers.UInt64Value {
	if value != "" {
		var i uint64
		i, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			return &wrappers.UInt64Value{Value: i}
		}
	}

	return nil
}

// WrapTime function returns nil when an error occurred or time is empty
// otherwise returns time in timestamp.Timestamp format.
func WrapTime(t *time.Time) *timestamp.Timestamp {
	if t != nil {
		var ts *timestamp.Timestamp
		ts, err := ptypes.TimestampProto(*t)
		if err == nil {
			return ts
		}
	}

	return nil
}

// WrapStr function returns nil when value is empty otherwise returns value in wrappers.StringValue format.
func WrapStr(value string) *wrappers.StringValue {
	if value != "" {
		return &wrappers.StringValue{Value: value}
	}

	return nil
}
