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

	d, err := sham.Generate([]byte(flag.Arg(0)))
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
