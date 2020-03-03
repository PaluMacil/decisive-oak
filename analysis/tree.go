package analysis

import (
	"github.com/PaluMacil/decisive-oak/parse"
)

func BuildTree(sample parse.Sample) Node {
	newSample := NewSample(sample)
	rootNode := build(newSample, nil)
	return rootNode
}

func build(sample Sample, parent *Node) Node {
	if len(sample.Targets) == 1 {
		return Node{
			parent:   parent,
			Children: nil,
			Sample:   sample,
			Label:    sample.Targets[0].Name,
			Terminal: true,
		}
	}
	// TODO: other node exits and non-exits

	return Node{}
}

type Attribute struct {
	TypeName     string
	Values       []string
	Entropy      float64
	PofAttribute float64
	Gain         float64
}

type Attributes []Attribute

type Node struct {
	parent   *Node
	Children []Node
	Sample   Sample
	Label    string
	Terminal bool
}

type Sample struct {
	Targets           Targets
	Entropy           float64
	AttributeAnalysis Attributes
	BestGainAttribute Attribute
	data              parse.Sample
}

func NewSample(sample parse.Sample) Sample {
	targetCounts := make(map[string]int)
	for _, eg := range sample.Examples {
		targetCounts[eg.Target] = targetCounts[eg.Target] + 1
	}
	targets := make([]Target, sample.NumTargets)
	// TODO: needed for EACH AT
	//targetOccurrences := sample.AttributeTypes.
	for i, t := range sample.Targets {
		targets[i] = Target{
			Name:        t,
			Occurrences: targetCounts[t],
		}
	}
	newSample := Sample{
		Targets: targets,
		//TODO: Entropy: entropy(targetOccurrences),
		data: sample,
	}
	newSample.AttributeAnalysis = newSample.getAttributes()
	// TODO: newSample.BestGainAttribute = newSample.getBestGainAttribute()

	return newSample
}

func (s Sample) getAttributes() Attributes {
	// TODO: attributes and values need the type sp[lit into two
	var attributes Attributes

	return attributes
}

// TODO: after I can calculate gains!
//func (s Sample) getBestGainAttribute() Attribute {
//}

type Target struct {
	Name        string
	Occurrences int
}

type Targets []Target
