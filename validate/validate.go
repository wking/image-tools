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

// Package validate provides a framework for validating OCI images.
package validate

import (
	"fmt"

	"github.com/mndrix/tap.go"
	"github.com/opencontainers/image-spec/specs-go"
	"github.com/opencontainers/image-tools/image/cas"
	"golang.org/x/net/context"
)

// Validate is a template for validating a CAS object.
type Validator func(ctx context.Context, harness *tap.T, engine cas.Engine, digest string, strict bool)

// Validators is a map from media types to an appropriate Validator function.
var Validators = map[string]Validator{
	"application/vnd.oci.image.layer.tar+gzip": ValidateGzippedLayer,
	"application/vnd.oci.image.layer.nondistributable.tar+gzip": ValidateGzippedLayer,
	// TODO: fill in other types
}

func Validate(ctx context.Context, harness *tap.T, engine cas.Engine, descriptor *specs.Descriptor, strict bool) {
	validator, ok := Validators[descriptor.MediaType]
	harness.Ok(ok, fmt.Sprintf("recognized media type %q", descriptor.MediaType))
	if !ok {
		return
	}
	validator(ctx, harness, engine, descriptor.Digest, strict)
}
