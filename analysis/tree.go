package analysis

import (
	"fmt"
	"github.com/PaluMacil/decisive-oak/parse"
)

func BuildTree(sample parse.Sample) (Node, error) {
	newSample, err := NewSample(sample)
	if err != nil {
		return Node{}, fmt.Errorf("building analyzed sample from parse sample")
	}
	rootNode, err := build(newSample, "", nil)
	if err != nil {
		return rootNode, fmt.Errorf("building root node: %w", err)
	}
	return rootNode, nil
}

func build(sample Sample, filterValue string, parent *Node) (Node, error) {
	if parent == nil {
		fmt.Println("starting first node")
	} else {
		fmt.Println("starting node", filterValue, parent.Sample.BestGainAttribute.Name)
	}
	var s Sample
	if filterValue != "" {
		if parent == nil {
			return Node{}, fmt.Errorf("parent cannot be nil when a filter value is given")
		}
		// data passed in before filtering
		unfilteredSampleData := sample.data
		// data filtered by parent's best gain attribute and the filter value for this child node
		filteredData, err := unfilteredSampleData.Filter(parent.Sample.BestGainAttribute.Name, filterValue)
		if err != nil {
			return Node{}, fmt.Errorf("filtering data from passed in sample: %w", err)
		}
		s, err = NewSample(filteredData)
		if err != nil {
			return Node{}, fmt.Errorf("creating a new analysis sample from filtered data: %w", err)
		}
	} else {
		// if not filtering the data
		s = sample
	}

	/*
		Terminal node definitions from https://en.wikipedia.org/wiki/ID3_algorithm
	*/

	// 1) Every element in the subset belongs to the same class;
	// in which case the node is turned into a leaf node and labelled with the class of the examples.
	if len(s.Targets) == 1 && len(s.data.Examples) > 0 {
		node := Node{
			parent:      parent,
			Children:    nil,
			Sample:      s,
			FilterValue: filterValue,
			Label:       s.Targets[0],
			Terminal:    true,
		}
		fmt.Println("completed node", node.FilterValue, node.Label)
		return node, nil
	}
	// 2) There are no more attributes to be selected, but the examples still do not belong to the same
	// class. In this case, the node is made a leaf node and labelled with the most common class of
	// the examples in the subset.
	if len(s.AttributeTypes) == 0 && len(s.data.Examples) > 0 {
		node := Node{
			parent:      parent,
			Children:    nil,
			Sample:      s,
			FilterValue: filterValue,
			Terminal:    true,
		}
		node.Label = node.mostCommonTarget()

		fmt.Println("completed node", node.FilterValue, node.Label)
		return node, nil
	}

	// 3) There are no examples in the subset, which happens when no example in the parent set was found
	// to match a specific value of the selected attribute. An example could be the absence of a person
	// among the population with age over 100 years. Then a leaf node is created and labelled with the
	// most common class of the examples in the parent node's set.
	if len(s.data.Examples) == 0 {
		if parent == nil {
			return Node{}, fmt.Errorf("parent cannot be nil when there are no remaining examples")
		}
		node := Node{
			parent:      parent,
			Children:    nil,
			Sample:      s,
			FilterValue: "",
			Label:       parent.mostCommonTarget(),
			Terminal:    true,
		}
		fmt.Println("completed node", node.FilterValue, node.Label)
		return node, nil
	}

	bestGainAttribute := s.BestGainAttribute
	fmt.Printf("\t%s has %d values\n", bestGainAttribute.Name, len(bestGainAttribute.Values))
	node := Node{
		parent:      parent,
		Sample:      s,
		FilterValue: filterValue,
		Label:       bestGainAttribute.Name,
		Terminal:    false,
	}

	var children []Node
	for _, value := range bestGainAttribute.Values {
		fmt.Println("\texamining value", value.Value, "of", bestGainAttribute.Name)
		fmt.Println("\tsecond attribute type name matches:", s.AttributeTypes[1].Name == s.data.AttributeTypes[1].Name)
		child, err := build(s, value.Value, &node)
		if err != nil {
			var label string
			if parent == nil {
				label = "root"
			} else {
				label = parent.Label
			}
			return Node{}, fmt.Errorf("building child node of %s: %w", label, err)
		}
		children = append(children, child)
	}
	node.Children = children

	fmt.Println("completed node", node.FilterValue, node.Label)
	return node, nil
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

func NewSample(sample parse.Sample) (Sample, error) {
	targetTotals := make([]int, sample.NumTargets)
	for _, eg := range sample.Examples {
		i, err := sample.Targets.Index(eg.Target)
		if err != nil {
			return Sample{}, fmt.Errorf("finding targets for calculating totals: %w", err)
		}
		targetTotals[i] = targetTotals[i] + 1
	}
	entropySet := entropy(targetTotals)
	attributeTypes, err := getAttributeTypes(sample, entropySet)
	if err != nil {
		return Sample{}, fmt.Errorf("getting analysis attribute types of new sample: %w", err)
	}
	newSample := Sample{
		Targets:           []string(sample.Targets),
		Entropy:           entropySet,
		AttributeTypes:    attributeTypes,
		BestGainAttribute: getBestGainAttribute(attributeTypes),
		data:              sample,
	}

	return newSample, nil
}

func getAttributeTypes(sample parse.Sample, entropySet float64) (AttributeTypes, error) {
	attributeTypes := make(AttributeTypes, sample.NumAttributes)
	for iAV, at := range sample.AttributeTypes {
		attributeOccurrenceLookup, err := at.OccurrencesInTargets(sample)
		if err != nil {
			return attributeTypes, fmt.Errorf("getting attribute type occurrences in targets for %s: %w", at.Name, err)
		}
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

	return attributeTypes, nil
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
