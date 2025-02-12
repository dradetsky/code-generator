// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

func NewModelForService(t *testing.T, serviceAlias string) *ackmodel.Model {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// We have subdirectories in pkg/generate and pkg/model that rely on the testdata
	// in pkg/generate. This code simply detects if we're running from one of
	// those subdirectories and if so, rebuilds the path to the API model files
	// in pkg/generate/testdata
	pathParts := strings.Split(path, "/")
	for x, pathPart := range pathParts {
		if pathPart == "generate" || pathPart == "model" {
			path = filepath.Join(pathParts[0:x]...)
			path = filepath.Join("/", path, "testdata")
			break
		}
	}
	sdkHelper := model.NewSDKHelper(path)
	sdkAPI, err := sdkHelper.API(serviceAlias)
	if err != nil {
		t.Fatal(err)
	}
	generatorConfigPath := filepath.Join(path, "models", "apis", serviceAlias, "0000-00-00", "generator.yaml")
	if _, err := os.Stat(generatorConfigPath); os.IsNotExist(err) {
		generatorConfigPath = ""
	}
	g, err := ackmodel.New(sdkAPI, "v1alpha1", generatorConfigPath, ackgenerate.DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	return g
}
