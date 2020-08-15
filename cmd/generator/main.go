package main

import (
	"flag"
	"log"

	"github.com/Feresey/gen-tgpotapi/generator"
)

func main() {
	var (
		outDir       string
		templateDir  string
		apSchemaPath string
	)

	flag.StringVar(&outDir, "o", "", "output dir")
	flag.StringVar(&templateDir, "t", "", "template dir")
	flag.StringVar(&apSchemaPath, "s", "", "api schema path")
	flag.Parse()

	gen, err := generator.NewGenerator(apSchemaPath, templateDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := gen.Generate(outDir); err != nil {
		log.Fatal("Generate go files: ", err)
	}
	log.Print("Success")
}
