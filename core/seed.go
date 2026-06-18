package generate

import (
	"hash/fnv"
	"strconv"
	"time"
)

func ParseSeed(s string) int64 {
	if s == "" {
		return 0
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		h := fnv.New64a()
		h.Write([]byte(s))
		return int64(h.Sum64())
	}
	return n
}

func RandomSeed() int64 {
	return time.Now().UnixNano()
}
