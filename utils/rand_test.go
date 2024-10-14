package utils

import "testing"

func TestRandomPick(t *testing.T) {
	list := []int{1, 2, 3, 4, 5}
	picked := RandomPick(list, 3)
	if len(picked) != 3 {
		t.Errorf("expected 3 elements, got %d", len(picked))
	}
	t.Logf("picked: %v", picked)

	list = []int{1, 2, 3}
	picked = RandomPick(list, 3)
	if len(picked) != 3 {
		t.Errorf("expected 3 elements, got %d", len(picked))
	}
	t.Logf("picked: %v", picked)

	list = []int{1, 2}
	picked = RandomPick(list, 3)
	if len(picked) != 3 {
		t.Errorf("expected 3 elements, got %d", len(picked))
	}
	t.Logf("picked: %v", picked)
}
