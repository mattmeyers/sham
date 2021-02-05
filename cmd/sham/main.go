package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mattmeyers/sham"
)

var (
	oPrettyPrint = flag.Bool("pretty", false, "pretty print the output")
)

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	ast, err := sham.NewParser(flag.Arg(0)).Parse()
	if err != nil {
		log.Fatal(err)
	}

	var out []byte
	if *oPrettyPrint {
		out, err = json.MarshalIndent(ast.Generate(), "", "    ")
	} else {
		out, err = json.Marshal(ast.Generate())
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
