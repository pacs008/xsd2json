## Background
In the world of bank-to-bank payments, the standard for message formats is ISO20022. This is an XML format, and there are many message types defined by XSDs at https://www.iso20022.org/. At the same time, there is increasing usage of APIs for payments. Hence there is a need to represent ISO20022 messages as XML. To ensure that the mapping is done correctly, a tool to convert XSD to JSON Schema was needed. **xsd2json** is that tool.
## Usage
**xsd2json -in XSDfilename -out JSONschemafilename [-dom domainname]**
Reads the input XSD, parses it into internal data strcutures, then writes it out as JSON schema. The optional domain parameter is used to generate the "$id" key for the file.
## Features
xsd2json supports the key XSD features, including:
- Mapping of XSD inbuilt types to JSON types
- Use of "$ref" to simplify the JSON schema
- Enforcing field presence via "required": [...]
- Enforcing strict compliance via "additionalProperties": false
- Restrictions on strings (length, pattern, enum)
- Restrictions on numbers (min, max)
- Support for XSD choices via "oneOf"
## Attributes
There is no direct support for attributes in JSON Schema, so the following mapping convention is followed:
- Map to an object type
- The object contains key "#name": value of XML text
- The object also contains keys "@Attribname", one per attribute.
Example:
`
"IntrBkSttlmAmt": {
   "#name": 1234,
    "@Ccy": "GBP"
},
`
## Version support
xsd2json generates schema files compatible with JSON Schema draft 4.
## Known limitations
xsd2json has not been extensively tested. XSD is a rich and compex standard, and there are undoubtedly many XSDs that will break the current version.
