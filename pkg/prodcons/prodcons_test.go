package prodcons_test

import (
	"cmp"
	"go-ipc/pkg/prodcons"
	"slices"
	"testing"
)

func TestProdConsRunner(t *testing.T) {
	var srcData []any = []any{1, 3, 4, 3}

	slices.SortFunc(srcData, func(a, b any) int {
		// You must assert the types to compare them
		return cmp.Compare(a.(int), b.(int))
	})

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		producer prodcons.ProducerI
		consumer prodcons.ConsumerI
		want     prodcons.ProdCondsRes
	}{
		struct {
			name     string
			producer prodcons.ProducerI
			consumer prodcons.ConsumerI
			want     prodcons.ProdCondsRes
		}{
			name:     "test1",
			producer: prodcons.NewProducer(srcData),
			consumer: prodcons.NewConsumer(),
			want: prodcons.ProdCondsRes{
				IsDone:   true,
				Consumed: srcData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := prodcons.NewProdCons(tt.producer, tt.consumer)
			got := pc.Runner()

			isConsumedEqual := slices.EqualFunc(got.Consumed, tt.want.Consumed, func(a, b any) bool {
				return a.(int) == b.(int)
			})

			if !isConsumedEqual || !got.IsDone {
				t.Errorf("Runner() = %v, want %v", got, tt.want)
			}
		})
	}
}
