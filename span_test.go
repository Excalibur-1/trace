package trace

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSpan(t *testing.T) {
	report := &mockReport{}
	t1 := NewTracer("service1", extendTag(), report, true)
	t.Run("test span string", func(t *testing.T) {
		sp1 := t1.New("testFinish").(*Span)
		assert.NotEmpty(t, fmt.Sprint(sp1))
	})
	t.Run("test fork", func(t *testing.T) {
		sp1 := t1.New("testFork").(*Span)
		sp2 := sp1.Fork("xxx", "opt_2").(*Span)
		assert.Equal(t, sp1.context.TraceId, sp2.context.TraceId)
		assert.Equal(t, sp1.context.SpanId, sp2.context.ParentId)
		t.Run("test max fork", func(t *testing.T) {
			sp3 := sp2.Fork("xx", "xxx")
			for i := 0; i < 100; i++ {
				sp3 = sp3.Fork("", "xxx")
			}
			assert.Equal(t, noopSpan{}, sp3)
		})
		t.Run("test max children", func(t *testing.T) {
			sp3 := sp2.Fork("xx", "xxx")
			for i := 0; i < 4096; i++ {
				sp3.Fork("", "xxx")
			}
			assert.Equal(t, noopSpan{}, sp3.Fork("xx", "xx"))
		})
	})
	t.Run("test finish", func(t *testing.T) {
		t.Run("test finish ok", func(t *testing.T) {
			sp1 := t1.New("testFinish").(*Span)
			time.Sleep(time.Millisecond)
			sp1.Finish(nil)
			assert.True(t, sp1.startTime.Unix() > 0)
			assert.True(t, sp1.duration > time.Microsecond)
		})
		t.Run("test finish error", func(t *testing.T) {
			sp1 := t1.New("testFinish").(*Span)
			time.Sleep(time.Millisecond)
			err := fmt.Errorf("🍻")
			sp1.Finish(&err)
			assert.True(t, sp1.startTime.Unix() > 0)
			assert.True(t, sp1.duration > time.Microsecond)
			errorTag := false
			for _, tag := range sp1.tags {
				if tag.Key == TagError && tag.Value != nil {
					errorTag = true
				}
			}
			assert.True(t, errorTag)
			messageLog := false
			for _, log := range sp1.logs {
				assert.True(t, log.Timestamp != 0)
				for _, field := range log.Fields {
					if field.Key == LogMessage && len(field.Value) != 0 {
						messageLog = true
					}
				}
			}
			assert.True(t, messageLog)
		})
		t.Run("test finish error stack", func(t *testing.T) {
			sp1 := t1.New("testFinish").(*Span)
			time.Sleep(time.Millisecond)
			err := fmt.Errorf("🍻")
			err = errors.WithStack(err)
			sp1.Finish(&err)
			ok := false
			for _, log := range sp1.logs {
				for _, field := range log.Fields {
					if field.Key == LogStack && len(field.Value) != 0 {
						ok = true
					}
				}
			}
			assert.True(t, ok, "LogStack set")
		})
		t.Run("test too many tags", func(t *testing.T) {
			sp1 := t1.New("testFinish").(*Span)
			for i := 0; i < 1024; i++ {
				sp1.SetTag(Tag{Key: strconv.Itoa(i), Value: "hello"})
			}
			assert.Len(t, sp1.tags, _maxTags+1)
			assert.Equal(t, sp1.tags[_maxTags].Key, "trace.error")
			assert.Equal(t, sp1.tags[_maxTags].Value, "too many tags")
		})
		t.Run("test too many logs", func(t *testing.T) {
			sp1 := t1.New("testFinish").(*Span)
			for i := 0; i < 1024; i++ {
				sp1.SetLog(LogField{Key: strconv.Itoa(i), Value: "hello"})
			}
			assert.Len(t, sp1.logs, _maxLogs+1)
			assert.Equal(t, sp1.logs[_maxLogs].Fields[0].Key, "trace.error")
			assert.Equal(t, sp1.logs[_maxLogs].Fields[0].Value, []byte("too many logs"))
		})
	})
}
