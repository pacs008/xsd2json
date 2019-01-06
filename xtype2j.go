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

// xtype2j
// map XSD builtin types to JSON types

package main

import (
	"strings"
)

var xtype2j = map[string]string{
	// XML Schema Built-In Numeric Datatypes:
	"decimal":            "number",
	"float":              "number",
	"double":             "number",
	"integer":            "number",
	"positiveInteger":    "number",
	"negativeInteger":    "number",
	"nonPositiveInteger": "number",
	"nonNegativeInteger": "number",
	"long":               "number",
	"int":                "number",
	"short":              "number",
	"byte":               "number",
	"unsignedLong":       "number",
	"unsignedInt":        "number",
	"unsignedShort":      "number",
	"unsignedByte":       "number",
	// XML Schema Built-In Date, Time, and Duration Datatypes:
	"dateTime":          "string",
	"dateTimeStamp":     "string",
	"date":              "string",
	"time":              "string",
	"gYearMonth":        "string",
	"gYear":             "string",
	"duration":          "string",
	"dayTimeDuration":   "string",
	"yearMonthDuration": "string",
	"gMonthDay":         "string",
	"gDay":              "string",
	"gMonth":            "string",
	//XML Schema String Datatypes:
	// string hallalujah, no need to map this!
	"normalizedString": "string",
	"token":            "string",
	"language":         "string",
	"NMTOKEN":          "string",
	"NMTOKENS":         "string",
	"Name":             "string",
	"NCName":           "string",
	// XML Schema "Magic" Datatypes:
	"ID":       "string",
	"IDREF":    "string",
	"IDREFS":   "string",
	"ENTITY":   "string",
	"ENTITIES": "string",
	// XML Schema Oddball Datatypes:
	"QName": "string",
	//boolean hallalujah, no need to map this!
	"hexBinary":    "string",
	"base64Binary": "string",
	"anyURI":       "string",
	"notation":     "string",
}

// map XML typenames to JSON
func mapTypename(name string) (string, bool) {
	idx := strings.Index(name, ":")
	if idx > -1 {
		name = name[idx+1:]
	}
	jname, mapped := xtype2j[name]
	if mapped {
		name = jname
	}
	return name, mapped
}
