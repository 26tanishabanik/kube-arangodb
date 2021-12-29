//
// DISCLAIMER
//
// Copyright 2016-2021 ArangoDB GmbH, Cologne, Germany
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

package k8sutil

import (
	"github.com/arangodb/kube-arangodb/pkg/util"
	"github.com/arangodb/kube-arangodb/pkg/util/constants"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/secret"
)

type License string

func (l License) IsV2Set() bool {
	return l != ""
}

func (l License) V2Hash() string {
	return util.SHA256FromString(string(l))
}

type LicenseSecret struct {
	V1 string
	V2 License
}

func GetLicenseFromSecret(secret secret.Inspector, name string) (LicenseSecret, bool) {
	s, ok := secret.Secret(name)
	if !ok {
		return LicenseSecret{}, false
	}

	var l LicenseSecret

	if v, ok := s.Data[constants.SecretKeyToken]; ok {
		l.V1 = string(v)
	}

	if v1, ok1 := s.Data[constants.SecretKeyV2License]; ok1 {
		l.V2 = License(v1)
	} else if v2, ok2 := s.Data[constants.SecretKeyV2Token]; ok2 {
		l.V2 = License(v2)
	}

	return l, true
}