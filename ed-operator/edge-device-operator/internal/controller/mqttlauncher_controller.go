/*
Copyright 2025.

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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "ed-operator/api/v1"
	domain_commands "ed-operator/internal/domain/commands"
	"ed-operator/internal/domain/dto"
)

// MQTTLauncherReconciler reconciles a MQTTLauncher object
type MQTTLauncherReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	mqttClient mqtt.Client
}

func (r *MQTTLauncherReconciler) PublishResult(ctx context.Context, Success bool, Message string) error {
	log := log.FromContext(ctx)
	if Success {
		log.Info("Command executed successfully")
	} else {
		log.Error(nil, "Command execution failed. Message: "+Message)
	}

	log.Info("Publishing result...")
	payload, err := json.Marshal(dto.ResultDto{Success: Success, Message: Message})
	if err != nil {
		log.Error(err, "Failed to marshal result")
		return nil
	}
	log.Info(string(payload))
	if token := r.mqttClient.Publish("deployments/results", 0, false, payload); token.Wait() && token.Error() != nil {
		log.Error(token.Error(), "Failed to publish result")
		return nil
	}
	log.Info("Result published")
	return nil
}

func (r *MQTTLauncherReconciler) Start(ctx context.Context) error {
	log := log.FromContext(ctx)

	log.Info("Preparing MQTT client...")
	opts := mqtt.NewClientOptions().AddBroker("tcp://nanomq:1883")
	rand.Seed(time.Now().UnixNano())
	opts.SetClientID(fmt.Sprintf("mqtt-operator-%d", rand.Intn(10000)))
	r.mqttClient = mqtt.NewClient(opts)

	log.Info("Connecting MQTT client...")
	if token := r.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error(), "Failed to connect MQTT client")
		return token.Error()
	}
	log.Info("MQTT client connected")

	log.Info("Subscribing MQTT client...")
	commandFactory := domain_commands.NewCommandFactory(ctx, r.Client)
	r.mqttClient.Subscribe("deployments/start", 0, func(c mqtt.Client, m mqtt.Message) {
		log.Info("Received message: " + string(m.Payload()))
		var commandDTO dto.CommandDTO
		if err := json.Unmarshal(m.Payload(), &commandDTO); err != nil {
			r.PublishResult(ctx, false, "Failed to interpret command: "+err.Error())
			return
		}

		log.Info("Executing command: " + commandDTO.Command)
		command, err := commandFactory.Make(commandDTO)
		if err == nil {
			if cerr := command.Execute(); cerr != nil {
				r.PublishResult(ctx, false, "Failed to execute command: "+cerr.Error())
				return
			}
			r.PublishResult(ctx, true, "")
			return
		}

		r.PublishResult(ctx, false, "Failed to execute command: "+err.Error())
	})

	return nil
}

// +kubebuilder:rbac:groups=core.apps.local,resources=mqttlaunchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.apps.local,resources=mqttlaunchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.apps.local,resources=mqttlaunchers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MQTTLauncher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *MQTTLauncherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MQTTLauncherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.MQTTLauncher{}).
		Named("mqttlauncher").
		Complete(r)
}
