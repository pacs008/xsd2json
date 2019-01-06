// xtype2j
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
