package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/alexey-medvedchikov/parser-from-scratch/internal/ast"
	"github.com/alexey-medvedchikov/parser-from-scratch/internal/parser"
	"github.com/alexey-medvedchikov/parser-from-scratch/internal/tokenizer"
)

func main() {
	var progCode string

	flag.StringVar(&progCode, "c", "", "Expression to parse")
	flag.Parse()

	if progCode == "" {
		args := flag.Args()
		if len(args) == 0 {
			flag.Usage()
			return
		}

		var err error
		progCode, err = readFiles(args)
		if err != nil {
			log.Fatalln(err)
		}
	}

	astTree, err := parse(progCode)
	if err != nil {
		log.Fatalln(err)
	}

	if err := dumpJSON(os.Stdout, astTree); err != nil {
		log.Fatalln(err)
	}
}

func readFiles(paths []string) (string, error) {
	var buf bytes.Buffer

	for _, fpath := range paths {
		if err := readFile(fpath, &buf); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func readFile(fpath string, w io.Writer) error {
	fp, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := fp.Close(); closeErr != nil {
			log.Printf("could not close file: %s", closeErr)
		}
	}()

	_, err = io.Copy(w, fp)
	return err
}

func parse(s string) (ast.Node, error) {
	var b ast.Builder

	tok := tokenizer.NewTokenizer(tokenizer.DefaultRules, s)
	p := parser.NewParser(tok, b)

	return p.Parse()
}

func dumpJSON(w io.Writer, tree ast.Node) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	return encoder.Encode(tree)
}
