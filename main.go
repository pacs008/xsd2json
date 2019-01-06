// xsd2json project main.go
package main

import (
	"fmt"
	"os"
)

func main() {

	ctxt := newContext()
	cmdLineParse(&ctxt)

	// open the input file
	fname := ctxt.inFile
	inf, err := os.Open(fname)
	if err != nil {
		fmt.Printf("File %v open err %v", fname, err)
		os.Exit(2)
	}
	defer inf.Close()

	// open the output file
	fname = ctxt.outFile
	outf, err := os.Create(fname)
	if err != nil {
		fmt.Printf("File %v open err %v", fname, err)
		os.Exit(2)
	}
	defer outf.Close()

	parseXml(inf, &ctxt)
	writeJson(outf, &ctxt)

}
