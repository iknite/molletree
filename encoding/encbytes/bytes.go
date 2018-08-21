package encbytes

import (
	"encoding/binary"
	"fmt"
)

func ToString(b []byte) string {
	return fmt.Sprintf("%X", b)
}

func ToPrettyId(id []byte) string {
	return fmt.Sprintf(
		"%d|%d",
		binary.BigEndian.Uint64(id[:8]),
		binary.BigEndian.Uint64(id[8:]),
	)
}
