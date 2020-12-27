{{.Head}}
{{range $method, $desc := .Methods}}
{{- if gt (len $desc.Arguments) 2 }}
// {{camel $method}}
// {{format $desc.Description.PlainText 0}}
type {{camel $method}}Config struct {
{{- range $idx, $put_required := $.RequiredOrder}}
{{- range $argname, $arg := $desc.Arguments}}
	{{- $type := get_type $argname $method $arg.Types }}

	{{- if eq $arg.Required $put_required}}
	// {{camel $argname}}
	// {{format $arg.Description.PlainText 1}}
	{{camel $argname}} {{if and (not $type.IsSimpleType) (not $arg.Required) (not $type.IsArray) -}}*{{end -}}
		{{- $type.GoType}} `json:"{{$argname}}{{if not $arg.Required}},omitempty{{end}}"`
	{{- end}}
{{- end}}
{{- end}}
}
{{- end}}

{{- if is_sendable $desc.Arguments}}
func (t {{camel $method}}Config) EncodeURL() (url.Values, error) {
	res := make(url.Values)
{{- range $argname, $arg := $desc.Arguments}}
	{{- $type := (get_type $argname $method $arg.Types)}}

	{{- if not (eq $type.GoType "InputFile")}}
		{{- if and $type.IsSimpleType (not (is_interface $type)) }}
	res.Add("{{$argname}}", {{format_url (print "t." (camel $argname)) false $type}})
		{{- else}}
	if t.{{camel $argname}} != nil {
			{{- if or (is_interface $type) (not $type.IsSimpleType) }}
	raw, err := json.Marshal({{print "t." (camel $argname)}})
	if err != nil {
		return nil, err
	}
	res.Add("{{$argname}}", string(raw))
			{{- else}}
	res.Add("{{$argname}}", {{format_url (print "t." (camel $argname)) true $type}})
			{{- end}}
	}
		{{- end}}
	{{- end}}
{{- end}}
	return res, nil
}
{{- end}}

// {{camel $method}}
// {{format $desc.Description.PlainText 0}}
{{- $input_file := ""}}
{{- $second := ""}}
{{- range $argname, $arg := $desc.Arguments}}
	{{- $type := (get_type $argname $method $arg.Types)}}
	{{- if and $arg.Required (eq $type.GoType "InputFile")}}
		{{- $input_file = $argname}}}}
	{{- else}}
		{{- $second = $argname}}
	{{- end}}
{{- end}}
func (api *API) {{camel $method}}(
	ctx context.Context,

{{- if gt (len $desc.Arguments) 2}}
	args *{{camel $method}}Config,
{{else}}
	{{range $argname, $arg := $desc.Arguments -}}
		// {{if not $arg.Required}}not {{end}}required.
		// {{format $arg.Description.PlainText 2}}
		{{- $type := (get_type $argname $method $arg.Types)}}
		{{lowercamel $argname}} {{if and (not $arg.Required) (not $type.IsArray)}}*{{end -}}
	{{$type.GoType}},
	{{end}}
{{- end -}}

) (
{{- $returns := false}}
{{- $return_stared := false}}
{{- with $desc.Returns}}
	{{- if (not (eq .GoType "True"))}}
		{{- if and (not .IsSimpleType) (not .IsArray)}}*{{$return_stared = true}}{{end}}
	{{- .GoType}},
		{{- $returns = true}}
	{{- end}}
{{- end -}}
error) {
{{- if not (empty $input_file)}}
{{- if gt (len $desc.Arguments) 2}}
	{{- $input_file = print "args." (camel $input_file)}}
{{- else}}
	{{- $input_file = lowercamel $input_file}}
{{- end}}
	if {{$input_file}}.Reader != nil {
		{{- if gt (len $desc.Arguments) 2}}
		values, err := args.EncodeURL()
		if err != nil {
			return {{if $returns}}nil,{{end}} err
		}
		{{- else}}
		values := url.Values{
			"{{$second}}" : []string{ {{- format_url (lowercamel $second) false (get_type $second $method (index $desc.Arguments $second).Types) -}} },
		}
		{{- end}}
		{{if $returns}}resp{{else}}_{{end}}, err := api.UploadFile(ctx, values, "{{$method}}","{{input_type $method}}", &{{$input_file}})
		{{- if not $returns}}
		return err
		{{- else}}
		if err != nil {
			return nil, err
		}

		var res {{$desc.Returns.GoType}}
		err = json.Unmarshal(resp.Result, &res)
		return {{if $return_stared}}&{{end}}res, err
		{{- end}}
	}
{{- end}}

{{- if not (gt (len $desc.Arguments) 2)}}
	{{- if gt (len $desc.Arguments) 0}}
	args := map[string]interface{} {
		{{- range $argname, $arg := $desc.Arguments}}
		"{{$argname}}" : {{lowercamel $argname}},
		{{- end}}
	}
	{{- end}}
{{- end}}
	{{if $returns}}resp{{else}}_{{end}}, err := api.MakeRequest(ctx, "{{$method}}", {{if eq (len $desc.Arguments) 0}}nil{{else}}args{{end}})
	{{- if not $returns}}
	return err
	{{- else}}
	if err != nil {
		return {{if $returns}}{{default_return $desc.Returns}}, {{end}}err
	}
	var data {{$desc.Returns.GoType}}
	err = json.Unmarshal(resp.Result, &data)
	return {{if $return_stared}}&{{end}}data, err
	{{- end}}
}
{{end}}