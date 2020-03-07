package analysis

import (
	"github.com/PaluMacil/decisive-oak/parse"
)

func BuildTree(sample parse.Sample) *Node {
	newSample := NewSample(sample)
	rootNode := build(newSample, "", nil)
	return rootNode
}

func build(sample Sample, filterValue string, parent *Node) *Node {
	s := NewSample(sample.data)
	if filterValue != "" {
		s = NewSample(sample.data.Filter(parent.Sample.BestGainAttribute.Name, filterValue))
	}

	/*
		Terminal node definitions from https://en.wikipedia.org/wiki/ID3_algorithm
	*/

	// 1) Every element in the subset belongs to the same class;
	// in which case the node is turned into a leaf node and labelled with the class of the examples.
	if len(s.Targets) == 1 && len(s.data.Examples) > 0 {
		return &Node{
			parent:      parent,
			Children:    nil,
			Sample:      s,
			FilterValue: "",
			Label:       s.Targets[0],
			Terminal:    true,
		}
	}
	// 2) There are no more attributes to be selected, but the examples still do not belong to the same
	// class. In this case, the node is made a leaf node and labelled with the most common class of
	// the examples in the subset.
	if len(s.AttributeTypes) == 0 && len(s.data.Examples) > 0 {
		node := &Node{
			parent:      parent,
			Children:    nil,
			Sample:      s,
			FilterValue: "",
			Terminal:    true,
		}
		node.Label = node.mostCommonTarget()

		return node
	}

	// 3) There are no examples in the subset, which happens when no example in the parent set was found
	// to match a specific value of the selected attribute. An example could be the absence of a person
	// among the population with age over 100 years. Then a leaf node is created and labelled with the
	// most common class of the examples in the parent node's set.
	if len(s.data.Examples) == 0 {
		return &Node{
			parent:      parent,
			Children:    nil,
			Sample:      s,
			FilterValue: "",
			Label:       parent.mostCommonTarget(),
			Terminal:    true,
		}
	}

	bestGainAttribute := s.BestGainAttribute
	node := &Node{
		parent:      parent,
		Sample:      s,
		FilterValue: filterValue,
		Label:       bestGainAttribute.Name,
		Terminal:    false,
	}
	var children []Node
	for _, value := range bestGainAttribute.Values {
		build(s, value.Value, node)
	}
	node.Children = children

	return node
}

type AttributeType struct {
	Name   string
	Gain   float64
	Values AttributeValues
}

type AttributeTypes []AttributeType

type AttributeValue struct {
	Value       string
	Entropy     float64
	Occurrences int
}

type AttributeValues []AttributeValue

type Node struct {
	parent      *Node
	Children    []Node
	Sample      Sample
	FilterValue string
	Label       string
	Terminal    bool
}

type Root Node

func (r Root) CountNodes() int {
	node := Node(r)
	total := new(int)
	return countNodes(&node, total)
}

func countNodes(currentNode *Node, runningTotal *int) int {
	*runningTotal += 1
	for _, n := range currentNode.Children {
		countNodes(&n, runningTotal)
	}

	return *runningTotal
}

func (n Node) Root() Root {
	node := n
	for {
		if node.parent == nil {
			return Root(node)
		}
		node = *node.parent
	}
}

func (n Node) mostCommonTarget() string {
	targetOccurrences := make(map[string]int)
	for _, eg := range n.Sample.data.Examples {
		t := eg.Target
		targetOccurrences[t] += 1
	}
	var highestName string
	var highestCount int
	for name, occurrences := range targetOccurrences {
		if occurrences > highestCount {
			highestName, highestCount = name, occurrences
		}
	}

	return highestName
}

type Sample struct {
	Targets           Targets
	Entropy           float64
	AttributeTypes    AttributeTypes
	BestGainAttribute AttributeType
	data              parse.Sample
}

func NewSample(sample parse.Sample) Sample {
	targetTotals := make([]int, sample.NumTargets)
	for _, eg := range sample.Examples {
		i := sample.Targets.Index(eg.Target)
		targetTotals[i] = targetTotals[i] + 1
	}
	entropySet := entropy(targetTotals)
	attributeTypes := getAttributeTypes(sample, entropySet)
	newSample := Sample{
		Targets:           []string(sample.Targets),
		Entropy:           entropySet,
		AttributeTypes:    attributeTypes,
		BestGainAttribute: getBestGainAttribute(attributeTypes),
		data:              sample,
	}

	return newSample
}

func getAttributeTypes(sample parse.Sample, entropySet float64) AttributeTypes {
	attributeTypes := make(AttributeTypes, sample.NumAttributes)
	for iAV, at := range sample.AttributeTypes {
		attributeOccurrenceLookup := at.OccurrencesInTargets(sample)
		attrValues := make(AttributeValues, at.NumValues)
		for iVal, v := range at.Values {
			targetOccurrences := attributeOccurrenceLookup[v]
			attrValues[iVal] = AttributeValue{
				Value:       v,
				Entropy:     entropy(targetOccurrences),
				Occurrences: attributeOccurrenceLookup.AttributeValueTotal(v),
			}
		}
		attributeTypes[iAV].Name = at.Name
		attributeTypes[iAV].Values = attrValues
		attributeTypes[iAV].Gain = gain(entropySet, attrValues...)
	}

	return attributeTypes
}

func getBestGainAttribute(attrTypes AttributeTypes) AttributeType {
	var attributeType AttributeType
	for _, at := range attrTypes {
		if at.Gain >= attributeType.Gain {
			attributeType = at
		}
	}

	return attributeType
}

type Targets []string
