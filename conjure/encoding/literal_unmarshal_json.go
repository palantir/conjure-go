package encoding

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
)

const (
	dataName = "data"
)

func UnmarshalJSONMethodBody(methodBody *jen.Group, receiverName, receiverType string, strict bool) *jen.Statement {
	methodBody.Id("ctx").Op(":=").Add(snip.ContextTODO()).Call()
	methodBody.If(jen.Op("!").Add(snip.GJSONValidBytes()).Call(jen.Id(dataName))).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"),
			jen.Lit(fmt.Sprintf("invalid JSON for %s", receiverType)),
		)),
	)
}
