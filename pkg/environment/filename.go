package environment

import (
	"fmt"
	"time"
)

type Filename string

func (f Filename) WithUnixSuffix() string {
	return fmt.Sprintf("%s_%d", f, unixNowMilli())
}

func (f Filename) String() string {
	return string(f)
}

func unixNowMilli() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}
