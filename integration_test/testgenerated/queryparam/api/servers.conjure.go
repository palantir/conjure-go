// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/v2/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
)

type TestService interface {
	EchoQuery(ctx context.Context, inputArg string, repsArg int, optionalArg *string, listParamArg []int, lastParamArg *string, aliasStringArg AliasString, aliasAliasStringArg AliasAliasString, optionalAliasStringArg OptionalAliasString, aliasIntegerArg AliasInteger, aliasAliasIntegerArg AliasAliasInteger, optionalAliasIntegerArg OptionalAliasInteger, uuidArg uuid.UUID, setUuidArg []uuid.UUID, setAliasUuidArg []AliasUuid, aliasUuidArg AliasUuid, aliasAliasUuidArg AliasAliasUuid, optionalAliasUuidArg OptionalAliasUuid, enumArg Enum, aliasOptionalEnumArg AliasOptionalEnum, aliasEnumArg AliasEnum, optionalAliasEnumArg OptionalAliasEnum, listAliasEnumArg ListAliasEnum) (Response, error)
	EchoHeader(ctx context.Context, inputArg string, repsArg int, optionalArg *string, lastParamArg *string, aliasStringArg AliasString, aliasAliasStringArg AliasAliasString, optionalAliasStringArg OptionalAliasString, aliasIntegerArg AliasInteger, aliasAliasIntegerArg AliasAliasInteger, optionalAliasIntegerArg OptionalAliasInteger, uuidArg uuid.UUID, aliasUuidArg AliasUuid, aliasAliasUuidArg AliasAliasUuid, optionalAliasUuidArg OptionalAliasUuid, enumArg Enum, aliasOptionalEnumArg AliasOptionalEnum, aliasEnumArg AliasEnum, optionalAliasEnumArg OptionalAliasEnum) (Response, error)
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Get("EchoQuery", "/echoQuery", httpserver.NewJSONHandler(handler.HandleEchoQuery, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add echoQuery route")
	}
	if err := resource.Get("EchoHeader", "/echoHeader", httpserver.NewJSONHandler(handler.HandleEchoHeader, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add echoHeader route")
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleEchoQuery(rw http.ResponseWriter, req *http.Request) error {
	reqQuery := req.URL.Query()
	var inputArg string
	if reqQuery.Has("input") {
		inputArg = reqQuery.Get("input")
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"input\" is required")
	}
	var repsArg int
	if reqQuery.Has("reps") {
		var err error
		repsArg, err = strconv.Atoi(reqQuery.Get("reps"))
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"reps\" as integer")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"reps\" is required")
	}
	var optionalArg *string
	if reqQuery.Has("optional") {
		if optionalArgStr := reqQuery.Get("optional"); optionalArgStr != "" {
			optionalArgInternal := optionalArgStr
			optionalArg = &optionalArgInternal
		}
	}
	var listParamArg []int
	if reqQuery.Has("listParam") {
		for _, v := range reqQuery["listParam"] {
			convertedVal, err := strconv.Atoi(v)
			if err != nil {
				return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"listParam\" as integer")
			}
			listParamArg = append(listParamArg, convertedVal)
		}
	}
	var lastParamArg *string
	if reqQuery.Has("lastParam") {
		if lastParamArgStr := reqQuery.Get("lastParam"); lastParamArgStr != "" {
			lastParamArgInternal := lastParamArgStr
			lastParamArg = &lastParamArgInternal
		}
	}
	var aliasStringArg AliasString
	if reqQuery.Has("aliasString") {
		aliasStringArg = AliasString(reqQuery.Get("aliasString"))
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasString\" is required")
	}
	var aliasAliasStringArg AliasAliasString
	if reqQuery.Has("aliasAliasString") {
		aliasAliasStringArg = AliasAliasString(reqQuery.Get("aliasAliasString"))
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasAliasString\" is required")
	}
	var optionalAliasStringArg OptionalAliasString
	if reqQuery.Has("optionalAliasString") {
		if err := optionalAliasStringArg.UnmarshalText([]byte(reqQuery.Get("optionalAliasString"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasString\" as OptionalAliasString")
		}
	}
	var aliasIntegerArg AliasInteger
	if reqQuery.Has("aliasInteger") {
		if err := safejson.Unmarshal([]byte(reqQuery.Get("aliasInteger")), &aliasIntegerArg); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasInteger\" as AliasInteger")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasInteger\" is required")
	}
	var aliasAliasIntegerArg AliasAliasInteger
	if reqQuery.Has("aliasAliasInteger") {
		if err := safejson.Unmarshal([]byte(reqQuery.Get("aliasAliasInteger")), &aliasAliasIntegerArg); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasAliasInteger\" as AliasAliasInteger")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasAliasInteger\" is required")
	}
	var optionalAliasIntegerArg OptionalAliasInteger
	if reqQuery.Has("optionalAliasInteger") {
		if err := safejson.Unmarshal([]byte(reqQuery.Get("optionalAliasInteger")), &optionalAliasIntegerArg); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasInteger\" as OptionalAliasInteger")
		}
	}
	var uuidArg uuid.UUID
	if reqQuery.Has("uuid") {
		var err error
		uuidArg, err = uuid.ParseUUID(reqQuery.Get("uuid"))
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"uuid\" as uuid")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"uuid\" is required")
	}
	var setUuidArg []uuid.UUID
	if reqQuery.Has("setUuid") {
		for _, v := range reqQuery["setUuid"] {
			convertedVal, err := uuid.ParseUUID(v)
			if err != nil {
				return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"setUuid\" as uuid")
			}
			setUuidArg = append(setUuidArg, convertedVal)
		}
	}
	var setAliasUuidArg []AliasUuid
	if reqQuery.Has("setAliasUuid") {
		for _, v := range reqQuery["setAliasUuid"] {
			var convertedVal AliasUuid
			if err := convertedVal.UnmarshalText([]byte(v)); err != nil {
				return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"setAliasUuid\" as AliasUuid")
			}
			setAliasUuidArg = append(setAliasUuidArg, convertedVal)
		}
	}
	var aliasUuidArg AliasUuid
	if reqQuery.Has("aliasUuid") {
		if err := aliasUuidArg.UnmarshalText([]byte(reqQuery.Get("aliasUuid"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasUuid\" as AliasUuid")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasUuid\" is required")
	}
	var aliasAliasUuidArg AliasAliasUuid
	if reqQuery.Has("aliasAliasUuid") {
		if err := aliasAliasUuidArg.UnmarshalText([]byte(reqQuery.Get("aliasAliasUuid"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasAliasUuid\" as AliasAliasUuid")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasAliasUuid\" is required")
	}
	var optionalAliasUuidArg OptionalAliasUuid
	if reqQuery.Has("optionalAliasUuid") {
		if err := optionalAliasUuidArg.UnmarshalText([]byte(reqQuery.Get("optionalAliasUuid"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasUuid\" as OptionalAliasUuid")
		}
	}
	var enumArg Enum
	if reqQuery.Has("enum") {
		if err := enumArg.UnmarshalText([]byte(reqQuery.Get("enum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"enum\" as Enum")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"enum\" is required")
	}
	var aliasOptionalEnumArg AliasOptionalEnum
	if reqQuery.Has("aliasOptionalEnum") {
		if err := aliasOptionalEnumArg.UnmarshalText([]byte(reqQuery.Get("aliasOptionalEnum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasOptionalEnum\" as AliasOptionalEnum")
		}
	}
	var aliasEnumArg AliasEnum
	if reqQuery.Has("aliasEnum") {
		if err := aliasEnumArg.UnmarshalText([]byte(reqQuery.Get("aliasEnum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasEnum\" as AliasEnum")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "query parameter \"aliasEnum\" is required")
	}
	var optionalAliasEnumArg OptionalAliasEnum
	if reqQuery.Has("optionalAliasEnum") {
		if err := optionalAliasEnumArg.UnmarshalText([]byte(reqQuery.Get("optionalAliasEnum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasEnum\" as OptionalAliasEnum")
		}
	}
	var listAliasEnumArg ListAliasEnum
	if reqQuery.Has("listAliasEnum") {
		for _, v := range reqQuery["listAliasEnum"] {
			var convertedVal AliasEnum
			if err := convertedVal.UnmarshalText([]byte(v)); err != nil {
				return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"listAliasEnum\" as AliasEnum")
			}
			listAliasEnumArg = append(listAliasEnumArg, convertedVal)
		}
	}
	respArg, err := t.impl.EchoQuery(req.Context(), inputArg, repsArg, optionalArg, listParamArg, lastParamArg, aliasStringArg, aliasAliasStringArg, optionalAliasStringArg, aliasIntegerArg, aliasAliasIntegerArg, optionalAliasIntegerArg, uuidArg, setUuidArg, setAliasUuidArg, aliasUuidArg, aliasAliasUuidArg, optionalAliasUuidArg, enumArg, aliasOptionalEnumArg, aliasEnumArg, optionalAliasEnumArg, listAliasEnumArg)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleEchoHeader(rw http.ResponseWriter, req *http.Request) error {
	var inputArg string
	if len(req.Header.Values("Input-String")) > 0 {
		inputArg = req.Header.Get("Input-String")
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Input-String\" is required")
	}
	var repsArg int
	if len(req.Header.Values("Reps-Integer")) > 0 {
		var err error
		repsArg, err = strconv.Atoi(req.Header.Get("Reps-Integer"))
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"reps\" as integer")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Reps-Integer\" is required")
	}
	var optionalArg *string
	if len(req.Header.Values("Optional-String")) > 0 {
		if optionalArgStr := req.Header.Get("Optional-String"); optionalArgStr != "" {
			optionalArgInternal := optionalArgStr
			optionalArg = &optionalArgInternal
		}
	}
	var lastParamArg *string
	if len(req.Header.Values("Last-Param")) > 0 {
		if lastParamArgStr := req.Header.Get("Last-Param"); lastParamArgStr != "" {
			lastParamArgInternal := lastParamArgStr
			lastParamArg = &lastParamArgInternal
		}
	}
	var aliasStringArg AliasString
	if len(req.Header.Values("Alias-String")) > 0 {
		aliasStringArg = AliasString(req.Header.Get("Alias-String"))
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-String\" is required")
	}
	var aliasAliasStringArg AliasAliasString
	if len(req.Header.Values("Alias-Alias-String")) > 0 {
		aliasAliasStringArg = AliasAliasString(req.Header.Get("Alias-Alias-String"))
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-Alias-String\" is required")
	}
	var optionalAliasStringArg OptionalAliasString
	if len(req.Header.Values("Optional-Alias-String")) > 0 {
		if err := optionalAliasStringArg.UnmarshalText([]byte(req.Header.Get("Optional-Alias-String"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasString\" as OptionalAliasString")
		}
	}
	var aliasIntegerArg AliasInteger
	if len(req.Header.Values("Alias-Integer")) > 0 {
		if err := safejson.Unmarshal([]byte(req.Header.Get("Alias-Integer")), &aliasIntegerArg); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasInteger\" as AliasInteger")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-Integer\" is required")
	}
	var aliasAliasIntegerArg AliasAliasInteger
	if len(req.Header.Values("Alias-Alias-Integer")) > 0 {
		if err := safejson.Unmarshal([]byte(req.Header.Get("Alias-Alias-Integer")), &aliasAliasIntegerArg); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasAliasInteger\" as AliasAliasInteger")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-Alias-Integer\" is required")
	}
	var optionalAliasIntegerArg OptionalAliasInteger
	if len(req.Header.Values("Optional-Alias-Integer")) > 0 {
		if err := safejson.Unmarshal([]byte(req.Header.Get("Optional-Alias-Integer")), &optionalAliasIntegerArg); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasInteger\" as OptionalAliasInteger")
		}
	}
	var uuidArg uuid.UUID
	if len(req.Header.Values("Uuid-Value")) > 0 {
		var err error
		uuidArg, err = uuid.ParseUUID(req.Header.Get("Uuid-Value"))
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"uuid\" as uuid")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Uuid-Value\" is required")
	}
	var aliasUuidArg AliasUuid
	if len(req.Header.Values("Alias-Uuid")) > 0 {
		if err := aliasUuidArg.UnmarshalText([]byte(req.Header.Get("Alias-Uuid"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasUuid\" as AliasUuid")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-Uuid\" is required")
	}
	var aliasAliasUuidArg AliasAliasUuid
	if len(req.Header.Values("Alias-Alias-Uuid")) > 0 {
		if err := aliasAliasUuidArg.UnmarshalText([]byte(req.Header.Get("Alias-Alias-Uuid"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasAliasUuid\" as AliasAliasUuid")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-Alias-Uuid\" is required")
	}
	var optionalAliasUuidArg OptionalAliasUuid
	if len(req.Header.Values("Optional-Alias-Uuid")) > 0 {
		if err := optionalAliasUuidArg.UnmarshalText([]byte(req.Header.Get("Optional-Alias-Uuid"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasUuid\" as OptionalAliasUuid")
		}
	}
	var enumArg Enum
	if len(req.Header.Values("Enum-Value")) > 0 {
		if err := enumArg.UnmarshalText([]byte(req.Header.Get("Enum-Value"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"enum\" as Enum")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Enum-Value\" is required")
	}
	var aliasOptionalEnumArg AliasOptionalEnum
	if len(req.Header.Values("Alias-Optional-Enum")) > 0 {
		if err := aliasOptionalEnumArg.UnmarshalText([]byte(req.Header.Get("Alias-Optional-Enum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasOptionalEnum\" as AliasOptionalEnum")
		}
	}
	var aliasEnumArg AliasEnum
	if len(req.Header.Values("Alias-Enum")) > 0 {
		if err := aliasEnumArg.UnmarshalText([]byte(req.Header.Get("Alias-Enum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"aliasEnum\" as AliasEnum")
		}
	} else {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "header parameter \"Alias-Enum\" is required")
	}
	var optionalAliasEnumArg OptionalAliasEnum
	if len(req.Header.Values("Optional-Alias-Enum")) > 0 {
		if err := optionalAliasEnumArg.UnmarshalText([]byte(req.Header.Get("Optional-Alias-Enum"))); err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"optionalAliasEnum\" as OptionalAliasEnum")
		}
	}
	respArg, err := t.impl.EchoHeader(req.Context(), inputArg, repsArg, optionalArg, lastParamArg, aliasStringArg, aliasAliasStringArg, optionalAliasStringArg, aliasIntegerArg, aliasAliasIntegerArg, optionalAliasIntegerArg, uuidArg, aliasUuidArg, aliasAliasUuidArg, optionalAliasUuidArg, enumArg, aliasOptionalEnumArg, aliasEnumArg, optionalAliasEnumArg)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}
