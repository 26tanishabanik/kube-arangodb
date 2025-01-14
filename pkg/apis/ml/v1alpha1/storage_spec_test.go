//
// DISCLAIMER
//
// Copyright 2023 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//

package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
)

func Test_ArangoMLStorageSpec(t *testing.T) {
	s := ArangoMLStorageSpec{
		ListenPort: nil,
		Resources:  core.ResourceRequirements{},
		S3:         nil,
	}
	s.SetDefaults()
	require.Error(t, s.Validate())

	s.S3 = &ArangoMLStorageS3Spec{
		Endpoint:              "some-endpoint",
		DisableSSL:            false,
		Region:                "",
		BucketName:            "test-bucket",
		CredentialsSecretName: "some-secret",
	}
	require.NoError(t, s.Validate())
}
