package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemBroadcast(t *testing.T) {
	ex := NewGlobalSynchronous(context.Background())
	sys := NewSystem(nil, ex)
	fooCount := 0
	foo := DeriverFunc(func(ev Event) bool {
		switch ev.(type) {
		case TestEvent:
			fooCount += 1
		case FooEvent:
			fooCount += 1
		default:
			return false
		}
		return true
	})
	barCount := 0
	bar := DeriverFunc(func(ev Event) bool {
		switch ev.(type) {
		case TestEvent:
			barCount += 1
		case BarEvent:
			barCount += 1
		default:
			return false
		}
		return true
	})
	fooEm := sys.Register("foo", foo, DefaultRegisterOpts())
	fooEm.Emit(TestEvent{})
	barEm := sys.Register("bar", bar, DefaultRegisterOpts())
	barEm.Emit(TestEvent{})
	// events are broadcast to every deriver, regardless who sends them
	require.NoError(t, ex.Drain())
	require.Equal(t, 2, fooCount)
	require.Equal(t, 2, barCount)
	// emit from bar, process in foo
	barEm.Emit(FooEvent{})
	require.NoError(t, ex.Drain())
	require.Equal(t, 3, fooCount)
	require.Equal(t, 2, barCount)
	// emit from foo, process in bar
	fooEm.Emit(BarEvent{})
	require.NoError(t, ex.Drain())
	require.Equal(t, 3, fooCount)
	require.Equal(t, 3, barCount)
}
