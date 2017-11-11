//  Copyright jean-françois PHILIPPE 2014-2016

package config

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// ParseError Parsing Error, memorise Line Number
type ParseError struct {
	LineNb int
	Msg    string
}

// Error return Error String
func (e *ParseError) Error() string {
	return "ParseError line:" + strconv.Itoa(e.LineNb) + " " + e.Msg
}

// Load Parse a conf from a Reader.
func (c *Config) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	// Boucle sur les lignes
	var lineNb int
	vals := c.values
	for scanner.Scan() {
		lineNb++
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				// Sections ?
				line = line[1 : len(line)-1]
				line = strings.TrimSpace(line)
				if line != "" {
					// Cree une nouvelle section
					vals = *c.sectionA(strings.Split(line, "."), true)
				} else {
					return &ParseError{lineNb, "Empty Section Name"}
				}
			} else {
				words := strings.SplitN(line, "=", 2)
				if len(words) != 2 {
					return &ParseError{lineNb, "Could not parse line"}
				}
				key := strings.TrimSpace(words[0])
				value := strings.TrimSpace(words[1])
				vals[key] = value
			}
		}
	}
	return nil
}

// LoadFile Load a file in current Config
func (c *Config) LoadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return c.Load(r)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
