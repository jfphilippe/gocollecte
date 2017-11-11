//  Copyright jean-françois PHILIPPE 2014-2016

package config

import (
	"strings"
	"testing"
)

// test de valeurs globales avec deux sections
func TestEval0(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2 \n[ section2 ]\n item = 3"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, err := conf.Eval("some")
		if err != nil || v != "some" {
			t.Error("Eval('some') error : found ", v, " expected some")
		}

		v, err = conf.Eval("${section.item}")
		if err != nil || v != "2" {
			t.Error("Eval('${section.item}') error : found ", v, " expected 2")
		}

		v, err = conf.Eval("-${some }-${ other}-")
		if err != nil || v != "-true-false-" {
			t.Error("Eval('s-${some }-${other}-') error : found ", v)
		}

		v, err = conf.Eval("-${some }-${ other-")
		if err != nil || v != "-true-${ other-" {
			t.Error("Eval('s-${some }-${ other-') error : found ", v)
		}

		v, err = conf.Eval("-${some -${ other}-")
		if err == nil {
			t.Error("Eval('-${some -${ other}-') error expected, found ", v)
		}

	}
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
