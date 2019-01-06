// xsd2json - convert XSD files to JSON schema
// Copyright (C) 2019  Tom Hay

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.package main

// structs
// data structures used throughout the program

package main

// anything that has a name
type named interface {
	getName() string
}

// anything that has attributes
type attributed interface {
	getAttrs() []attribute
}

// any element
type element struct {
	name      string
	etype     string
	minOccurs int64
	maxOccurs int64
}

// any attribute
type attribute struct {
	name     string
	atype    string
	adefault string
	fixed    string
	required bool
}

// definition of a simple type
type simpleType struct {
	name string
	// restrictions
	base           string
	attrs          []attribute
	enum           []string
	minExclusive   int64
	minInclusive   int64
	maxExclusive   int64
	maxInclusive   int64
	totalDigits    int64
	fractionDigits int64
	length         int64
	minLength      int64
	maxLength      int64
	whiteSpace     string // preserve | replace | collapse
	pattern        string
}

// definition of a complex type
type complexType struct {
	name       string
	attrs      []attribute
	etype      string // sequence | choice
	elems      []element
	simpleBase *simpleType
	anyFlag    bool //does the type allow "any" extension?
}

// data being worked on
type context struct {
	inFile      string
	outFile     string
	inFileBase  string // base part of path
	outFileBase string
	domain      string
	smplType    *simpleType
	cplxType    *complexType
	elem        *element
	// the dictionary
	root         *element
	simpleTypes  map[string]simpleType
	complexTypes map[string]complexType
}

// initialise the context
func newContext() context {
	c := context{}
	c.simpleTypes = make(map[string]simpleType)
	c.complexTypes = make(map[string]complexType)
	return c
}

// implement interfaces
func (e element) getName() string {
	return e.name
}
func (s simpleType) getName() string {
	return s.name
}
func (c complexType) getName() string {
	return c.name
}
func (s simpleType) getAttrs() []attribute {
	return s.attrs
}
func (c complexType) getAttrs() []attribute {
	return c.attrs
}

// create a new element
func newElement() *element {
	return &element{
		name:      "",
		etype:     "",
		minOccurs: -1,
		maxOccurs: -1,
	}
}

// create a new simple type
func newSimpleType(aname string) *simpleType {
	return &simpleType{
		name:           aname,
		base:           "",
		attrs:          make([]attribute, 0),
		enum:           make([]string, 0),
		minExclusive:   -1,
		minInclusive:   -1,
		maxExclusive:   -1,
		maxInclusive:   -1,
		totalDigits:    -1,
		fractionDigits: -1,
		length:         -1,
		minLength:      -1,
		maxLength:      -1,
		whiteSpace:     "",
		pattern:        "",
	}
}

// create a new complex type
func newComplexType(aname string) *complexType {
	return &complexType{
		name:       aname,
		attrs:      make([]attribute, 0),
		etype:      "",
		elems:      make([]element, 0),
		simpleBase: nil,
	}
}

// create a clone of simpleType with a new name if specified
// Deep copy attrs & enum so they can be
// mutated without affecting the original.
func (s *simpleType) clone(name *string) *simpleType {
	n := s
	if name != nil {
		n.name = *name
	}
	n.attrs = append(make([]attribute, 0), s.attrs...)
	n.enum = append(make([]string, 0), s.enum...)
	return n
}

// create a clone of complexType with a new name if specified
// Deep copy attrs & elems so they can be
// mutated without affecting the original.
func (c *complexType) clone(name *string) *complexType {
	// n := newComplexType(c.name)
	// if name != nil {
	// 	n.name = *name
	// }
	// n.attrs = append(n.attrs, c.attrs...)
	// n.etype = c.etype
	// n.elems = append(n.elems, c.elems...)
	// n.simpleBase = c.simpleBase
	n := c
	if name != nil {
		n.name = *name
	}
	n.attrs = append(make([]attribute, 0), c.attrs...)
	n.elems = append(make([]element, 0), c.elems...)
	return n
}
