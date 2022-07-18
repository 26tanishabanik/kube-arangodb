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

package deployment

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	"github.com/arangodb/arangosync-client/client"
	agencydriver "github.com/arangodb/go-driver/agency"

	api "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v1"
	"github.com/arangodb/kube-arangodb/pkg/deployment/acs"
	"github.com/arangodb/kube-arangodb/pkg/deployment/acs/sutil"
	"github.com/arangodb/kube-arangodb/pkg/deployment/agency"
	"github.com/arangodb/kube-arangodb/pkg/deployment/chaos"
	deploymentClient "github.com/arangodb/kube-arangodb/pkg/deployment/client"
	memberState "github.com/arangodb/kube-arangodb/pkg/deployment/member"
	"github.com/arangodb/kube-arangodb/pkg/deployment/patch"
	"github.com/arangodb/kube-arangodb/pkg/deployment/reconcile"
	"github.com/arangodb/kube-arangodb/pkg/deployment/reconciler"
	"github.com/arangodb/kube-arangodb/pkg/deployment/resilience"
	"github.com/arangodb/kube-arangodb/pkg/deployment/resources"
	"github.com/arangodb/kube-arangodb/pkg/deployment/resources/inspector"
	"github.com/arangodb/kube-arangodb/pkg/logging"
	"github.com/arangodb/kube-arangodb/pkg/operator/scope"
	"github.com/arangodb/kube-arangodb/pkg/util"
	"github.com/arangodb/kube-arangodb/pkg/util/arangod"
	"github.com/arangodb/kube-arangodb/pkg/util/arangod/conn"
	"github.com/arangodb/kube-arangodb/pkg/util/errors"
	"github.com/arangodb/kube-arangodb/pkg/util/globals"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil"
	inspectorInterface "github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector"
	"github.com/arangodb/kube-arangodb/pkg/util/kclient"
	"github.com/arangodb/kube-arangodb/pkg/util/trigger"
)

// Config holds configuration settings for a Deployment
type Config struct {
	ServiceAccount            string
	AllowChaos                bool
	ScalingIntegrationEnabled bool
	OperatorImage             string
	ArangoImage               string
	Scope                     scope.Scope
}

// Dependencies holds dependent services for a Deployment
type Dependencies struct {
	EventRecorder record.EventRecorder

	Client kclient.Client
}

// deploymentEventType strongly typed type of event
type deploymentEventType string

const (
	eventArangoDeploymentUpdated deploymentEventType = "ArangoDeploymentUpdated"
)

// deploymentEvent holds an event passed from the controller to the deployment.
type deploymentEvent struct {
	Type       deploymentEventType
	Deployment *api.ArangoDeployment
}

const (
	deploymentEventQueueSize = 256
	minInspectionInterval    = 250 * util.Interval(time.Millisecond) // Ensure we inspect the generated resources no less than with this interval
	maxInspectionInterval    = 10 * util.Interval(time.Second)       // Ensure we inspect the generated resources no less than with this interval
)

type deploymentStatusObject struct {
	version int32
	last    api.DeploymentStatus // Internal status copy of the CR
}

// Deployment is the in process state of an ArangoDeployment.
type Deployment struct {
	log logging.Logger

	name      string
	namespace string

	apiObject *api.ArangoDeployment // API object
	status    struct {
		mutex sync.Mutex
		deploymentStatusObject
	}
	config Config
	deps   Dependencies

	eventCh chan *deploymentEvent
	stopCh  chan struct{}
	stopped int32

	inspectTrigger            trigger.Trigger
	inspectCRDTrigger         trigger.Trigger
	updateDeploymentTrigger   trigger.Trigger
	clientCache               deploymentClient.Cache
	agencyCache               agency.Cache
	recentInspectionErrors    int
	clusterScalingIntegration *clusterScalingIntegration
	reconciler                *reconcile.Reconciler
	resilience                *resilience.Resilience
	resources                 *resources.Resources
	chaosMonkey               *chaos.Monkey
	acs                       sutil.ACS
	syncClientCache           client.ClientCache
	haveServiceMonitorCRD     bool

	memberState memberState.StateInspector

	metrics Metrics
}

func (d *Deployment) WithArangoMember(cache inspectorInterface.Inspector, timeout time.Duration, name string) reconciler.ArangoMemberModContext {
	return reconciler.NewArangoMemberModContext(cache, timeout, name)
}

func (d *Deployment) WithCurrentArangoMember(name string) reconciler.ArangoMemberModContext {
	return d.WithArangoMember(d.acs.CurrentClusterCache(), globals.GetGlobals().Timeouts().Kubernetes().Get(), name)
}

func (d *Deployment) GetMembersState() memberState.StateInspector {
	return d.memberState
}

func (d *Deployment) GetAgencyCache() (agency.State, bool) {
	return d.agencyCache.Data()
}

func (d *Deployment) GetAgencyHealth() (agency.Health, bool) {
	return d.agencyCache.Health()
}

func (d *Deployment) RefreshAgencyCache(ctx context.Context) (uint64, error) {
	if d.apiObject.Spec.Mode.Get() == api.DeploymentModeSingle {
		return 0, nil
	}

	if info := d.apiObject.Status.Agency; info != nil {
		if size := info.Size; size != nil {
			lCtx, c := globals.GetGlobalTimeouts().Agency().WithTimeout(ctx)
			defer c()

			rsize := int(*size)

			clients := make(map[string]agencydriver.Agency)
			for _, m := range d.GetStatusSnapshot().Members.Agents {
				a, err := d.GetAgency(lCtx, m.ID)
				if err != nil {
					return 0, err
				}

				clients[m.ID] = a
			}

			return d.agencyCache.Reload(lCtx, rsize, clients)
		}
	}

	return 0, errors.Newf("Agency not yet established")
}

func (d *Deployment) SetAgencyMaintenanceMode(ctx context.Context, enabled bool) error {
	if !d.GetMode().HasAgents() {
		return nil
	}

	ctxChild, cancel := globals.GetGlobalTimeouts().ArangoD().WithTimeout(ctx)
	defer cancel()
	client, err := d.GetDatabaseClient(ctxChild)
	if err != nil {
		return err
	}

	data := "on"
	if !enabled {
		data = "off"
	}

	conn := client.Connection()
	r, err := conn.NewRequest(http.MethodPut, "/_admin/cluster/maintenance")
	if err != nil {
		return err
	}

	if _, err := r.SetBody(data); err != nil {
		return err
	}

	resp, err := conn.Do(ctx, r)
	if err != nil {
		return err
	}

	if err := resp.CheckStatus(http.StatusOK); err != nil {
		return err
	}

	return nil
}

// New creates a new Deployment from the given API object.
func New(config Config, deps Dependencies, apiObject *api.ArangoDeployment) (*Deployment, error) {
	if err := apiObject.Spec.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	i := inspector.NewInspector(inspector.NewDefaultThrottle(), deps.Client, apiObject.GetNamespace(), apiObject.GetName())

	d := &Deployment{
		apiObject:   apiObject,
		name:        apiObject.GetName(),
		namespace:   apiObject.GetNamespace(),
		config:      config,
		deps:        deps,
		eventCh:     make(chan *deploymentEvent, deploymentEventQueueSize),
		stopCh:      make(chan struct{}),
		agencyCache: agency.NewCache(apiObject.GetNamespace(), apiObject.GetName(), apiObject.Spec.Mode),
		acs:         acs.NewACS(apiObject.GetUID(), i),
	}

	d.log = logger.WrapObj(d)

	d.memberState = memberState.NewStateInspector(d)

	d.clientCache = deploymentClient.NewClientCache(d, conn.NewFactory(d.getAuth, d.getConnConfig))

	d.status.last = *(apiObject.Status.DeepCopy())
	d.reconciler = reconcile.NewReconciler(apiObject.GetNamespace(), apiObject.GetName(), d)
	d.resilience = resilience.NewResilience(apiObject.GetNamespace(), apiObject.GetName(), d)
	d.resources = resources.NewResources(apiObject.GetNamespace(), apiObject.GetName(), d)
	if d.status.last.AcceptedSpec == nil {
		// We've validated the spec, so let's use it from now.
		d.status.last.AcceptedSpec = apiObject.Spec.DeepCopy()
	}

	localInventory.Add(d)

	go d.run()
	go d.listenForPodEvents(d.stopCh)
	go d.listenForPVCEvents(d.stopCh)
	go d.listenForSecretEvents(d.stopCh)
	go d.listenForServiceEvents(d.stopCh)
	go d.listenForCRDEvents(d.stopCh)
	if apiObject.Spec.GetMode() == api.DeploymentModeCluster {
		ci := newClusterScalingIntegration(d)
		d.clusterScalingIntegration = ci
		go ci.ListenForClusterEvents(d.stopCh)
	}
	if config.AllowChaos {
		d.chaosMonkey = chaos.NewMonkey(apiObject.GetNamespace(), apiObject.GetName(), d)
		go d.chaosMonkey.Run(d.stopCh)
	}

	return d, nil
}

// Update the deployment.
// This sends an update event in the deployment event queue.
func (d *Deployment) Update(apiObject *api.ArangoDeployment) {
	d.send(&deploymentEvent{
		Type:       eventArangoDeploymentUpdated,
		Deployment: apiObject,
	})
}

// Stop the deployment.
// Called when the deployment was deleted by the user.
func (d *Deployment) Stop() {
	d.log.Info("deployment is deleted by user")
	if atomic.CompareAndSwapInt32(&d.stopped, 0, 1) {
		close(d.stopCh)
	}
}

// send given event into the deployment event queue.
func (d *Deployment) send(ev *deploymentEvent) {
	select {
	case d.eventCh <- ev:
		l, ecap := len(d.eventCh), cap(d.eventCh)
		if l > int(float64(ecap)*0.8) {
			d.log.
				Int("used", l).
				Int("capacity", ecap).
				Warn("event queue buffer is almost full")
		}
	case <-d.stopCh:
	}
}

// run is the core the core worker.
// It processes the event queue and polls the state of generated
// resource on a regular basis.
func (d *Deployment) run() {
	log := d.log

	// Create agency mapping
	if err := d.createAgencyMapping(context.TODO()); err != nil {
		d.CreateEvent(k8sutil.NewErrorEvent("Failed to create agency mapping members", err, d.GetAPIObject()))
	}

	if d.GetPhase() == api.DeploymentPhaseNone {
		// Create service monitor
		if d.haveServiceMonitorCRD {
			if err := d.resources.EnsureServiceMonitor(context.TODO()); err != nil {
				d.CreateEvent(k8sutil.NewErrorEvent("Failed to create service monitor", err, d.GetAPIObject()))
			}
		}

		// Create initial topology
		if err := d.createInitialTopology(context.TODO()); err != nil {
			d.CreateEvent(k8sutil.NewErrorEvent("Failed to create initial topology", err, d.GetAPIObject()))
		}

		status, lastVersion := d.GetStatus()
		status.Phase = api.DeploymentPhaseRunning
		if err := d.UpdateStatus(context.TODO(), status, lastVersion); err != nil {
			log.Err(err).Warn("update initial CR status failed")
		}
		log.Info("start running...")
	}

	d.lookForServiceMonitorCRD()

	// Execute inspection for first time without delay of 10s
	log.Debug("Initially inspect deployment...")
	inspectionInterval := d.inspectDeployment(minInspectionInterval)
	log.Str("interval", inspectionInterval.String()).Debug("...deployment inspect started")

	for {
		select {
		case <-d.stopCh:
			// We're being stopped.
			return

		case event := <-d.eventCh:
			// Got event from event queue
			switch event.Type {
			case eventArangoDeploymentUpdated:
				d.updateDeploymentTrigger.Trigger()
			default:
				panic("unknown event type" + event.Type)
			}

		case <-d.inspectTrigger.Done():
			log.Trace("Inspect deployment...")
			inspectionInterval = d.inspectDeployment(inspectionInterval)
			log.Str("interval", inspectionInterval.String()).Trace("...inspected deployment")

		case <-d.inspectCRDTrigger.Done():
			d.lookForServiceMonitorCRD()
		case <-d.updateDeploymentTrigger.Done():
			inspectionInterval = minInspectionInterval
			if err := d.handleArangoDeploymentUpdatedEvent(context.TODO()); err != nil {
				d.CreateEvent(k8sutil.NewErrorEvent("Failed to handle deployment update", err, d.GetAPIObject()))
			}

		case <-inspectionInterval.After():
			// Trigger inspection
			d.inspectTrigger.Trigger()
			// Backoff with next interval
			inspectionInterval = inspectionInterval.Backoff(1.5, maxInspectionInterval)
		}
	}
}

// handleArangoDeploymentUpdatedEvent is called when the deployment is updated by the user.
func (d *Deployment) handleArangoDeploymentUpdatedEvent(ctx context.Context) error {
	log := d.log.Str("deployment", d.apiObject.GetName())

	// Get the most recent version of the deployment from the API server
	ctxChild, cancel := globals.GetGlobalTimeouts().Kubernetes().WithTimeout(ctx)
	defer cancel()

	current, err := d.deps.Client.Arango().DatabaseV1().ArangoDeployments(d.apiObject.GetNamespace()).Get(ctxChild, d.apiObject.GetName(), meta.GetOptions{})
	if err != nil {
		log.Err(err).Debug("Failed to get current version of deployment from API server")
		if k8sutil.IsNotFound(err) {
			return nil
		}
		return errors.WithStack(err)
	}

	specBefore := d.apiObject.Spec
	status := d.status.last
	if d.status.last.AcceptedSpec != nil {
		specBefore = *status.AcceptedSpec.DeepCopy()
	}
	newAPIObject := current.DeepCopy()
	newAPIObject.Spec.SetDefaultsFrom(specBefore)
	newAPIObject.Spec.SetDefaults(d.apiObject.GetName())

	resetFields := specBefore.ResetImmutableFields(&newAPIObject.Spec)
	if len(resetFields) > 0 {
		log.Strs("fields", resetFields...).Debug("Found modified immutable fields")
		newAPIObject.Spec.SetDefaults(d.apiObject.GetName())
	}
	if err := newAPIObject.Spec.Validate(); err != nil {
		d.CreateEvent(k8sutil.NewErrorEvent("Validation failed", err, d.apiObject))
		// Try to reset object
		if err := d.updateCRSpec(ctx, d.apiObject.Spec); err != nil {
			log.Err(err).Error("Restore original spec failed")
			d.CreateEvent(k8sutil.NewErrorEvent("Restore original failed", err, d.apiObject))
		}
		return nil
	}
	if len(resetFields) > 0 {
		for _, fieldName := range resetFields {
			log.Str("field", fieldName).Debug("Reset immutable field")
			d.CreateEvent(k8sutil.NewImmutableFieldEvent(fieldName, d.apiObject))
		}
	}

	// Save updated spec
	if err := d.updateCRSpec(ctx, newAPIObject.Spec); err != nil {
		return errors.WithStack(errors.Newf("failed to update ArangoDeployment spec: %v", err))
	}
	// Save updated accepted spec
	{
		status, lastVersion := d.GetStatus()
		if newAPIObject.Status.IsForceReload() {
			log.Warn("Forced status reload!")
			status = newAPIObject.Status
			status.ForceStatusReload = nil
		}
		status.AcceptedSpec = newAPIObject.Spec.DeepCopy()
		if err := d.UpdateStatus(ctx, status, lastVersion); err != nil {
			return errors.WithStack(errors.Newf("failed to update ArangoDeployment status: %v", err))
		}
	}

	// Notify cluster of desired server count
	if ci := d.clusterScalingIntegration; ci != nil {
		ci.SendUpdateToCluster(d.apiObject.Spec)
	}

	// Trigger inspect
	d.inspectTrigger.Trigger()

	return nil
}

// CreateEvent creates a given event.
// On error, the error is logged.
func (d *Deployment) CreateEvent(evt *k8sutil.Event) {
	d.deps.EventRecorder.Event(evt.InvolvedObject, evt.Type, evt.Reason, evt.Message)
}

// Update the status of the API object from the internal status
func (d *Deployment) updateCRStatus(ctx context.Context, force ...bool) error {
	if len(force) == 0 || !force[0] {
		if d.apiObject.Status.Equal(d.status.last) {
			// Nothing has changed
			return nil
		}
	}

	// Send update to API server
	depls := d.deps.Client.Arango().DatabaseV1().ArangoDeployments(d.GetNamespace())
	attempt := 0
	for {
		attempt++
		q := patch.NewPatch(patch.ItemReplace(patch.NewPath("status"), d.status.last))

		if d.apiObject.GetDeletionTimestamp() == nil {
			if ensureFinalizers(d.apiObject) {
				q.ItemAdd(patch.NewPath("metadata", "finalizers"), d.apiObject.Finalizers)
			}
		}

		var newAPIObject *api.ArangoDeployment
		err := globals.GetGlobalTimeouts().Kubernetes().RunWithTimeout(ctx, func(ctxChild context.Context) error {
			p, err := q.Marshal()
			if err != nil {
				return err
			}

			newAPIObject, err = depls.Patch(ctxChild, d.GetName(), types.JSONPatchType, p, meta.PatchOptions{})

			return err
		})
		if err == nil {
			// Update internal object
			d.apiObject = newAPIObject
			return nil
		}
		if attempt < 10 {
			continue
		}
		if err != nil {
			d.log.Err(err).Debug("failed to patch ArangoDeployment status")
			return errors.WithStack(errors.Newf("failed to patch ArangoDeployment status: %v", err))
		}
	}
}

// Update the spec part of the API object (d.apiObject)
// to the given object, while preserving the status.
// On success, d.apiObject is updated.
func (d *Deployment) updateCRSpec(ctx context.Context, newSpec api.DeploymentSpec, force ...bool) error {

	if len(force) == 0 || !force[0] {
		if d.apiObject.Spec.Equal(&newSpec) {
			d.log.Trace("Nothing to update in updateCRSpec")
			// Nothing to update
			return nil
		}
	}

	// Send update to API server
	update := d.apiObject.DeepCopy()
	attempt := 0
	for {
		attempt++
		update.Spec = newSpec
		update.Status = d.status.last
		ns := d.apiObject.GetNamespace()
		var newAPIObject *api.ArangoDeployment
		err := globals.GetGlobalTimeouts().Kubernetes().RunWithTimeout(ctx, func(ctxChild context.Context) error {
			var err error
			newAPIObject, err = d.deps.Client.Arango().DatabaseV1().ArangoDeployments(ns).Update(ctxChild, update, meta.UpdateOptions{})

			return err
		})
		if err == nil {
			// Update internal object
			d.apiObject = newAPIObject
			return nil
		}
		if attempt < 10 && k8sutil.IsConflict(err) {
			// API object may have been changed already,
			// Reload api object and try again
			var current *api.ArangoDeployment

			err = globals.GetGlobalTimeouts().Kubernetes().RunWithTimeout(ctx, func(ctxChild context.Context) error {
				var err error
				current, err = d.deps.Client.Arango().DatabaseV1().ArangoDeployments(ns).Get(ctxChild, update.GetName(), meta.GetOptions{})

				return err
			})
			if err == nil {
				update = current.DeepCopy()
				continue
			}
		}
		if err != nil {
			d.log.Err(err).Debug("failed to patch ArangoDeployment spec")
			return errors.WithStack(errors.Newf("failed to patch ArangoDeployment spec: %v", err))
		}
	}
}

// isOwnerOf returns true if the given object belong to this deployment.
func (d *Deployment) isOwnerOf(obj meta.Object) bool {
	ownerRefs := obj.GetOwnerReferences()
	if len(ownerRefs) < 1 {
		return false
	}
	return ownerRefs[0].UID == d.apiObject.UID
}

// lookForServiceMonitorCRD checks if there is a CRD for the ServiceMonitor
// CR and sets the flag haveServiceMonitorCRD accordingly. This is called
// once at creation time of the deployment and then always if the CRD
// informer is triggered.
func (d *Deployment) lookForServiceMonitorCRD() {
	var err error
	if d.GetScope().IsNamespaced() {
		_, err = d.acs.CurrentClusterCache().ServiceMonitor().V1()
		if k8sutil.IsForbiddenOrNotFound(err) {
			return
		}
	} else {
		_, err = d.deps.Client.KubernetesExtensions().ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), "servicemonitors.monitoring.coreos.com", meta.GetOptions{})
	}
	log := d.log
	log.Debug("Looking for ServiceMonitor CRD...")
	if err == nil {
		if !d.haveServiceMonitorCRD {
			log.Info("...have discovered ServiceMonitor CRD")
		}
		d.haveServiceMonitorCRD = true
		d.triggerInspection()
		return
	} else if k8sutil.IsNotFound(err) {
		if d.haveServiceMonitorCRD {
			log.Info("...ServiceMonitor CRD no longer there")
		}
		d.haveServiceMonitorCRD = false
		return
	}
	log.Err(err).Warn("Error when looking for ServiceMonitor CRD")
}

// SetNumberOfServers adjust number of DBservers and coordinators in arangod
func (d *Deployment) SetNumberOfServers(ctx context.Context, noCoordinators, noDBServers *int) error {
	ctxChild, cancel := globals.GetGlobalTimeouts().ArangoD().WithTimeout(ctx)
	defer cancel()
	c, err := d.clientCache.GetDatabase(ctxChild)
	if err != nil {
		return errors.WithStack(err)
	}

	err = globals.GetGlobalTimeouts().ArangoD().RunWithTimeout(ctx, func(ctxChild context.Context) error {
		return arangod.SetNumberOfServers(ctxChild, c.Connection(), noCoordinators, noDBServers)
	})

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Deployment) ApplyPatch(ctx context.Context, p ...patch.Item) error {
	parser := patch.Patch(p)

	data, err := parser.Marshal()
	if err != nil {
		return err
	}

	c := d.deps.Client.Arango().DatabaseV1().ArangoDeployments(d.apiObject.GetNamespace())

	ctxChild, cancel := globals.GetGlobalTimeouts().Kubernetes().WithTimeout(ctx)
	defer cancel()
	depl, err := c.Patch(ctxChild, d.apiObject.GetName(), types.JSONPatchType, data, meta.PatchOptions{})
	if err != nil {
		return err
	}

	d.apiObject = depl

	return nil
}
