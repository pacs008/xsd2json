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

// xsd2json project main.go
// main function for xsd2json

package main

import (
	"fmt"
	"os"
)

func main() {

	// license notice
	fmt.Println("xsd2json Copyright (C) 2019  Tom Hay")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Println("under certain conditions; see COPYING.txt for details.")

	// initialise
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
