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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/mndrix/tap.go"
	"github.com/opencontainers/image-tools/image/cas"
	"golang.org/x/net/context"
)

func ValidateGzippedLayer(ctx context.Context, harness *tap.T, engine cas.Engine, digest string, strict bool) {
	reader, err := engine.Get(ctx, digest)
	harness.Ok(err == nil, fmt.Sprintf("retrieve %s from CAS", digest))
	if err != nil {
		return
	}

	gzipReader, err := gzip.NewReader(reader)
	harness.Ok(err == nil, fmt.Sprintf("gzip reader for %s", digest))
	if err != nil {
		return
	}

	tarReader := tar.NewReader(gzipReader)
	for {
		select {
		case <-ctx.Done():
			harness.Ok(false, ctx.Err().Error())
			return
		default:
		}

		header, err := tarReader.Next()
		if err == io.EOF {
			return
		}
		harness.Ok(err == nil, fmt.Sprintf("read tar header from %s", digest))
		if err != nil {
			return
		}

		message := fmt.Sprintf("ustar typeflag in %s (%q)", digest, header.Typeflag)
		if header.Typeflag < tar.TypeReg || header.Typeflag > tar.TypeFifo {
			harness.Ok(!strict, message)
		} else {
			harness.Ok(true, message)
		}
	}

	return
}
