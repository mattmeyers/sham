package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mattmeyers/sham"
)

// format is a custom flag for defining the output format of the generated data.
// This type implements the flag.Value interface and acts as an enum. If an
// invalid value is provided on the command line, then and error will be returned.
// The currently supported output formats are (case insensitive):
//		- json (default)
//		- xml
type format string

func (f *format) Set(s string) error {
	s = strings.ToLower(s)

	if s != "json" && s != "xml" {
		return errors.New("unknown output format")
	}

	*f = format(s)

	return nil
}

func (f *format) Get() interface{} { return string(*f) }

func (f *format) String() string { return string(*f) }

var (
	oPrettyPrint bool
	oCount       int
	oOutFormat   format = format("json")
)

func initCLIApp() {
	flag.Usage = func() {
		fmt.Println(`sham is a tool for generating random data

Usage:

	sham [options] <schema>

Options:
	-f value	set the output format: json, xml (default json)
	-n int		the number of generations to perform (default 1)		
	-pretty		pretty print the result
	-h, --help	show this help message`)
	}

	flag.BoolVar(&oPrettyPrint, "pretty", false, "pretty print the output")
	flag.IntVar(&oCount, "n", 1, "the number of generations to perform")
	flag.Var(&oOutFormat, "f", "set the output format: json, xml")
	flag.Parse()
}

func main() {
	initCLIApp()
	rand.Seed(time.Now().Unix())

	schema, err := readFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	if n := flag.NArg(); n > 1 || (schema != nil && n > 0) {
		log.Fatal("only a single schema can be processed")
	} else if n == 1 {
		schema = []byte(flag.Arg(0))
	}

	p, err := sham.NewDefaultParser(schema).Parse()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < oCount; i++ {
		e, err := encoders[string(oOutFormat)](p.Generate())
		if err != nil {
			log.Fatal(err)
		}
		writeToStdout(e)
	}
}

func readFromStdin() ([]byte, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() == 0 {
		return nil, nil
	}

	return ioutil.ReadAll(os.Stdin)
}

type encoder func(interface{}) ([]byte, error)

var encoders = map[string]encoder{
	"json": encodeJSON,
	"xml":  encodeXML,
}

func encodeJSON(d interface{}) ([]byte, error) {
	var out []byte
	var err error

	if oPrettyPrint {
		out, err = json.MarshalIndent(d, "", "    ")
	} else {
		out, err = json.Marshal(d)
	}

	if err != nil {
		return nil, err
	}

	return out, nil
}

func encodeXML(d interface{}) ([]byte, error) {
	var out []byte
	var err error

	if oPrettyPrint {
		out, err = xml.MarshalIndent(d, "", "    ")
	} else {
		out, err = xml.Marshal(d)
	}

	if err != nil {
		return nil, err
	}

	return out, nil
}

func writeToStdout(d []byte) { fmt.Fprintln(os.Stdout, string(d)) }
