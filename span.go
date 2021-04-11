package trace

import (
	"fmt"
	"time"

	"github.com/Excalibur-1/trace/proto"
)

var _ Trace = &Span{}

const (
	_maxChildren = 1024
	_maxTags     = 128
	_maxLogs     = 256
)

// Span is a trace span.
type Span struct {
	dapper        *dapper
	context       spanContext
	operationName string
	startTime     time.Time
	duration      time.Duration
	tags          []Tag
	logs          []*proto.Log
	children      int
}

func (s *Span) TraceId() string {
	return s.context.String()
}

func (s *Span) Fork(serviceName, operationName string) Trace {
	if s.children > _maxChildren {
		// if child span more than max children set return noopSpan
		return noopSpan{}
	}
	s.children++
	// 为了兼容临时为 New 的 Span 设置 span.kind
	return s.dapper.newSpanWithContext(operationName, s.context).SetTag(tagString(TagSpanKind, "client"))
}

func (s *Span) Follow(serviceName, operationName string) Trace {
	return s.Fork(serviceName, operationName).SetTag(tagString(TagSpanKind, "producer"))
}

func (s *Span) Finish(perr *error) {
	s.duration = time.Since(s.startTime)
	if perr != nil && *perr != nil {
		err := *perr
		s.SetTag(tagBool(TagError, true))
		s.SetLog(Log(LogMessage, err.Error()))
		if err, ok := err.(stackTracer); ok {
			s.SetLog(Log(LogStack, fmt.Sprintf("%+v", err.StackTrace())))
		}
	}
	s.dapper.report(s)
}

func (s *Span) SetTag(tags ...Tag) Trace {
	if !s.context.isSampled() && !s.context.isDebug() {
		return s
	}
	if len(s.tags) < _maxTags {
		s.tags = append(s.tags, tags...)
	}
	if len(s.tags) == _maxTags {
		s.tags = append(s.tags, Tag{Key: "trace.error", Value: "too many tags"})
	}
	return s
}

// SetLog LogFields是一种有效且经过类型检查的方式来记录key:value
// 注意:当前不支持
func (s *Span) SetLog(logs ...LogField) Trace {
	if !s.context.isSampled() && !s.context.isDebug() {
		return s
	}
	if len(s.logs) < _maxLogs {
		s.setLog(logs...)
	}
	if len(s.logs) == _maxLogs {
		s.setLog(LogField{Key: "trace.error", Value: "too many logs"})
	}
	return s
}

// Visit visits the k-v pair in trace, calling fn for each.
func (s *Span) Visit(fn func(k, v string)) {
	fn(FosTraceID, s.context.String())
}

// SetTitle reset trace title
func (s *Span) SetTitle(operationName string) {
	s.operationName = operationName
}

func (s *Span) String() string {
	return s.context.String()
}

func (s *Span) ServiceName() string {
	return s.dapper.serviceName
}

func (s *Span) OperationName() string {
	return s.operationName
}

func (s *Span) StartTime() time.Time {
	return s.startTime
}

func (s *Span) Duration() time.Duration {
	return s.duration
}

func (s *Span) Context() spanContext {
	return s.context
}

func (s *Span) Tags() []Tag {
	return s.tags
}

func (s *Span) Logs() []*proto.Log {
	return s.logs
}

func (s *Span) setLog(logs ...LogField) Trace {
	protoLog := &proto.Log{
		Timestamp: time.Now().UnixNano(),
		Fields:    make([]*proto.Field, len(logs)),
	}
	for i := range logs {
		protoLog.Fields[i] = &proto.Field{Key: logs[i].Key, Value: []byte(logs[i].Value)}
	}
	s.logs = append(s.logs, protoLog)
	return s
}
