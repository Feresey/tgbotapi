package generator

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeMapping_GoType(t *testing.T) {
	tests := []struct {
		t    string
		want string
	}{
		{
			t:    "int",
			want: "int64",
		},
		{
			t:    "str",
			want: "string",
		},
		{
			t:    "not-standart-type",
			want: "not-standart-type",
		},
		{
			t:    "array(type)",
			want: "[]type",
		},
		{
			t:    "array(array(array(array)))",
			want: "[][][]array",
		},
		{
			t:    "array(array(PhotoSize))",
			want: "[][]PhotoSize",
		},
	}
	for _, tt := range tests {
		got := TypeMapping(tt.t).GoType()
		require.Equal(t, tt.want, got)
	}
}

func TestTypes(t *testing.T) {
	raw, err := ioutil.ReadFile("../schema/public/types.json")
	if err != nil {
		println(err)
		println("no file == no test")
		return
	}

	var decoded map[string]Type
	err = json.Unmarshal(raw, &decoded)
	require.NoError(t, err)

	reEncoded, err := json.Marshal(decoded)
	require.NoError(t, err)

	require.JSONEq(t, string(raw), string(reEncoded))
}

func TestMethods(t *testing.T) {
	raw, err := ioutil.ReadFile("../schema/public/methods.json")
	if err != nil {
		println(err)
		println("no file == no test")
		return
	}

	var decoded map[string]Method
	err = json.Unmarshal(raw, &decoded)
	require.NoError(t, err)

	reEncoded, err := json.Marshal(decoded)
	require.NoError(t, err)

	require.JSONEq(t, string(raw), string(reEncoded))
}

func TestArticles(t *testing.T) {
	raw, err := ioutil.ReadFile("../schema/public/articles.json")
	if err != nil {
		println(err)
		println("no file == no test")
		return
	}

	var decoded map[string]Article
	err = json.Unmarshal(raw, &decoded)
	require.NoError(t, err)

	reEncoded, err := json.Marshal(decoded)
	require.NoError(t, err)

	require.JSONEq(t, string(raw), string(reEncoded))
}
