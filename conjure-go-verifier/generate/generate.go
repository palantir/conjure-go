// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/pkg/errors"
)

const conjureVerifierVersion = "0.18.5"

func main() {
	const versionFilePath = "version_test.go"
	newVersionFileContent := fmt.Sprintf(`// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a generated file: do not edit by hand.
// To update this file, run the generator in conjure-go-verifier.
package verifier_test

const verificationServerVersion = "%s"
`, conjureVerifierVersion)

	const clientVerificationAPIFile = "verification-server-api.conjure.json"
	const clientTestCasesFile = "verification-server-test-cases.json"

	// if version file exists and is in desired state, assume that all downloaded content is in desired state
	if currVersionFileContent, err := ioutil.ReadFile(versionFilePath); err != nil || string(currVersionFileContent) != newVersionFileContent {
		if err := downloadFile(clientTestCasesFile, fmt.Sprintf("https://palantir.bintray.com/releases/com/palantir/conjure/verification/verification-server-test-cases/%s/verification-server-test-cases-%s.json", conjureVerifierVersion, conjureVerifierVersion)); err != nil {
			panic(err)
		}
		if err := downloadFile(clientVerificationAPIFile, fmt.Sprintf("https://palantir.bintray.com/releases/com/palantir/conjure/verification/verification-server-api/%s/verification-server-api-%s.conjure.json", conjureVerifierVersion, conjureVerifierVersion)); err != nil {
			panic(err)
		}
		// update version in circle.yml
		if err := updateVersionInCircleConfig("../.circleci/config.yml", conjureVerifierVersion); err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(versionFilePath, []byte(newVersionFileContent), 0644); err != nil {
			panic(err)
		}
	}
	conjureDefinition, err := conjure.FromIRFile(clientVerificationAPIFile)
	if err != nil {
		panic(err)
	}
	if err := conjure.Generate(conjureDefinition, conjure.OutputConfiguration{
		GenerateServer:      true,
		LiteralJSONMethods:  true,
		OutputDir:           ".",
		GenerateYAMLMethods: true,
	}); err != nil {
		panic(err)
	}

	if err := generateServerImpl(".", conjureDefinition); err != nil {
		panic(err)
	}
}

var dockerImgRegexp = regexp.MustCompile(`(?m)^(.*?- image: palantirtechnologies/conjure-verification-server):(.+)$`)

func updateVersionInCircleConfig(cfgFile, version string) error {
	bytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return errors.WithStack(err)
	}
	cfgContent := string(bytes)

	updated := dockerImgRegexp.ReplaceAllString(cfgContent, fmt.Sprintf("$1:%s", version))
	if cfgContent == updated {
		return nil
	}
	if err := ioutil.WriteFile(cfgFile, []byte(updated), 0644); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GET %s received bad response status: %s", url, resp.Status)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

const (
	serverImport = "github.com/palantir/conjure-go/v6/conjure-go-verifier/conjure/verification/server"
)

func generateServerImpl(outputDir string, conjureDefinition spec.ConjureDefinition) error {
	def, err := types.NewConjureDefinition(outputDir, conjureDefinition)
	if err != nil {
		return err
	}

	// Declare impl type
	for _, pkgDef := range def.Packages {
		file := jen.NewFilePath(pkgDef.ImportPath)
		if len(pkgDef.Services) == 0 {
			continue
		}
		for _, serviceDef := range pkgDef.Services {
			if serviceDef.Name == "AutoDeserializeService" {
				generateAutoDeserializeServiceServerImpl(file.Group, serviceDef)
				buf := &bytes.Buffer{}
				if err := file.Render(buf); err != nil {
					return err
				}
				path := filepath.Join(pkgDef.OutputDir, "autodeserialize_service_impl.go")
				if err := ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return errors.New("AutoDeserializeService not found")
}

func generateAutoDeserializeServiceServerImpl(file *jen.Group, serviceDef *types.ServiceDefinition) {
	file.Type().Id("AutoDeserializeServiceImpl").Struct(
		jen.Id("TestCases").Map(jen.Qual(serverImport, "EndpointName")).Qual(serverImport, "PositiveAndNegativeTestCases"),
	)
	for _, endpointDef := range serviceDef.Endpoints {
		file.Func().
			Params(jen.Id("a").Id("AutoDeserializeServiceImpl")).
			Id(transforms.Export(endpointDef.EndpointName)).
			Params(jen.Id("_").Add(snip.Context()), jen.Id("indexArg").Int()).
			Params((*endpointDef.Returns).Code(), jen.Error()).
			Block(
				jen.Var().Id("value").Add((*endpointDef.Returns).Code()),
				jen.Err().Op(":=").Id("value").Dot("UnmarshalJSON").Call(
					jen.Id("a").Dot("testCaseBytes").Call(jen.Lit(endpointDef.EndpointName), jen.Id("indexArg")),
				),
				jen.Return(jen.Id("value"), jen.Err()),
			)
	}
	file.Func().
		Params(jen.Id("a").Id("AutoDeserializeServiceImpl")).
		Id("testCaseBytes").
		Params(jen.Id("endpointName").Qual(serverImport, "EndpointName"), jen.Id("i").Int()).
		Params(jen.Op("[]").Byte()).
		Block(
			jen.Id("cases").Op(":=").Id("a").Dot("TestCases").Index(jen.Id("endpointName")),
			jen.List(jen.Id("posLen"), jen.Id("negLen")).Op(":=").List(
				jen.Len(jen.Id("cases").Dot("Positive")),
				jen.Len(jen.Id("cases").Dot("Negative")),
			),
			jen.If(jen.Id("i").Op("<").Id("posLen")).Block(
				jen.Return(jen.Op("[]").Byte().Call(jen.Id("cases").Dot("Positive").Index(jen.Id("i")))),
			).Else().If(jen.Id("i").Op("<").Id("posLen").Op("+").Id("negLen")).Block(
				jen.Return(jen.Op("[]").Byte().Call(jen.Id("cases").Dot("Negative").Index(jen.Id("i").Op("-").Id("posLen")))),
			),
			jen.Panic(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("invalid test case index %s[%d]"), jen.Id("endpointName"), jen.Id("i"),
			)),
		)

}
