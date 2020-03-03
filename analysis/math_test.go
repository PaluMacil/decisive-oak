package analysis

import (
	"fmt"
	"testing"
)

func TestEntropy(t *testing.T) {
	if fmt.Sprintf("%.2f", entropy([]int{3, 2})) != "0.97" {
		t.Errorf("incorrect entropy for 3, 2")
	}
	if fmt.Sprintf("%.2f", entropy([]int{1, 0})) != "0.00" {
		t.Errorf("incorrect entropy for 1, 0")
	}
	if fmt.Sprintf("%.2f", entropy([]int{1, 2})) != "0.92" {
		t.Errorf("incorrect entropy for 1, 2")
	}
	if fmt.Sprintf("%.2f", entropy([]int{1, 1})) != "1.00" {
		t.Errorf("incorrect entropy for 1, 1")
	}
}
