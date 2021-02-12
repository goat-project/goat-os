package util

import (
	"net"
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

// IsPublicIPv4 function returns true when IP is public IPv4 otherwise returns false.
func IsPublicIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}

	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}

	return false
}

// Contains check if a slice contains an element.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
