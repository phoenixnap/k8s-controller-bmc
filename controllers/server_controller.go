/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/oauth2/clientcredentials"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bmcv1 "github.com/phoenixnap/k8s-bmc/api/v1"
)

// ServerReconciler reconciles a Server object
type ServerReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

// +kubebuilder:rbac:groups=bmc.api.phoenixnap.com,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bmc.api.phoenixnap.com,resources=servers/status,verbs=get;update;patch

var (
	bmcServerIDAnnotation = `bmc.api.phoenixnap.com/server_id`

	finalizerName = `server.finalizers.bmc.api.phoenixnap.com`

	ENV_BMC_CLIENT_ID     = `BMC_CLIENT_ID`
	ENV_BMC_CLIENT_SECRET = `BMC_CLIENT_SECRET`
	ENV_BMC_TOKEN_URL     = `BMC_TOKEN_URL`
	ENV_BMC_ENDPOINT_URL  = `BMC_ENDPOINT_URL`

	requeueAfter1Min = ctrl.Result{RequeueAfter: 1 * time.Minute}
	requeueAfter2Min = ctrl.Result{RequeueAfter: 2 * time.Minute}
	requeueAfter5Min = ctrl.Result{RequeueAfter: 5 * time.Minute}

	// UpperCamelCase
	EventReasonCleanupError   = `CleanupError`
	EventReasonCleanupSuccess = `CleanupSuccess`

	EventReasonCreated              = `Created`
	EventReasonCreateError          = `CreateError`
	EventReasonCreateErrorPermanent = `CreateErrorPermanent`
	EventReasonCreateErrorInventory = `CreateErrorInventory`
	EventReasonCreateFailure        = `CreateServerFailure`

	EventReasonResourceOrphaned = `ResourceOrphaned`
	EventReasonPollFailure      = `PollingFailure`
	EventReasonStatusChange     = `StatusChange`

	StatusIrreconcilable = `irreconcilable`
	StatusOrphaned       = `orphaned`
	StatusStale          = `stale`
)

func (r *ServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("server", req.NamespacedName)

	// 1. get the Server
	var server bmcv1.Server
	if err := r.Get(ctx, req.NamespacedName, &server); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Load a BMC client
	if len(os.Getenv(ENV_BMC_CLIENT_ID)) <= 0 ||
		len(os.Getenv(ENV_BMC_CLIENT_SECRET)) <= 0 ||
		len(os.Getenv(ENV_BMC_TOKEN_URL)) <= 0 ||
		len(os.Getenv(ENV_BMC_ENDPOINT_URL)) <= 0 {
		log.Error(fmt.Errorf(`incomplete BMC connection configuration`), `incomplete BMC connection configuration`)
		os.Exit(1)
	}
	bmcConfig := clientcredentials.Config{
		ClientID:     os.Getenv(ENV_BMC_CLIENT_ID),
		ClientSecret: os.Getenv(ENV_BMC_CLIENT_SECRET),
		TokenURL:     os.Getenv(ENV_BMC_TOKEN_URL),
		Scopes:       []string{"bmc", "bmc.read"}}
	bmc := bmcConfig.Client(context.Background())

	// 2. Check for delettion activity and finalizer
	if server.ObjectMeta.DeletionTimestamp.IsZero() {
		// Not deleted, verify that our finalizer is present
		found := false
		for _, finalizer := range server.ObjectMeta.Finalizers {
			if finalizer == finalizerName {
				found = true
			}
		}
		if !found {
			// add the finalizer
			log.Info(`attaching finalizer`)
			server.ObjectMeta.Finalizers = append(server.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		log.Info(`finalizing`)

		bmcServerID := server.Annotations[bmcServerIDAnnotation]
		// skip finalization for orphaned resources
		if server.Status.BMCStatus != StatusOrphaned && len(bmcServerID) > 0 {
			// Do BMC cleanup
			apiReq, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%sservers/%s", os.Getenv(ENV_BMC_ENDPOINT_URL), bmcServerID),
				nil)
			if err != nil {
				r.Recorder.Event(&server, `Warning`, EventReasonCleanupError, err.Error())
				return ctrl.Result{}, err
			}
			apiResp, err := bmc.Do(apiReq)
			if err != nil {
				r.Recorder.Event(&server, `Warning`, EventReasonCleanupError, err.Error())
				return ctrl.Result{}, err
			}
			defer apiResp.Body.Close()

			body, err := ioutil.ReadAll(apiResp.Body)
			if err != nil {
				r.Recorder.Event(&server, `Warning`, EventReasonCleanupError, err.Error())
				return ctrl.Result{}, err
			}

			switch apiResp.StatusCode {
			case 400:
				// bad data, or controller/API incompatibility
				log.Info("unable to delete", `code`, 400, `body`, string(body))
				server.Status.BMCStatus = StatusIrreconcilable
				if err := r.Update(ctx, &server); err != nil {
					return ctrl.Result{}, err
				}
				return requeueAfter2Min, nil
			case 401:
				// bad credentials
				log.Info("unable to delete", `code`, 401, `body`, string(body))
				server.Status.BMCStatus = StatusIrreconcilable
				if err := r.Update(ctx, &server); err != nil {
					return ctrl.Result{}, err
				}
				return requeueAfter2Min, nil
			case 403:
				// unauthorized (also 404)
				log.Info("unable to delete", `code`, 403, `body`, string(body))
				server.Status.BMCStatus = StatusOrphaned
				if err := r.Update(ctx, &server); err != nil {
					return ctrl.Result{}, err
				}
				return requeueAfter2Min, nil
			case 500:
				// temporarily unavailable, backoff and retry
				log.Info(`BMC temporarily unavailable`, `body`, string(body))
				return requeueAfter2Min, nil
			case 200:
				fallthrough
			case 201:
				fallthrough
			case 202:
				fallthrough
			case 204:
				// the call was successful, do nothing and continue reconciliation
				r.Recorder.Eventf(&server, `Normal`, EventReasonCleanupSuccess, "Deleted BMC server %s", bmcServerID)
			default:
				r.Recorder.Eventf(&server, `Warning`, EventReasonCleanupError, "Unexpected response from API: %v", apiResp.StatusCode)
				return requeueAfter2Min, fmt.Errorf("unexpected response during server delete: %v", apiResp.StatusCode)
			}
		}

		for i, finalizer := range server.ObjectMeta.Finalizers {
			if finalizer == finalizerName {
				server.ObjectMeta.Finalizers[i] = server.ObjectMeta.Finalizers[len(server.ObjectMeta.Finalizers)-1]
				server.ObjectMeta.Finalizers = server.ObjectMeta.Finalizers[:len(server.ObjectMeta.Finalizers)-1]
				break
			}
		}
		if err := r.Update(ctx, &server); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// 3. Create, poll, or update? Branch on the bmcServerID annotation
	bmcServerID := server.Annotations[bmcServerIDAnnotation]
	if len(bmcServerID) == 0 {
		log.Info(`creating`)
		createBody, err := json.Marshal(server.Spec)
		if err != nil {
			return ctrl.Result{}, err
		}

		apiResp, err := bmc.Post(fmt.Sprintf("%sservers", os.Getenv(ENV_BMC_ENDPOINT_URL)), `application/json`, bytes.NewBuffer(createBody))
		if err != nil {
			r.Recorder.Event(&server, `Warning`, EventReasonCreateError, err.Error())
			return ctrl.Result{}, err
		}
		defer apiResp.Body.Close()

		body, err := ioutil.ReadAll(apiResp.Body)
		if err != nil {
			r.Recorder.Event(&server, `Warning`, EventReasonCreateError, err.Error())
			return ctrl.Result{}, err
		}

		switch apiResp.StatusCode {
		case 400:
			// bad data, or controller/API incompatibility
			r.Recorder.Eventf(&server, `Warning`, EventReasonCreateErrorPermanent, `Code: %v`, apiResp.StatusCode)
			log.Info("unable to reconcile", `code`, 400, `body`, string(body))
			server.Status.BMCStatus = StatusIrreconcilable
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			// something is wrong with the controller or input, stop polling
			return ctrl.Result{}, nil
		case 401:
			// bad credentials
			r.Recorder.Eventf(&server, `Warning`, EventReasonCreateErrorPermanent, `Code: %v`, apiResp.StatusCode)
			log.Info("unable to reconcile", `code`, 401, `body`, string(body))
			server.Status.BMCStatus = StatusIrreconcilable
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			// something is wrong with the controller or input, stop polling
			return ctrl.Result{}, nil
		case 403:
			// unauthorized (also 404)
			r.Recorder.Eventf(&server, `Warning`, EventReasonCreateErrorPermanent, `Code: %v`, apiResp.StatusCode)
			log.Info("unable to reconcile", `code`, 403, `body`, string(body))
			server.Status.BMCStatus = StatusIrreconcilable
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			// something is wrong with the controller or input, stop polling
			return ctrl.Result{}, nil
		case 406:
			// no inventory, backoff and retry
			r.Recorder.Eventf(&server, `Warning`, EventReasonCreateErrorInventory, `Code: %v`, apiResp.StatusCode)
			log.Info("temporary no inventory", `code`, 406, `body`, string(body))
			return requeueAfter5Min, err
		case 409:
			// something is wrong; incompatible state
			r.Recorder.Eventf(&server, `Warning`, EventReasonCreateErrorPermanent, `Code: %v`, apiResp.StatusCode)
			log.Info("unable to reconcile", `code`, 409, `body`, string(body))
			server.Status.BMCStatus = StatusIrreconcilable
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			return requeueAfter2Min, nil
		case 500:
			// temporarily unavailable, backoff and retry
			r.Recorder.Event(&server, `Warning`, EventReasonCreateFailure, `Temporary API failure`)
			log.Info(`BMC temporarily unavailable`, `body`, string(body))
			return requeueAfter2Min, err
		case 200:
			fallthrough
		case 201:
			// the call was successful, do nothing and continue reconciliation
		default:
			r.Recorder.Eventf(&server, `Warning`, EventReasonCreateError, `Unexpected response from API: %v`, apiResp.StatusCode)
			return ctrl.Result{}, fmt.Errorf("unexpected response during server create: %v", apiResp.StatusCode)
		}

		// Set the resulting server ID in the annotation and set status
		var ss bmcv1.ServerStatus
		err = json.Unmarshal(body, &ss)
		if err != nil {
			r.Recorder.Event(&server, `Warning`, EventReasonCreateError, err.Error())
			return ctrl.Result{}, err
		}

		r.Recorder.Eventf(&server, `Normal`, EventReasonCreated, "creatd BMC server %s", ss.BMCServerID)

		server.Status = ss
		server.Annotations[bmcServerIDAnnotation] = ss.BMCServerID
		if err := r.Update(ctx, &server); err != nil {
			return ctrl.Result{}, err
		}
		return requeueAfter1Min, nil

	} else {
		log.Info(`polling`)
		apiResp, err := bmc.Get(fmt.Sprintf("%sservers/%s", os.Getenv(ENV_BMC_ENDPOINT_URL), bmcServerID))
		if err != nil {
			if server.Status.BMCStatus != StatusStale {
				r.Recorder.Eventf(&server, `Normal`, EventReasonStatusChange, `%v -> %v`, server.Status.BMCStatus, StatusStale)
			}
			server.Status.BMCStatus = StatusStale
			if ierr := r.Update(ctx, &server); ierr != nil {
				return ctrl.Result{}, ierr
			}
			return requeueAfter2Min, err
		}
		defer apiResp.Body.Close()

		body, err := ioutil.ReadAll(apiResp.Body)
		if err != nil {
			if server.Status.BMCStatus != StatusStale {
				r.Recorder.Eventf(&server, `Normal`, EventReasonStatusChange, `%v -> %v`, server.Status.BMCStatus, StatusStale)
			}
			server.Status.BMCStatus = StatusStale
			if ierr := r.Update(ctx, &server); ierr != nil {
				return ctrl.Result{}, ierr
			}
			return requeueAfter2Min, err
		}

		switch apiResp.StatusCode {
		case 400:
			// bad data, or controller/API incompatibility
			r.Recorder.Eventf(&server, `Warning`, EventReasonPollFailure, `Code: %v`, apiResp.StatusCode)
			log.Info("unable to reconcile", `code`, 400, `body`, string(body))
			server.Status.BMCStatus = StatusIrreconcilable
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			return requeueAfter5Min, nil
		case 401:
			// bad credentials
			r.Recorder.Eventf(&server, `Warning`, EventReasonPollFailure, `Code: %v`, apiResp.StatusCode)
			log.Info("unable to reconcile", `code`, 401, `body`, string(body))
			server.Status.BMCStatus = StatusIrreconcilable
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			return requeueAfter5Min, nil
		case 403:
			// unauthorized (also 404)
			r.Recorder.Event(&server, `Warning`, EventReasonResourceOrphaned, `Access to BMC resource was denied`)
			log.Info("unable to reconcile", `code`, 403, `body`, string(body))
			server.Status.BMCStatus = StatusOrphaned
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		case 500:
			// temporarily unavailable, backoff and retry
			r.Recorder.Eventf(&server, `Warning`, EventReasonPollFailure, `Code: %v`, apiResp.StatusCode)
			log.Info(`BMC temporarily unavailable`, `body`, string(body))
			server.Status.BMCStatus = StatusStale
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
			return requeueAfter5Min, nil
		case 200:
		default:
			r.Recorder.Eventf(&server, `Warning`, EventReasonPollFailure, `Unexpected response from API: %v`, apiResp.StatusCode)
			return ctrl.Result{}, fmt.Errorf("unexpected response during server poll: %v", apiResp.StatusCode)
		}

		// Update the status
		var ss bmcv1.ServerStatus
		err = json.Unmarshal(body, &ss)
		if err != nil {
			return requeueAfter2Min, err
		}

		// detect a status delta
		if server.Status.BMCStatus != ss.BMCStatus {
			r.Recorder.Eventf(&server, `Normal`, EventReasonStatusChange, `%v -> %v`, server.Status.BMCStatus, ss.BMCStatus)
		}

		server.Status = ss
		if err := r.Update(ctx, &server); err != nil {
			return requeueAfter2Min, err
		}

		// BMC server details are mostly immutable. However servers do have
		// power state and can have SSH and other OS configuration, "reset."
		// A ValidatingWebhook should prevent users or other controllers from
		// changing a server resource.

		// Poll timing based on status and expected change
		switch ss.BMCStatus {
		case `powered-on`:
			return requeueAfter2Min, nil
		default:
			return requeueAfter1Min, nil
		}
	}
}

func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bmcv1.Server{}).
		Complete(r)
}
