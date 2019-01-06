// cmdLineParse
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func cmdLineParse(ctxt *context) {
	inFilePtr := flag.String("in", "", "input file name")
	outFilePtr := flag.String("out", "", "output file name")
	domainPtr := flag.String("dom", "", "domain name for $id")

	flag.Parse()

	if *inFilePtr == "" || *outFilePtr == "" {
		fmt.Printf("Usage: %s -in xsdfile -out jsonfile [-dom domain]", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	ctxt.inFile = *inFilePtr
	ctxt.inFileBase = filepath.Base(ctxt.inFile)
	ctxt.outFile = *outFilePtr
	ctxt.outFileBase = filepath.Base(ctxt.outFile)
	ctxt.domain = *domainPtr
}
