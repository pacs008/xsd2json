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

// writeJson
// Take the populated data structures and output JSON schema v4

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const tsz = 3 // tab size

// entry point for writing
func writeJson(f io.Writer, ctxt *context) {
	inPrintf(f, 0, "{\n")

	writeHdrs(f, ctxt, tsz)
	rootType := ctxt.complexTypes[ctxt.root.getName()]
	writeComplexBody(rootType, f, ctxt, tsz)

	writeDefinitions(f, ctxt, tsz)

	inPrintf(f, 0, "}\n")
}

// write the name of an element or type
func writeName(n named, f io.Writer, ctxt *context, indent int) {
	inPrintf(f, indent, "\"%s\": {\n", n.getName())
}

// write an element
// if multiple occurrences are allowed, make it an array of items
// of the specified type
func writeElement(el element, f io.Writer, ctxt *context, indent int) {
	if el.maxOccurs > 1 {
		inPrintf(f, indent, "\"%s\": {\n", el.getName())
		inPrintf(f, indent+tsz, "\"type\": \"array\",\n")
		inPrintf(f, indent+tsz, "\"items\": {\n")
		inPrintf(f, indent+tsz+tsz, "\"$ref\": \"#/definitions/%s\",\n", el.etype)
		inPrintf(f, indent+tsz, "},\n")
		inPrintf(f, indent, "},\n")
	} else {
		inPrintf(f, indent, "\"%s\": {\"$ref\": \"#/definitions/%s\" },\n", el.getName(), el.etype)
	}
}

// write the body of a simple type
// if it has attributes, turn it into an object
// the #value element represents the base type
// each attribute forms a separate element named @Attributename
func writeSimpleBody(simple simpleType, f io.Writer, ctxt *context, indent int) {
	if len(simple.attrs) > 0 {
		inPrintf(f, indent, "\"type\": \"object\",\n")
		inPrintf(f, indent, "\"properties\": {\n")
		inPrintf(f, indent+tsz, "\"#value\": {\n")
		writeSimpleProperties(simple, f, ctxt, indent+tsz+tsz)
		inPrintf(f, indent+tsz, "},\n")
		required := writeAttrs(simple, f, ctxt, indent+tsz)
		inPrintf(f, indent, "},\n")
		inPrintf(f, indent, "\"required\": %s\n,", arrayString(required))
		inPrintf(f, indent, "\"additionalProperties\": false,\n")
	} else {
		writeSimpleProperties(simple, f, ctxt, indent)
	}
}

// write the properties of a simple type
func writeSimpleProperties(simple simpleType, f io.Writer, ctxt *context, indent int) {
	jtype, mapped := mapTypename(simple.base)
	inPrintf(f, indent, "\"type\": \"%s\",\n", jtype)
	if mapped {
		inPrintf(f, indent, "\"$comment\": \"XML datatype was %s\",\n", simple.base)
	}
	// string constraints
	if simple.minLength > -1 {
		inPrintf(f, indent, "\"minLength\": %d,\n", simple.minLength)
	}
	if simple.maxLength > -1 {
		inPrintf(f, indent, "\"maxLength\": %d,\n", simple.maxLength)
	}
	if simple.length > -1 {
		inPrintf(f, indent, "\"minLength\": %d,\n", simple.length)
		inPrintf(f, indent, "\"maxLength\": %d,\n", simple.length)
	}
	if len(simple.enum) > 0 {
		inPrintf(f, indent, "\"enum\": %s,\n", arrayString(simple.enum))
	}
	if simple.pattern != "" {
		// double all slashes to make valid JSON escapes
		escaped := strings.Replace(simple.pattern, "\\", "\\\\", -1)
		inPrintf(f, indent, "\"pattern\": \"%s\",\n", escaped)
	}
	// number constraints
	if simple.minInclusive > -1 {
		inPrintf(f, indent, "\"minimum\": %d,\n", simple.minInclusive)
	}
	if simple.minExclusive > -1 {
		inPrintf(f, indent, "\"exclusiveMinimum\": %d,\n", simple.minExclusive)
	}
	if simple.maxInclusive > -1 {
		inPrintf(f, indent, "\"maximum\": %d,\n", simple.maxInclusive)
	}
	if simple.maxExclusive > -1 {
		inPrintf(f, indent, "\"exclusiveMaximum\": %d,\n", simple.maxExclusive)
	}
	// JSON schema can't handle these rules
	if simple.totalDigits > -1 {
		inPrintf(f, indent, "\"$comment\": \"XML specified totalDigits=%d\",\n", simple.totalDigits)
	}
	if simple.fractionDigits > -1 {
		inPrintf(f, indent, "\"$comment\": \"XML specified fractionDigits=%d\",\n", simple.fractionDigits)
	}
	if simple.whiteSpace != "" {
		inPrintf(f, indent, "\"$comment\": \"XML specified whiteSpace=%s\",\n", simple.whiteSpace)
	}
}

// close the braces
func writeClose(f io.Writer, ctxt *context, indent int) {
	inPrintf(f, indent, "},\n")
}

// write the file headers
func writeHdrs(f io.Writer, ctxt *context, indent int) {
	domain := "https://example.com"
	when := time.Now().Format(time.RFC1123)
	if ctxt.domain != "" {
		domain = ctxt.domain
	}
	hdrs := [...]string{
		"\"$id\": \"" + domain + "/" + ctxt.outFileBase + "\",\n",
		"\"$schema\": \"http://json-schema.org/draft-04/schema#\",\n",
		"\"title\": \"" + ctxt.outFileBase + "\",\n",
		"\"description\": \"Derived from " + ctxt.inFileBase + " by '" + filepath.Base(os.Args[0]) + "' on " + when + ".\",\n",
		"\n",
	}
	for _, str := range hdrs {
		inPrintf(f, indent, str)
	}
}

// write all the type definitions
func writeDefinitions(f io.Writer, ctxt *context, indent int) {
	inPrintf(f, indent, "\n")
	inPrintf(f, indent, "\"$comment\": \"---Type definitions---\",\n")
	inPrintf(f, indent, "\"definitions\": {\n")

	// print all the simple type definitions first
	inPrintf(f, indent, "\n")
	inPrintf(f, indent+tsz, "\"$comment\": \"---Simple type definitions---\",\n")
	for _, simple := range ctxt.simpleTypes {
		writeSimple(simple, f, ctxt, indent+tsz)
	}
	// now print all the complex type definitions
	inPrintf(f, indent, "\n")
	inPrintf(f, indent+tsz, "\"$comment\": \"---Complex type definitions---\",\n")
	for _, cmplx := range ctxt.complexTypes {
		writeComplex(cmplx, f, ctxt, indent+tsz)
	}
	inPrintf(f, indent, "}\n")
}

// write a simple type definition
func writeSimple(simple simpleType, f io.Writer, ctxt *context, indent int) {
	writeName(simple, f, ctxt, indent)
	writeSimpleBody(simple, f, ctxt, indent+tsz)
	writeClose(f, ctxt, indent)
}

// write a complex type definition
func writeComplex(cmplx complexType, f io.Writer, ctxt *context, indent int) {
	writeName(cmplx, f, ctxt, indent)
	writeComplexBody(cmplx, f, ctxt, indent+tsz)
	writeClose(f, ctxt, indent)
}

// write the body of a complex type
func writeComplexBody(cmplx complexType, f io.Writer, ctxt *context, indent int) {
	// if it's based on simple, do simple body
	if cmplx.simpleBase != nil {
		fmt.Printf("Doing simple body for %s: %v\n", cmplx.name, *cmplx.simpleBase)
		writeSimpleBody(*cmplx.simpleBase, f, ctxt, indent+tsz)
		return
	}
	inPrintf(f, indent, "\"type\": \"object\",\n")
	inPrintf(f, indent, "\"properties\": {\n")
	if len(cmplx.attrs) > 0 {
		fmt.Printf("Doing attrs for complex %s\n", cmplx.name)
		writeAttrs(cmplx, f, ctxt, indent+tsz)
	}
	required := make([]string, 0)
	for _, el := range cmplx.elems {
		writeElement(el, f, ctxt, indent+tsz)
		if el.minOccurs != 0 {
			required = append(required, el.getName())
		}
	}
	inPrintf(f, indent, "},\n")

	switch cmplx.etype {
	case "choice":
		// XSD choice maps to JSON schema thus:
		// "oneOf": [
		// {"required": ["Cd"] },
		// {"required": ["Prtry"] },
		// ],
		inPrintf(f, indent, "\"oneOf\": [\n")
		for _, el := range cmplx.elems {
			inPrintf(f, indent+tsz, "{\"required\": [\"%s\"]},\n", el.getName())
		}
		inPrintf(f, indent, "],\n")

	default:
		if len(required) > 0 {
			inPrintf(f, indent, "\"required\": %s,\n", arrayString(required))
		}
	}
	if !cmplx.anyFlag {
		inPrintf(f, indent, "\"additionalProperties\": false,\n")
	} else {
		inPrintf(f, indent, "\"$comment\": \"XSD allows 'any', so properties not restricted\",\n")
	}
}

func writeAttrs(attd attributed, f io.Writer, ctxt *context, indent int) []string {
	attrs := attd.getAttrs()
	required := []string{"#name"}
	for _, attr := range attrs {
		if attr.required {
			required = append(required, "@"+attr.name)
		}
		inPrintf(f, indent, "\"@%s\": {\n", attr.name)
		// atype must be either builtin or simple ...
		if _, ok := ctxt.simpleTypes[attr.atype]; ok {
			inPrintf(f, indent+tsz, "\"$ref\": \"#/definitions/%s\",\n", attr.atype)
		} else {
			inPrintf(f, indent+tsz, "\"type\": \"%s\",\n", attr.atype)
		}
		if attr.adefault != "" {
			inPrintf(f, indent+tsz, "\"default\": \"%s\",\n", attr.adefault)
		}
		if attr.fixed != "" {
			inPrintf(f, indent+tsz, "\"$comment\": \"XML specified fixed value %s\",\n", attr.fixed)
		}
		inPrintf(f, indent, "},\n")
	}
	return required
}

// print an arbitrary thing with an indent
func inPrintf(f io.Writer, indent int, s string, v ...interface{}) (int, error) {
	var n1, n2 int
	var err error

	n1, err = fmt.Fprintf(f, "%*s", indent, "")
	if err == nil {
		n2, err = fmt.Fprintf(f, s, v...)
	}
	if err != nil {
		fmt.Printf("Write failed: %v\n", err)
	}
	return n1 + n2, err
}

func arrayString(strs []string) string {
	s := "["
	for _, val := range strs {
		s += fmt.Sprintf("\"%s\",", val)
	}
	s = s[:len(s)-1]
	s += "]"
	return s
}
