package trace

import (
	"context"
	"encoding/binary"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var _hostHash byte

type ctxKey string

var _ctxKey ctxKey = "mos/trace.trace"

func init() {
	rand.Seed(time.Now().UnixNano())
	Hostname, err := os.Hostname()
	if err != nil {
		Hostname = strconv.Itoa(int(time.Now().UnixNano()))
	}
	_hostHash = byte(oneAtTimeHash(Hostname))
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// FromContext returns the trace bound to the context, if any.
func FromContext(ctx context.Context) (t Trace, ok bool) {
	t, ok = ctx.Value(_ctxKey).(Trace)
	return
}

// NewContext new a trace context.
// NOTE: This method is not thread safe.
func NewContext(ctx context.Context, t Trace) context.Context {
	return context.WithValue(ctx, _ctxKey, t)
}

func oneAtTimeHash(s string) (hash uint32) {
	b := []byte(s)
	for i := range b {
		hash += uint32(b[i])
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return
}

func genID() uint64 {
	var b [8]byte
	// i think this code will not survive to 2106-02-07
	binary.BigEndian.PutUint32(b[4:], uint32(time.Now().Unix())>>8)
	b[4] = _hostHash
	binary.BigEndian.PutUint32(b[:4], uint32(rand.Int31()))
	return binary.BigEndian.Uint64(b[:])
}
