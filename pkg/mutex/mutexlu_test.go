package mutex_test

import (
	"go-ipc/pkg/mutex"
	"testing"
)

func TestMutexRunner(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want mutex.CounterRes
	}{
		// TODO: Add test cases.
		{
			name: "correct",
			want: mutex.CounterRes{
				FinalIncrement: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var m mutex.ServiceMutexI = mutex.NewMutex()
			got := m.Runner()
			// TODO: update the condition below to compare got with tt.want.
			if tt.want.FinalIncrement != 100 {
				t.Errorf("Runner() = %v, want %v", got, tt.want)
			}
		})
	}
}
