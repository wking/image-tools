// Copyright 2016 The Linux Foundation
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
	"fmt"
	"os"

	"github.com/mndrix/tap.go"
	"github.com/opencontainers/image-spec/specs-go"
	"github.com/opencontainers/image-tools/image/cas/layout"
	"github.com/opencontainers/image-tools/validate"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type validateCmd struct {
	mediaType string // the type to validate, can be empty string
	digest string
	strict bool
}

func main() {
	validator := &validateCmd{}

	cmd := &cobra.Command{
		Use:   "oci-image-validate PATH DIGEST",
		Short: "Validate an OCI image",
		Run:   validator.Run,
	}

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func (validator *validateCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		err := cmd.Usage()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	path := args[0]
	validator.mediaType = "application/vnd.oci.image.layer.tar+gzip"
	validator.digest = args[1]
	validator.strict = true

	ctx := context.Background()

	engine, err := layout.NewEngine(ctx, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer engine.Close()

	harness := tap.New()
	harness.Header(0)

	descriptor := &specs.Descriptor{
		MediaType: validator.mediaType,
		Digest: validator.digest,
		Size: -1,
	}
	validate.Validate(ctx, harness, engine, descriptor, validator.strict)

	harness.AutoPlan()

	return
}
