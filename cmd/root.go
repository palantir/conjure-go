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

package cmd

import (
	"strings"

	"github.com/palantir/pkg/cobracli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/conjure-go/conjure"
)

const (
	outputDirFlagName = "output"
	serverFlagName    = "server"
)

var (
	version          = "unspecified"
	debug            bool
	outputDirFlagVar string
	serverFlagVar    bool
)

var rootCmd = &cobra.Command{
	Use:   "conjure",
	Short: "Generates Go files based on a specified input Conjure IR file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Generate(args[0], outputDirFlagVar)
	},
}

func Execute() int {
	return cobracli.ExecuteWithDebugVarAndDefaultParams(rootCmd, &debug, cobracli.VersionFlagParam(version))
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debug output")
	rootCmd.Flags().StringVar(&outputDirFlagVar, outputDirFlagName, ".", "base directory into which generated Conjure is written")
	rootCmd.Flags().BoolVar(&serverFlagVar, serverFlagName, false, "enable witchcraft-go server generation")
}

func Generate(irFile, outDir string) error {
	if !strings.HasSuffix(irFile, ".json") {
		return errors.Errorf(`IR file %s does not have suffix ".json"`, irFile)
	}
	conjureDefinition, err := conjure.FromIRFile(irFile)
	if err != nil {
		return err
	}
	output := conjure.OutputConfiguration{OutputDir: outDir}
	if serverFlagVar {
		output.ServerType = conjure.WitchcraftServer
	}
	if err := conjure.Generate(conjureDefinition, output); err != nil {
		return errors.Wrapf(err, "failed to generate Conjure")
	}
	return nil
}
