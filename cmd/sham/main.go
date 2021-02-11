package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mattmeyers/sham"
)

var (
	oPrettyPrint = flag.Bool("pretty", false, "pretty print the output")
)

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	schema, err := readFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	if n := flag.NArg(); n > 1 || (schema != nil && n > 0) {
		log.Fatal("only a single schema can be processed")
	} else if n == 1 {
		schema = []byte(flag.Arg(0))
	}

	d, err := sham.Generate(schema)
	if err != nil {
		log.Fatal(err)
	}

	var out []byte
	if *oPrettyPrint {
		out, err = json.MarshalIndent(d, "", "    ")
	} else {
		out, err = json.Marshal(d)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
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
