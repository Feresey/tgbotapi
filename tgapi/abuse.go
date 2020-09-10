package tgapi

import (
	"encoding/json"
	"fmt"
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

func NewInt(i int64) IntStr {
	return IntStr{Int: i}
}

func NewStr(s string) IntStr {
	return IntStr{String: s}
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

func (t *MessageEntity) IsCommand() bool {
	return t.Type == EntityTypeBotCommand
}

// IsCommand returns true if message starts with a "bot_command" entity.
func (t *Message) IsCommand() bool {
	if t.Entities == nil || len(t.Entities) == 0 {
		return false
	}

	entity := t.Entities[0]
	return entity.Offset == 0 && entity.IsCommand()
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is removed. Use
// CommandWithAt() if you do not want that.
func (t *Message) Command() string {
	command := t.CommandWithAt()

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandWithAt checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is not removed. Use Command()
// if you want that.
func (t *Message) CommandWithAt() string {
	if !t.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := t.Entities[0]
	return t.GetText()[1:entity.Length]
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
//
// Note: The first character after the command name is omitted:
// - "/foo bar baz" yields "bar baz", not " bar baz"
// - "/foo-bar baz" yields "bar baz", too
// Even though the latter is not a command conforming to the spec, the API
// marks "/foo" as command entity.
func (t *Message) CommandArguments() string {
	if !t.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := t.Entities[0]
	if int64(len(t.GetText())) == entity.Length {
		return "" // The command makes up the whole message
	}

	return t.GetText()[entity.Length+1:]
}

// String displays a simple text version of a user.
//
// It is normally a user's username, but falls back to a first/last
// name as available.
func (t *User) String() string {
	if t == nil {
		return ""
	}
	if t.Username != nil {
		return t.GetUsername()
	}

	name := t.FirstName
	if t.LastName != nil {
		name += " " + t.GetLastName()
	}

	return name
}

func AskUser(user *User) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", user.String(), user.ID)
}
