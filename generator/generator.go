package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/iancoleman/strcase"

	"golang.org/x/tools/imports"
)

const maxInLine = 100

func getEnumName(s string) string {
	switch s {
	case "MessageEntity":
		return "Entity"
	case "KeyboardButtonPollType":
		return "KeyboardButton"
	}
	return strings.Split(strcase.ToDelimited(s, ' '), " ")[0]
}

var (
	markupType = []TypeMapping{
		"InlineKeyboardMarkup",
		"ReplyKeyboardMarkup",
		"ReplyKeyboardRemove",
		"ForceReply",
	}
	inputType = []TypeMapping{
		"InputFile",
		"str",
	}
	intStr = []TypeMapping{
		"int",
		"str",
	}
)

func multitype(ss []TypeMapping) TypeMapping {
	if reflect.DeepEqual(ss, markupType) {
		return "ReplyMarkup"
	}
	if reflect.DeepEqual(ss, inputType) {
		return "InputDataType"
	}
	if reflect.DeepEqual(ss, intStr) {
		return "IntStr"
	}
	return ""
}

func getType(fieldName string, typeName string, types []TypeMapping) TypeMapping {
	if len(types) > 1 {
		return multitype(types)
	}
	if len(types) != 1 {
		return ""
	}
	if fieldName == "type" {
		return TypeMapping(strcase.ToCamel(fmt.Sprintf("%s_type", getEnumName(typeName))))
	}
	return types[0]
}

var funcs = template.FuncMap{
	"get_type":   getType,
	"camel":      strcase.ToCamel,
	"lowercamel": strcase.ToLowerCamel,
	"first":      getEnumName,
	"inc":        func(i int) int { return i + 1 },
	"format": func(s string, tabs int) string {
		s = strings.TrimPrefix(s, "Optional. ")
		s = strings.ReplaceAll(s, ".Example", ".\nExample")
		var (
			text    = strings.Fields(s)
			res     strings.Builder
			prefix  = "\n" + strings.Repeat("\t", tabs) + "// "
			padding = len(prefix) + 3*tabs
			count   = padding
		)
		for _, field := range text {
			// хз что тут будет с юникодом, но пофиг
			if count+len(field) >= maxInLine {
				res.WriteString(prefix)
				count = padding
			}
			if count != padding {
				res.WriteByte(' ')
				count += 1
			}
			res.WriteString(field)
			count += len(field)
		}
		return res.String()
	},
}

type Generator struct {
	schema *ApiSchema
	tmpl   *template.Template
}

func NewGenerator(schemaFile, tempaltesDir string) (*Generator, error) {
	file, err := os.Open(schemaFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var schema ApiSchema
	err = json.NewDecoder(file).Decode(&schema)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("").
		Funcs(sprig.TxtFuncMap()).
		Funcs(funcs)

	tmpl, err = tmpl.ParseGlob(filepath.Join(tempaltesDir, "*.tmpl"))
	if err != nil {
		return nil, err
	}
	return &Generator{schema: &schema, tmpl: tmpl}, nil
}

func (g *Generator) Generate(outDir string) error {
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}

	// TODO
	head := fmt.Sprintf("// Generated by: %s\n\npackage %s",
		strings.Join(
			append([]string{"go", "run", "github.com/Feresey/gen-tgbotapi/cmd/generator"},
				os.Args[1:]...,
			),
			" "),
		"api",
	)

	templateData := struct {
		Head          string
		RequiredOrder []bool
		*ApiSchema
		// map[typename][]value
		EnumTypes map[string][]string
	}{
		Head:          head,
		ApiSchema:     g.schema,
		RequiredOrder: []bool{true, false},
		EnumTypes:     g.getEnums(),
	}

	for _, tmpl := range g.tmpl.Templates() {
		name := strings.TrimSuffix(tmpl.Name(), ".tmpl")
		if name == "" {
			// main template
			continue
		}
		outName := filepath.Join(outDir, name+".go")
		log.Printf("Generating template %s", outName)

		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, templateData); err != nil {
			return err
		}

		formatted, err := imports.Process("", buf.Bytes(), nil)
		if err != nil {
			log.Print("Format generated code: ", err)
			formatted = buf.Bytes()
		}

		if err := ioutil.WriteFile(outName, formatted, 0666); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) getEnums() map[string][]string {
	res := make(map[string][]string)
	for typename, typeDesc := range g.schema.Types {
		field, ok := typeDesc.Fields["type"]
		if !ok {
			continue
		}
		enumName := strcase.ToCamel(getEnumName(typename) + "_type")
		res[enumName] = append(res[enumName], oneof(field.Description.PlainText)...)
	}
	for typename, typeDesc := range g.schema.Methods {
		field, ok := typeDesc.Arguments["type"]
		if !ok {
			continue
		}
		enumName := strcase.ToCamel(getEnumName(typename) + "_type")
		res[enumName] = append(res[enumName], oneof(field.Description.PlainText)...)
	}

	for idx, list := range res {
		res[idx] = unique(list)
	}
	return res
}

func unique(ss []string) []string {
	unique := make(map[string]struct{})
	for _, s := range ss {
		unique[s] = struct{}{}
	}
	res := make([]string, 0, len(unique))
	for s := range unique {
		res = append(res, s)
	}
	sort.Strings(res)
	return res
}

func oneof(text string) []string {
	parts := strings.Split(strings.ToLower(strings.TrimSuffix(text, ".")), "one of ")
	if len(parts) == 2 {
		return parseEntity(parts[1])
	}
	if len(parts) != 1 {
		log.Print("Not valid `oneof` retarded syntax: ", parts)
		return nil
	}

	one := "type of the result, must be "
	if strings.HasPrefix(parts[0], one) {
		return []string{strings.TrimPrefix(parts[0], one)}
	}

	if strings.Contains(parts[0], "quiz") && strings.Contains(parts[0], "regular") {
		return []string{"regular", "quiz"}
	}

	return parseEntity(parts[0])
}

func parseEntity(s string) []string {
	re := regexp.MustCompile(`"([[:word:]]+)"`)
	matches := re.FindAllStringSubmatch(s, -1)
	res := make([]string, 0, len(matches))
	for _, match := range matches {
		res = append(res, match[1])
	}
	return res
}
