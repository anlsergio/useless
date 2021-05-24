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
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	uselessv1 "useless/api/v1"
)

// MachineReconciler reconciles a Machine object
type MachineReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=useless.my.domain,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=useless.my.domain,resources=machines/status,verbs=get;update;patch

func (r *MachineReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("machine", req.NamespacedName)

	// your logic here
	var machine uselessv1.Machine

	if err := r.Get(ctx, req.NamespacedName, &machine); err != nil {
		log.Info("Something went wrong while trying to get the resource", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if machine.Status.Status == "" {
		machine.Status.Status = "HOWDY"
	}

	if err := r.Status().Update(ctx, &machine); err != nil {
		log.Info("Something went wrong while trying to update the resource", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if machine.Status.Status == "OK" {
		return ctrl.Result{}, nil
	}

	if machine.Status.Status == "DELETE" {
		if err := r.Delete(ctx, &machine); err != nil {
			log.Info("Something went wrong while trying to delete the resource", "name", req.NamespacedName)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		log.Info("The Useless Machine has been DELETED!", "name", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	mtype := machine.Spec.MachineType
	switch mtype {
	case "useful":
		machine.Status.Status = "OK"
	case "useless":
		machine.Status.Status = "DELETE"
	case "cheeky":
		r := "o"
		s := `\`

		if len(machine.Status.Status) > 10 {
			machine.Status.Status = strings.Trim(machine.Status.Status[0:10], " ")
		}

		if machine.Status.Status == "HOWDY" {
			machine.Status.Status = r
			break
		}

		if machine.Status.Status == fill(10, r) {
			time.Sleep(time.Second)
			machine.Status.Status = "DELETE"
			break
		}

		time.Sleep(time.Millisecond * 500)
		machine.Status.Status = fill(plusminus(len(machine.Status.Status)), r)
		machine.Status.Status = machine.Status.Status + fill(10-len(machine.Status.Status), " ") + s
	}

	if err := r.Status().Update(ctx, &machine); err != nil {
		log.Info("Something went wrong while trying to update the resource", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Hello from machine ctrl", "name", req.NamespacedName)

	return ctrl.Result{}, nil
}

func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uselessv1.Machine{}).
		Complete(r)
}

func fill(n int, c string) string {
	r := ""
	for i := 0; i < n; i++ {
		r += c
	}
	return r
}

func plusminus(counter int) int {
	if counter == 0 {
		return 1
	}

	if counter == 1 {
		return 2
	}

	n := rand.Intn(2)

	if n == 0 {
		n = -1
	}

	n += counter
	return n
}
