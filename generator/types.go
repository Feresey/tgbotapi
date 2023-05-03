package generator

import "strings"

type TypeMapping string

// cat ../schema/public/types.json | jq '[.[].fields | to_entries | .[].value.types[]] | unique'

func (t TypeMapping) GoType() string {
	if t.IsArray() {
		return t.SplitArray()
	}
	//nolint:goconst // useless
	switch t {
	case "str", "String":
		return "string"
	case "int", "Int":
		// Ну на самом деле chat_id должен быть int64, а message_id и int хватит.
		return "int64"
	case "True":
		return "bool"
	case "Float", "float":
		return "float64"
	case "InputMediaAudio, InputMediaDocument, InputMediaPhoto and InputMediaVideo":
		return "InputMediaGraphics"
	// retarded parsing
	case "Messages":
		return " Message"
	}
	return rename(string(t))
}

func (t TypeMapping) IsSimpleType() bool {
	switch t {
	case "str", "String", "int", "bool", "IntStr":
		return true
	case "CaptionEntities":
		return false
	default:
		return isInterface(t)
	}
}

func (t TypeMapping) IsArray() bool { return strings.HasPrefix(string(t), "array(") }

func (t TypeMapping) ArrayType() TypeMapping {
	return TypeMapping(strings.TrimSuffix(strings.TrimPrefix(string(t), "array("), ")"))
}

func (t TypeMapping) SplitArray() string {
	var res strings.Builder
	for t.IsArray() {
		res.WriteString("[]")
		t = t.ArrayType()
	}
	return res.String() + t.GoType()
}

type Field struct {
	Types       []TypeMapping `json:"types"`
	Description Description   `json:"description"`
	Required    bool          `json:"required"`
}

type Description struct {
	PlainText string `json:"plaintext"`
	Markdown  string `json:"markdown"`
	HTML      string `json:"html"`
}

type Type struct {
	Fields      map[string]Field `json:"fields"`
	Description Description      `json:"description"`
	Category    string           `json:"category"`
}

type Method struct {
	Arguments map[string]Field `json:"arguments"`
	// null available
	Returns     *TypeMapping `json:"returns"`
	Description Description  `json:"description"`
	Category    string       `json:"category"`
}
type Article struct {
	Description
	Title    string `json:"title"`
	Category string `json:"category"`
}

type APISchema struct {
	Articles map[string]Article `json:"articles"`
	Methods  map[string]Method  `json:"methods"`
	Types    map[string]Type    `json:"types"`
	Version  string             `json:"version"`

	// ну на это мне как-то посрать

	BuildInfo  map[string]interface{} `json:"build_info"`
	Changelogs map[string]interface{} `json:"changelogs"`
}
