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

func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels(in *jlexer.Lexer, out *RegReqWithUserRespList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(RegReqWithUserRespList, 0, 0)
			} else {
				*out = RegReqWithUserRespList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 RegReqWithUserResp
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels(out *jwriter.Writer, in RegReqWithUserRespList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v RegReqWithUserRespList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegReqWithUserRespList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegReqWithUserRespList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegReqWithUserRespList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels1(in *jlexer.Lexer, out *RegReqWithUserResp) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "user":
			(out.User).UnmarshalEasyJSON(in)
		case "manager":
			if in.IsNull() {
				in.Skip()
				out.Manager = nil
			} else {
				if out.Manager == nil {
					out.Manager = new(UserRes)
				}
				(*out.Manager).UnmarshalEasyJSON(in)
			}
		case "type":
			out.Type = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "time_in_queue":
			out.TimeInQueue = uint32(in.Uint32())
		case "create_time":
			out.CreateTime = uint64(in.Uint64())
		case "message":
			out.Message = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels1(out *jwriter.Writer, in RegReqWithUserResp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		(in.User).MarshalEasyJSON(out)
	}
	if in.Manager != nil {
		const prefix string = ",\"manager\":"
		out.RawString(prefix)
		(*in.Manager).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"time_in_queue\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.TimeInQueue))
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CreateTime))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegReqWithUserResp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegReqWithUserResp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegReqWithUserResp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegReqWithUserResp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels1(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels2(in *jlexer.Lexer, out *RegReqWithUser) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "user":
			(out.User).UnmarshalEasyJSON(in)
		case "manager":
			if in.IsNull() {
				in.Skip()
				out.Manager = nil
			} else {
				if out.Manager == nil {
					out.Manager = new(User)
				}
				(*out.Manager).UnmarshalEasyJSON(in)
			}
		case "type":
			out.Type = RegReqType(in.Int8())
		case "status":
			out.Status = string(in.String())
		case "time_in_queue":
			out.TimeInQueue = uint32(in.Uint32())
		case "create_time":
			out.CreateTime = uint64(in.Uint64())
		case "message":
			out.Message = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels2(out *jwriter.Writer, in RegReqWithUser) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		(in.User).MarshalEasyJSON(out)
	}
	if in.Manager != nil {
		const prefix string = ",\"manager\":"
		out.RawString(prefix)
		(*in.Manager).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.Int8(int8(in.Type))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"time_in_queue\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.TimeInQueue))
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CreateTime))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegReqWithUser) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegReqWithUser) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegReqWithUser) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegReqWithUser) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels2(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels3(in *jlexer.Lexer, out *RegReqRespList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(RegReqRespList, 0, 0)
			} else {
				*out = RegReqRespList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 RegReqResp
			(v4).UnmarshalEasyJSON(in)
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels3(out *jwriter.Writer, in RegReqRespList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			(v6).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v RegReqRespList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegReqRespList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegReqRespList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegReqRespList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels3(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels4(in *jlexer.Lexer, out *RegReqResp) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "user_id":
			out.UserID = uint64(in.Uint64())
		case "manager_id":
			out.ManagerID = uint64(in.Uint64())
		case "type":
			out.Type = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "create_time":
			out.CreateTime = uint64(in.Uint64())
		case "message":
			out.Message = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels4(out *jwriter.Writer, in RegReqResp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.UserID))
	}
	if in.ManagerID != 0 {
		const prefix string = ",\"manager_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.ManagerID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CreateTime))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegReqResp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegReqResp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegReqResp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegReqResp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels4(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels5(in *jlexer.Lexer, out *RegReqFull) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "user_id":
			out.UserID = uint64(in.Uint64())
		case "manager_id":
			out.ManagerID = uint64(in.Uint64())
		case "type":
			out.Type = RegReqType(in.Int8())
		case "status":
			out.Status = string(in.String())
		case "create_time":
			out.CreateTime = uint64(in.Uint64())
		case "message":
			out.Message = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels5(out *jwriter.Writer, in RegReqFull) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.UserID))
	}
	if in.ManagerID != 0 {
		const prefix string = ",\"manager_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.ManagerID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.Int8(int8(in.Type))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CreateTime))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegReqFull) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegReqFull) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegReqFull) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegReqFull) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels5(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels6(in *jlexer.Lexer, out *ParentPassportReq) {
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
		case "passport":
			out.Passport = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels6(out *jwriter.Writer, in ParentPassportReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"passport\":"
		out.RawString(prefix[1:])
		out.String(string(in.Passport))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ParentPassportReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ParentPassportReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ParentPassportReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ParentPassportReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels6(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels7(in *jlexer.Lexer, out *FixParentPassportReq) {
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
		case "req_id":
			out.ReqID = uint64(in.Uint64())
		case "passport":
			out.Passport = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels7(out *jwriter.Writer, in FixParentPassportReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"req_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ReqID))
	}
	{
		const prefix string = ",\"passport\":"
		out.RawString(prefix)
		out.String(string(in.Passport))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FixParentPassportReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FixParentPassportReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FixParentPassportReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FixParentPassportReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels7(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels8(in *jlexer.Lexer, out *FixChildThirdRegReq) {
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
		case "req_id":
			out.ReqID = uint64(in.Uint64())
		case "child":
			(out.Child).UnmarshalEasyJSON(in)
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels8(out *jwriter.Writer, in FixChildThirdRegReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"req_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ReqID))
	}
	{
		const prefix string = ",\"child\":"
		out.RawString(prefix)
		(in.Child).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FixChildThirdRegReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FixChildThirdRegReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FixChildThirdRegReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FixChildThirdRegReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels8(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels9(in *jlexer.Lexer, out *FixChildSecondRegReq) {
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
		case "req_id":
			out.ReqID = uint64(in.Uint64())
		case "child":
			(out.Child).UnmarshalEasyJSON(in)
		case "relationship":
			out.Relationship = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels9(out *jwriter.Writer, in FixChildSecondRegReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"req_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ReqID))
	}
	{
		const prefix string = ",\"child\":"
		out.RawString(prefix)
		(in.Child).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"relationship\":"
		out.RawString(prefix)
		out.String(string(in.Relationship))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FixChildSecondRegReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FixChildSecondRegReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FixChildSecondRegReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FixChildSecondRegReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels9(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels10(in *jlexer.Lexer, out *FixChildFirstRegReq) {
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
		case "req_id":
			out.ReqID = uint64(in.Uint64())
		case "child":
			(out.Child).UnmarshalEasyJSON(in)
		case "is_student":
			out.IsStudent = bool(in.Bool())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels10(out *jwriter.Writer, in FixChildFirstRegReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"req_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ReqID))
	}
	{
		const prefix string = ",\"child\":"
		out.RawString(prefix)
		(in.Child).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"is_student\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsStudent))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FixChildFirstRegReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FixChildFirstRegReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FixChildFirstRegReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FixChildFirstRegReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels10(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels11(in *jlexer.Lexer, out *FailedReq) {
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
		case "req_id":
			out.ReqId = uint64(in.Uint64())
		case "failed_message":
			out.FailedMessage = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels11(out *jwriter.Writer, in FailedReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"req_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ReqId))
	}
	{
		const prefix string = ",\"failed_message\":"
		out.RawString(prefix)
		out.String(string(in.FailedMessage))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FailedReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FailedReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FailedReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FailedReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels11(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels12(in *jlexer.Lexer, out *ChildThirdRegReq) {
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
		case "child":
			(out.Child).UnmarshalEasyJSON(in)
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels12(out *jwriter.Writer, in ChildThirdRegReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"child\":"
		out.RawString(prefix[1:])
		(in.Child).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChildThirdRegReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChildThirdRegReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChildThirdRegReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChildThirdRegReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels12(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels13(in *jlexer.Lexer, out *ChildSecondRegReq) {
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
		case "child":
			(out.Child).UnmarshalEasyJSON(in)
		case "relationship":
			out.Relationship = string(in.String())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels13(out *jwriter.Writer, in ChildSecondRegReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"child\":"
		out.RawString(prefix[1:])
		(in.Child).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"relationship\":"
		out.RawString(prefix)
		out.String(string(in.Relationship))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChildSecondRegReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels13(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChildSecondRegReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels13(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChildSecondRegReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels13(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChildSecondRegReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels13(l, v)
}
func easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels14(in *jlexer.Lexer, out *ChildFirstRegReq) {
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
		case "child":
			(out.Child).UnmarshalEasyJSON(in)
		case "is_student":
			out.IsStudent = bool(in.Bool())
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
func easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels14(out *jwriter.Writer, in ChildFirstRegReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"child\":"
		out.RawString(prefix[1:])
		(in.Child).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"is_student\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsStudent))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChildFirstRegReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels14(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChildFirstRegReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3c9d2b01EncodeGithubComVoyakinHLokleBackendInternalModels14(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChildFirstRegReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels14(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChildFirstRegReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3c9d2b01DecodeGithubComVoyakinHLokleBackendInternalModels14(l, v)
}
