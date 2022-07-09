// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9e1087fdDecodeGithubComVoyakinHLokleBackendInternalModels(in *jlexer.Lexer, out *Credentials) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "email":
			out.Email = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9e1087fdEncodeGithubComVoyakinHLokleBackendInternalModels(out *jwriter.Writer, in Credentials) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Credentials) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComVoyakinHLokleBackendInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Credentials) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComVoyakinHLokleBackendInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Credentials) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComVoyakinHLokleBackendInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Credentials) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComVoyakinHLokleBackendInternalModels(l, v)
}
