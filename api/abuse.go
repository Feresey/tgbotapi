package api

import (
	"encoding/json"
	"strconv"
	"strings"
)

type True struct{}

func (True) MarshalText() ([]byte, error) { return []byte("true"), nil }
func (*True) UnmarshalText([]byte) error  { return nil }

// one of InputMediaPhoto or InputMediaVideo
type InputMediaGraphics interface{}

type IntStr struct {
	Int    int64
	String string
}

func (i IntStr) MarshalText() ([]byte, error) {
	if strings.HasPrefix(i.String, "@") {
		return []byte(i.String), nil
	}
	return []byte(strconv.FormatInt(i.Int, 10)), nil
}

type InputDataType struct {
	File   *InputFile
	String string
}

func (i InputDataType) MarshalJSON() ([]byte, error) {
	if i.File != nil {
		return json.Marshal(i.File)
	}
	return []byte(i.String), nil
}

// one of
/*
[
    "InlineKeyboardMarkup",
    "ReplyKeyboardMarkup",
    "ReplyKeyboardRemove",
    "ForceReply"
]
*/
// мне влом писать, сами думайте.
type ReplyMarkup interface{}
