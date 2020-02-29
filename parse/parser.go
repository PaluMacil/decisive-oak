package parse

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func FromFile(filename, delimiter string) (Sample, error) {
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
	
}
