package analysis

import (
	"github.com/PaluMacil/decisive-oak/parse"
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

func TestBuildTree(t *testing.T) {
	sample, err := parse.FromFile("../new-treatment.data.txt")
	if err != nil {
		t.Errorf("failed parsing file new-treatment.data.txt: %v", err)
	}
	tree := BuildTree(sample)
	nodeCount := tree.Root().CountNodes()
	if nodeCount != 3 {
		t.Errorf("node count: expected 3, got %d", nodeCount)
	}
}

func TestMain(m *testing.M) {
	n2 = &n1.Children[0]
	n2.parent = n1

	n3 = &n1.Children[1]
	n3.parent = n1

	n4 = &n2.Children[0]
	n4.parent = n2

	n5 = &n3.Children[0]
	n5.parent = n3

	n6 = &n3.Children[1]
	n6.parent = n3

	os.Exit(m.Run())
}

var n1 = &Node{
	Sample:      Sample{},
	FilterValue: "",
	Label:       "n1",
	Terminal:    false,
	parent:      nil,
	Children: []Node{
		{
			Sample:      Sample{},
			FilterValue: "",
			Label:       "n2",
			Terminal:    false,
			Children: []Node{
				{
					Sample:      Sample{},
					FilterValue: "",
					Label:       "n4",
					Terminal:    true,
				},
			},
		},
		{
			Sample:      Sample{},
			FilterValue: "",
			Label:       "n3",
			Terminal:    false,
			Children: []Node{
				{
					Sample:      Sample{},
					FilterValue: "",
					Label:       "n5",
					Terminal:    true,
				},
				{
					Sample:      Sample{},
					FilterValue: "",
					Label:       "n6",
					Terminal:    true,
				},
			},
		},
	},
}
var n2 *Node
var n3 *Node
var n4 *Node
var n5 *Node
var n6 *Node
