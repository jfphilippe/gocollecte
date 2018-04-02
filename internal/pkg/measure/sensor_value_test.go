/*
  Copyright jean-françois PHILIPPE 2014-2018
*/
package measure

import (
	"testing"
)

func TestNew(t *testing.T) {
	vals := NewNodeValues(50, 12)
	if len(vals.Values) > 0 {
		t.Error(
			"For len", vals,
			"expected", 0,
			"got", len(vals.Values),
		)
	}
	if vals.When != 50 {
		t.Error(
			"For When", vals,
			"expected", 50,
			"got", vals.When,
		)
	}
	if vals.Node != 12 {
		t.Error(
			"For Node_Id", vals,
			"expected", 12,
			"got", vals.Node,
		)
	}
}

func TestEmpty0(t *testing.T) {
	var vals NodeValues
	if len(vals.Values) > 0 {
		t.Error(
			"For len", vals,
			"expected", 0,
			"got", len(vals.Values),
		)
	}
}

func TestCopyEmpty(t *testing.T) {
	var val0 NodeValues
	vals := val0
	if len(vals.Values) > 0 {
		t.Error(
			"For len", vals,
			"expected", 0,
			"got", len(vals.Values),
		)
	}
}

func TestAppend0(t *testing.T) {
	var val0 NodeValues
	// append est declare sur *NodeValues
	vals := &val0
	vals.AppendValue(13, 0)
	if len(vals.Values) != 1 {
		t.Error(
			"For len", vals,
			"expected", 1,
			"got", len(vals.Values),
		)
	}
}

func TestCopy(t *testing.T) {
	var val0 NodeValues
	vals := &val0
	vals.AppendValue(13, 0)
	if len(vals.Values) != 1 {
		t.Error(
			"For len", vals,
			"expected", 1,
			"got", len(vals.Values),
		)
	}
	// Effectue une copie !
	val1 := val0
	if len(val1.Values) != 1 {
		t.Error(
			"For len", val1,
			"expected", 1,
			"got", len(val1.Values),
		)
	}
}

/* vi:set fileencodings=utf-8 tabstop=4 ai sw=2: */
