// parseXml
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

// parse the XML file handle and populate the context
func parseXml(f io.Reader, ctxt *context) {

	// start parsing
	decoder := xml.NewDecoder(f)
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil { // EOF
			break
		}
		// Inspect the type of the token just read.
		switch el := t.(type) {
		case xml.StartElement:
			startElement(&el, ctxt)
		case xml.EndElement:
			endElement(&el, ctxt)
		case xml.CharData:
			// fmt.Printf("charData: %v\n", el)
		case xml.Comment:
			// fmt.Printf("comment: %v\n", el)
		case xml.Directive:
			// fmt.Printf("directive: %v\n", el)
		case xml.ProcInst:
			// fmt.Printf("procInst: %v=>%v\n", el.Target, el.Inst)
		// case io.EOF:
		// 	fmt.Printf("End Of File")
		default:
			fmt.Printf("WTF? %v", el)
		}
	}

}

func startElement(el *xml.StartElement, ctxt *context) {
	// convert attrs into map (duplicate attrs will be lost)
	attrs := make(map[string]string)
	for _, attr := range el.Attr {
		attrs[attr.Name.Local] = attr.Value
	}
	switch el.Name.Local {
	case "element":
		// fmt.Printf("xml element\n")
		ctxt.elem = newElement()
		elem := ctxt.elem
		for name, value := range attrs {
			switch name {
			case "name":
				elem.name = value
			case "type":
				elem.etype = value
			case "minOccurs":
				elem.minOccurs, _ = strconv.ParseInt(value, 10, 64)
			case "maxOccurs":
				if value == "unbounded" {
					elem.maxOccurs = 9999999
				} else {
					elem.maxOccurs, _ = strconv.ParseInt(value, 10, 64)
				}
			default:
				fmt.Printf("\t%v: %v\n", name, value)
			}
		}

		// fmt.Printf("xml element %v: %v\n", elem.name, elem.etype)
		if ctxt.cplxType == nil {
			ctxt.root = elem
		} else {
			found := false
			// over-write if already exists
			for i, old := range ctxt.cplxType.elems {
				if old.name == elem.name {
					ctxt.cplxType.elems[i] = *elem
					found = true
				}
			}
			if !found {
				ctxt.cplxType.elems = append(ctxt.cplxType.elems, *elem)
			}
		}
	case "attribute":
		attr := attribute{}
		for name, value := range attrs {
			switch name {
			case "name":
				attr.name = value
			case "type":
				attr.atype = value
			case "default":
				attr.adefault = value
			case "fixed":
				attr.fixed = value
			case "use":
				attr.required = (value == "required")
			}
		}
		if ctxt.smplType != nil {
			ctxt.smplType.attrs = append(ctxt.smplType.attrs, attr)
		} else {
			ctxt.cplxType.attrs = append(ctxt.cplxType.attrs, attr)
		}
	case "sequence": // sequence and choice can also occur in extensions!
		fallthrough
	case "choice":
		ctxt.cplxType.etype = el.Name.Local
	case "restriction": // mandatory base attribute
		fallthrough
	case "extension":
		baseName := attrs["base"]
		if ctxt.smplType != nil { // we're doing a simple type
			ctxt.smplType.base = baseName
		} else {
			simpleBase, isSimple := ctxt.simpleTypes[baseName]
			complexBase, isComplex := ctxt.complexTypes[baseName]
			switch {
			case isSimple:
				// We are going to change this to a simple type
				smplName := ctxt.cplxType.name
				ctxt.cplxType = nil
				ctxt.smplType = simpleBase.clone(&smplName)
				//				ctxt.cplxType.simpleBase = &simpleBase
			case isComplex:
				// deep copy of the base type, then we can over-write / add to elements
				ctxt.cplxType = complexBase.clone(&ctxt.cplxType.name)
			default:
				fmt.Printf("Whoops! Complex %s no base type %s found\n", ctxt.cplxType.name, baseName)
			}
		}
	case "enumeration": // always nested within a simpleType
		ctxt.smplType.enum = append(ctxt.smplType.enum, el.Attr[0].Value)
	case "minInclusive":
		ctxt.smplType.minInclusive, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "maxInclusive":
		ctxt.smplType.maxInclusive, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "minExclusive":
		ctxt.smplType.minExclusive, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "maxExclusive":
		ctxt.smplType.maxExclusive, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "totalDigits":
		ctxt.smplType.totalDigits, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "fractionDigits":
		ctxt.smplType.fractionDigits, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "length":
		ctxt.smplType.length, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "minLength":
		ctxt.smplType.minLength, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "maxLength":
		ctxt.smplType.maxLength, _ = strconv.ParseInt(el.Attr[0].Value, 10, 64)
	case "whitespace":
		ctxt.smplType.whiteSpace = el.Attr[0].Value
	case "pattern":
		ctxt.smplType.pattern = el.Attr[0].Value
	case "simpleType":
		ctxt.smplType = newSimpleType(el.Attr[0].Value)
	case "complexType":
		ctxt.cplxType = newComplexType(el.Attr[0].Value)
	case "simpleContent": // holder for extension or restriction
		break
	case "any":
		ctxt.cplxType.anyFlag = true
	case "schema": // do nothing
	default:
		fmt.Printf("startElement: %v\n", el.Name.Local)
		for _, attr := range el.Attr {
			fmt.Printf("\t%v: %v\n", attr.Name.Local, attr.Value)
		}
	}
}

func endElement(el *xml.EndElement, ctxt *context) {
	switch el.Name.Local {
	//all the above do nothing
	case "enumeration":
	case "choice":
	case "sequence":
	case "restriction":
	case "minInclusive":
	case "maxInclusive":
	case "minExclusive":
	case "maxExclusive":
	case "totalDigits":
	case "fractionDigits":
	case "length":
	case "minLength":
	case "maxLength":
	case "whitespace":
	case "pattern":
	case "simpleContent":
	case "attribute":
	case "extension":
	case "any":
	case "schema":
		//all the above do nothing
	case "element":
		ctxt.elem = nil // force an error if assignment attempted
	case "simpleType":
		ctxt.simpleTypes[ctxt.smplType.name] = *ctxt.smplType
		// fmt.Printf("simpleType %+v", ctxt.smplType)
		ctxt.smplType = nil // force an error if assignment attempted
	case "complexType":
		if ctxt.smplType != nil {
			ctxt.simpleTypes[ctxt.smplType.name] = *ctxt.smplType
			ctxt.smplType = nil // force an error if assignment attempted
		} else {
			ctxt.complexTypes[ctxt.cplxType.name] = *ctxt.cplxType
			// fmt.Printf("complexType %+v", ctxt.cplxType)
			ctxt.cplxType = nil // force an error if assignment attempted
		}
	default:
		fmt.Printf("Unclassified endElement: %v\n", el.Name.Local)
	}
}
