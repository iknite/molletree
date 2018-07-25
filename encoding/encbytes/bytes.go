package encbytes

import "fmt"

func ToString(b []byte) string {
	return fmt.Sprintf("%X", b)
}
