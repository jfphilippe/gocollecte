//  Copyright jean-françois PHILIPPE 2014-2018

package config

import (
	"strings"
	"testing"
)

func TestNew0(t *testing.T) {
	conf := New(nil)
	if conf == nil {
		t.Error("New returned nil")
	}
}

func TestNew1(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	if conf == nil {
		t.Error("New returned nil")
	} else {
		v, ok := conf.Raw("root")
		if !ok || v != "/tmp" {
			t.Error("'root' not found", v)
		}

		v, ok = conf.Raw("PROC")
		if !ok || v != "4" {
			t.Error("'PROC' not found", v)
		}

		v, ok = conf.Raw("nope")
		if ok || v != "" {
			t.Error("'nope' found", v, " ok :", ok)
		}

		v, ok = conf.Raw("section.nope")
		if ok || v != "" {
			t.Error("'nope' found '", v, "' ok :", ok)
		}

		v, ok = conf.Raw(".PROC")
		if !ok || v != "4" {
			t.Error("'.PROC' not found", v)
		}

	}
}

func TestLoad0(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := "some = true\nother = false"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, ok := conf.Raw("other")
		if !ok || v != "false" {
			t.Error("Raw('other') error : found ", v, " expected false")
		}

	}
}

// test de valeurs globales avec un commentaire
func TestLoad1(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, ok := conf.Raw("root")
		if !ok || v != "/tmp" {
			t.Error("'root' not found", v)
		}

		v, ok = conf.Raw("PROC")
		if !ok || v != "4" {
			t.Error("'PROC' not found", v)
		}

		v, ok = conf.Raw("nope")
		if ok || v != "" {
			t.Error("'nope' found", v, " ok :", ok)
		}

		v, ok = conf.Raw("other")
		if !ok || v != "false" {
			t.Error("Raw('other') error : found ", v, " expected false")
		}

		v, ok = conf.Raw("some")
		if !ok || v != "true" {
			t.Error("Raw('some') error : found ", v, " expected true")
		}

	}
}

// test de valeurs globales avec une section
func TestLoad2(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, ok := conf.Raw("root")
		if !ok || v == "" {
			t.Error("'root' not found", v)
		}

		v, ok = conf.Raw("PROC")
		if !ok || v == "" {
			t.Error("'PROC' not found", v)
		}

		v, ok = conf.Raw("other")
		if !ok || v == "" || v != "false" {
			t.Error("Raw('other') error : found ", v, " expected false")
		}

		v, ok = conf.Raw("some")
		if !ok || v == "" || v != "true" {
			t.Error("Raw('some') error : found ", v, " expected true")
		}

		v, ok = conf.Raw("section.item")
		if !ok || v == "" || v != "2" {
			t.Error("Raw('section.item') error : found ", v, " expected 2, ok :", ok, " for :", conf)
		}

	}
}

// test de valeurs globales avec deux sections
func TestLoad3(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2 \n[ section2 ]\n item = 3"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, ok := conf.Raw("some")
		if !ok || v == "" || v != "true" {
			t.Error("Raw('some') error : found ", v, " expected true")
		}

		v, ok = conf.Raw("section.item")
		if !ok || v == "" || v != "2" {
			t.Error("Raw('section.item') error : found ", v, " expected 2, ok :", ok, " for :", conf)
		}

		v, ok = conf.Raw("section2.item")
		if !ok || v == "" || v != "3" {
			t.Error("Raw('section2.item') error : found ", v, " expected 3, ok :", ok, " for :", conf)
		}

	}
}

// test d acces aux sections
func TestSection(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2 \n[ section.sub ]\n item = 3"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		// Lecture section existante
		section := conf.Section("section")
		if section == nil {
			t.Error("Section 'section' not found")
		}

		v, ok := section.Raw("item")
		if !ok || v == "" || v != "2" {
			t.Error("Section.Raw('item') error : found ", v, " expected 2, ok :", ok, " for :", conf)
		}

		v, ok = conf.Raw("section.sub.item")
		if !ok || v == "" || v != "3" {
			t.Error("Raw('section.sub.item') error : found ", v, " expected 3, ok :", ok, " for :", conf)
		}

		v, ok = section.Raw("sub.item")
		if !ok || v == "" || v != "3" {
			t.Error("Section.Raw('sub.item') error : found ", v, " expected 3, ok :", ok, " for :", conf)
		}

		section = conf.Section("none")
		if section == nil {
			t.Error("Section 'none' not found")
		}

	}
}

// test de valeurs globales avec deux sections
func TestString1(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2 \n[ section2 ]\n item = 3"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, err := conf.String("some")
		if err != nil || v != "true" {
			t.Error("Raw('some') error : found ", v, " expected true")
		}

		v, err = conf.String("section.item")
		if err != nil || v != "2" {
			t.Error("Raw('section.item') error : found ", v, " expected 2, for :", conf)
		}

		v, err = conf.String("section2.item")
		if err != nil || v != "3" {
			t.Error("Raw('section2.item') error : found ", v, " expected 3, for :", conf)
		}

		v, err = conf.String("section2.item2")
		if err == nil || v != "" {
			t.Error("Raw('section2.item2') error : found ", v, " expected '', for :", conf)
		}

	}
}

// test de valeurs globales avec deux sections
func TestString2(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2 \n[ section2 ]\n item = 3"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v, err := conf.String("some", "deflt")
		if err != nil || v != "true" {
			t.Error("Raw('some') error : found ", v, " expected true")
		}

		v, err = conf.String("section.item", "deflt")
		if err != nil || v != "2" {
			t.Error("Raw('section.item') error : found ", v, " expected 2, for :", conf)
		}

		v, err = conf.String("section2.item", "deflt")
		if err != nil || v != "3" {
			t.Error("Raw('section2.item') error : found ", v, " expected 3, for :", conf)
		}

		v, err = conf.String("section2.item2", "dflt")
		if err != nil || v != "dflt" {
			t.Error("Raw('section2.item2') error : found ", v, " expected 3, for :", conf)
		}

		v, err = conf.String("section3.item2", "dflt")
		if err != nil || v != "dflt" {
			t.Error("Raw('section3.item2') error : found ", v, " expected dflt, for :", conf)
		}

	}
}

// test de valeurs globales avec deux sections
func TestSections(t *testing.T) {
	m := map[string]string{"root": "/tmp",
		"PROC": "4",
	}
	conf := New(m)
	str := " # Comment \nsome = true\nother = false \n [ section ] \nitem=2 \n[ section2 ]\n item = 3"
	err := conf.Load(strings.NewReader(str))
	if err != nil {
		t.Error("Echec Load", err)
	} else {
		v := conf.Sections()
		if len(v) != 2 {
			t.Error("Sections() error : found ", v, " expected len == 2")
		}

		c := v["section"]
		if c == nil {
			t.Error("Sections() error : subsection 'section' not found ", v)
		} else {
			v, err := c.String("item")
			if err != nil || v != "2" {
				t.Error("Section('section2.item') error : found ", v, " expected 2, for :", c)
			}
		}

		c = v["section2"]
		if c == nil {
			t.Error("Sections() error : subsection 'section2' not found ", v)
		}

		c = v["section3"]
		if c != nil {
			t.Error("Sections() error : subsection 'section3' found ", v)
		}
	}
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
