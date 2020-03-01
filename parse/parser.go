package parse

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func FromFile(filename string) (Sample, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Sample{}, fmt.Errorf("reading %s: %w", filename, err)
	}
	defer file.Close()
	return Parse(file)
}

func Parse(reader io.Reader) (Sample, error) {
	var sample Sample
	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty lines
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	if len(lines) < 6 {
		return sample, fmt.Errorf("a data file with less than 6 lines cannot be valid")
	}
	var err error
	sample.NumTargets, err = strconv.Atoi(lines[0])
	if err != nil {
		return sample, fmt.Errorf("parsing number of targets: %w", err)
	}
	sample.Targets = strings.Split(lines[1], ",")
	sample.NumAttributes, err = strconv.Atoi(lines[2])
	if err != nil {
		return sample, fmt.Errorf("parsing number of attributes: %w", err)
	}
	requiredLineCount := 3 + sample.NumAttributes + 1
	if len(lines) < requiredLineCount {
		return sample, fmt.Errorf("not enough lines left for attributes and example count")
	}
	indexOfNumExamples := sample.NumAttributes + 3
	attributeTypes, err := parseAttributeTypes(lines[3:indexOfNumExamples])
	if err != nil {
		return Sample{}, err
	}
	sample.AttributeTypes = attributeTypes
	sample.NumExamples, err = strconv.Atoi(lines[indexOfNumExamples])
	if err != nil {
		return sample, fmt.Errorf("parsing number of targets: %w", err)
	}
	requiredLineCount = requiredLineCount + sample.NumExamples
	if len(lines) < requiredLineCount {
		return sample, fmt.Errorf("not enough lines left for examples")
	}

	examples, err := parseExamples(
		attributeTypes,
		sample.Targets,
		lines[indexOfNumExamples+1:requiredLineCount],
	)
	if err != nil {
		return Sample{}, err
	}
	sample.Examples = examples

	return sample, nil
}

func parseAttributeTypes(lines []string) (AttributeTypes, error) {
	types := make([]AttributeType, len(lines))
	for i, line := range lines {
		splits := strings.Split(line, ",")
		if len(splits) < 2 {
			return types, fmt.Errorf("not enough data in line %d of attribute lines", i)
		}
		types[i].Name = splits[0]
		if splits[1] == "real" {
			types[i].Real = true
			continue
		}
		numValues, err := strconv.Atoi(splits[1])
		if err != nil {
			return types, fmt.Errorf("parsing number of attribute values for %s: %w", splits[0], err)
		}
		types[i].NumValues = numValues
		types[i].Values = splits[2:]
		foundValues := len(types[i].Values)
		if foundValues != numValues {
			return types, fmt.Errorf("incorrect number of values for %s: expected %d but found %d",
				splits[0], numValues, foundValues)
		}
	}

	return types, nil
}

func parseExamples(attributeTypes AttributeTypes, targets Targets, lines []string) (Examples, error) {
	examples := make([]Example, len(lines))
	for idxExample, line := range lines {
		splits := strings.Split(line, ",")
		if len(splits) != len(attributeTypes)+1 {
			return examples, fmt.Errorf("expected %d attributes and one target in splits, got %d total",
				len(attributeTypes), len(splits))
		}
		lastSplitIndex := len(splits) - 1
		target := splits[lastSplitIndex]
		if !targets.IsValid(target) {
			return examples, fmt.Errorf("got invalid target \"%s\", expected %s",
				target, attributeTypes)
		}
		examples[idxExample].Target = target
		for idxExampleAttribute, v := range splits[:lastSplitIndex] {
			if attributeTypes[idxExampleAttribute].Real {
				realValue, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return examples, fmt.Errorf("parsing the real value for attribute %s which is attribute numnber %d: %w",
						attributeTypes[idxExample].Name, idxExample, err)
				}
				examples[idxExample].RealValues = append(examples[idxExample].RealValues, realValue)
				continue
			}
			if !attributeTypes[idxExampleAttribute].IsValidValue(v) {
				return examples, fmt.Errorf("got invalid attribute value \"%s\", expected %v",
					v, attributeTypes[idxExample].Values)
			}
			examples[idxExample].StringValues = append(examples[idxExample].StringValues, v)
		}
	}

	return examples, nil
}
