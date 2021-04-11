package trace

import (
	"fmt"
	"testing"
)

func TestMockTrace(t *testing.T) {
	mockTrace := &MockTrace{}
	_ = mockTrace.Inject(nil, nil, nil)
	_, _ = mockTrace.Extract(nil, nil)
	root := mockTrace.New("test")
	root.Fork("", "")
	root.Follow("", "")
	root.Finish(nil)
	err := fmt.Errorf("test")
	root.Finish(&err)
	root.SetTag()
	root.SetLog()
	root.Visit(func(k, v string) {})
	root.SetTitle("")
}
