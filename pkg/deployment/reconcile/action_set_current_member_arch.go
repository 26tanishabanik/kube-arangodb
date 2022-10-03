//
// DISCLAIMER
//
// Copyright 2016-2022 ArangoDB GmbH, Cologne, Germany
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

package reconcile

import (
	"context"

	api "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v1"
	"github.com/arangodb/kube-arangodb/pkg/util/errors"
)

// newSetCurrentMemberArchAction creates a new Action that implements the given planned SetMemberCurrentArch action.
func newSetCurrentMemberArchAction(action api.Action, actionCtx ActionContext) Action {
	a := &actionSetCurrentMemberArch{}

	a.actionImpl = newActionImplDefRef(action, actionCtx)

	return a
}

// actionSetCurrentMemberArch implements an SetMemberCurrentArch.
type actionSetCurrentMemberArch struct {
	// actionImpl implement timeout and member id functions
	actionImpl
}

// Start performs the start of the action.
// Returns true if the action is completely finished, false in case
// the start time needs to be recorded and a ready condition needs to be checked.
func (a *actionSetCurrentMemberArch) Start(ctx context.Context) (bool, error) {
	ready, _, err := a.CheckProgress(ctx)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return ready, nil
}

// CheckProgress checks the progress of the action.
// Returns true if the action is completely finished, false otherwise.
func (a *actionSetCurrentMemberArch) CheckProgress(ctx context.Context) (bool, bool, error) {
	arch := a.action.Architecture

	if err := a.actionCtx.WithStatusUpdate(ctx, func(s *api.DeploymentStatus) bool {
		m, g, found := s.Members.ElementByID(a.action.MemberID)
		if !found {
			a.log.Error("No such member")
			return false
		}

		if !m.Architecture.Equal(&arch) {
			m.Architecture = &arch
		}

		if err := s.Members.Update(m, g); err != nil {
			a.log.Error("Member update failed")
			return false
		}

		return true
	}); err != nil {
		a.log.Error("Member failed")
		return true, false, nil
	}

	return true, false, nil
}
