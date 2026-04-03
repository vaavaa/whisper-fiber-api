package tritonwhisper

import "testing"

func TestArgmaxLastTimeStep(t *testing.T) {
	// [1,2,4] — два шага, vocab 4
	logits := []float32{
		0, 0, 1, 0,
		0, 2, 0, 0,
	}
	shape := []int64{1, 2, 4}
	idx, err := ArgmaxLastTimeStep(logits, shape)
	if err != nil {
		t.Fatal(err)
	}
	if idx != 1 {
		t.Fatalf("ожидали индекс 1, получили %d", idx)
	}
}
