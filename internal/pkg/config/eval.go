//  Copyright jean-françois PHILIPPE 2014-2018

package config

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

const (
	// MaxRecursion Max recursive eval
	MaxRecursion = 10
)

// ErrMaxRecursion Error returned when max recursion is reached
var ErrMaxRecursion = errors.New("Max recursion reached :" + strconv.Itoa(MaxRecursion))

// Find matching } of ${ in a string.
// val is the remaining of the string. i.e : after ${
// return pos of } in string or -1
func (c *Config) matchEnd(val string) int {
	// For now first } found
	// Later may handle ${xx${yy}zz} (nested items)
	return strings.Index(val, "}")
}

func (c *Config) eval(buffer *bytes.Buffer, val string, deep uint) error {
	// Safe guard against infinite recursion
	if deep > MaxRecursion {
		return ErrMaxRecursion
	}
	remain := val // remaining of current string.

	var start, end int // Of ${xxx}

	// Search for ${
	start = strings.Index(remain, "${")
	for ; start >= 0; start = strings.Index(remain, "${") {
		buffer.WriteString(remain[:start])
		remain = remain[start+2:]
		end = c.matchEnd(remain)
		if end >= 0 {
			key := strings.TrimSpace(remain[:end])
			remain = remain[end+1:]
			subs, exists := c.Find(key)
			if exists {
				err := c.eval(buffer, subs, deep+1)
				if err != nil {
					return err
				}
			} else {
				return errors.New("Missing key '" + key + "'")
			}
		} else {
			buffer.WriteString("${")
		}
	}
	buffer.WriteString(remain)
	return nil
}

// Eval Evalue les interpolations eventuelles
// D une chaine.
func (c *Config) Eval(val string) (string, error) {
	start := strings.Index(val, "${")
	if start >= 0 {
		// May need substitution ...
		buffer := bytes.NewBufferString(val[:start])
		// Reserve some place.
		buffer.Grow(len(val))
		err := c.eval(buffer, val[start:], 0)
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	}
	// No subst pattern.
	return val, nil
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
