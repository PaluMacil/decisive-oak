package analysis

import (
	"fmt"
	"os"
	"testing"
)

func TestNode_Root(t *testing.T) {
	lbl1, lbl3, lbl6 := n1.Root().Label, n3.Root().Label, n6.Root().Label
	if lbl1 != "n1" {
		t.Errorf("expected label of root to be n1, got %s", lbl1)
	}
	if lbl3 != "n1" {
		t.Errorf("expected label of root to be n1, got %s", lbl3)
	}
	if lbl6 != "n1" {
		t.Errorf("expected label of root to be n1, got %s", lbl6)
	}
}

func TestRoot_CountNodes(t *testing.T) {
	tree := n6
	count := tree.Root().CountNodes()
	if count != 6 {
		t.Errorf("expected 6, got %d", count)
	}
}

func TestMain(m *testing.M) {
	n1.Children = []Node{*n2, *n3}
	n2.Children = []Node{*n4}
	n3.Children = []Node{*n5, *n6}
	fmt.Println(1, len(n1.Children), len(n2.Children), len(n3.Children))
	os.Exit(m.Run())
}

var n1 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n1",
	Terminal:    false,
	parent:      nil,
}

var n2 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n2",
	Terminal:    false,
	parent:      n1,
}

var n3 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n3",
	Terminal:    false,
	parent:      n1,
}

var n4 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n4",
	Terminal:    true,
	parent:      n2,
}

var n5 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n5",
	Terminal:    true,
	parent:      n3,
}

var n6 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n6",
	Terminal:    true,
	parent:      n3,
}
