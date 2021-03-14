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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var serverlog = logf.Log.WithName("server-resource")

func (r *Server) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-bmc-api-phoenixnap-com-v1-server,mutating=true,failurePolicy=fail,groups=bmc.api.phoenixnap.com,resources=servers,verbs=create;update,versions=v1,name=mserver.kb.io

var _ webhook.Defaulter = &Server{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Server) Default() {
	if r.Spec.OS == `` {
		r.Spec.OS = UbuntuBionic
	}
	if r.Spec.NetworkType == `` {
		r.Spec.NetworkType = PublicAndPrivate
	}
	if r.Spec.Location == `` {
		r.Spec.Location = Phoenix
	}
	if r.Spec.Type == `` {
		r.Spec.Type = S1C1Small
	}
	if r.Spec.InstallDefaultSSHKeys == nil {
		r.Spec.InstallDefaultSSHKeys = new(bool)
		*r.Spec.InstallDefaultSSHKeys = true
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-bmc-api-phoenixnap-com-v1-server,mutating=false,failurePolicy=fail,groups=bmc.api.phoenixnap.com,resources=servers,versions=v1,name=vserver.kb.io

var _ webhook.Validator = &Server{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Server) ValidateCreate() error {
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Server) ValidateUpdate(old runtime.Object) error {
	prev := old.(*Server)
	serverlog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	if r.Spec.Hostname != prev.Spec.Hostname {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`hostname`), `immutable`))
	}
	if r.Spec.Description != prev.Spec.Description {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`description`), `immutable`))
	}
	if r.Spec.OS != prev.Spec.OS {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`os`), `immutable`))
	}
	if r.Spec.Type != prev.Spec.Type {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`type`), `immutable`))
	}
	if r.Spec.Location != prev.Spec.Location {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`location`), `immutable`))
	}
	if r.Spec.NetworkType != prev.Spec.NetworkType {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`networkType`), `immutable`))
	}
	if len(r.Spec.SSHKeyIDs) != len(prev.Spec.SSHKeyIDs) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`sshKeyIds`), `immutable`))
	} else {
		for i, _ := range r.Spec.SSHKeyIDs {
			if r.Spec.SSHKeyIDs[i] != prev.Spec.SSHKeyIDs[i] {
				allErrs = append(allErrs, field.Forbidden(field.NewPath(`spec`).Child(`sshKeyIds`), `immutable`))
			}
		}
	}
	if len(allErrs) <= 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: `bmc.api.phoenixnap.com`, Kind: `Server`}, r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Server) ValidateDelete() error {
	return nil
}
