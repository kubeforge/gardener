// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubevirtbotanist

import (
	"fmt"

	"github.com/gardener/gardener/pkg/apis/garden/v1beta1"

	"github.com/gardener/gardener/pkg/operation"
	"github.com/gardener/gardener/pkg/operation/common"
	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

// GetMachineClassInfo returns the name of the class kind, the plural of it and the name of the Helm chart which
// contains the machine class template.
func (b *KubeVirtBotanist) GetMachineClassInfo() (classKind, classPlural, classChartName string) {
	classKind = "KubeVirtMachineClass"
	classPlural = "kubevirtmachineclasses"
	classChartName = "kubevirt-machineclass"
	return
}

// GenerateMachineClassSecretData generates the secret data for the machine class secret (except the userData field
// which is computed elsewhere).
func (b *KubeVirtBotanist) GenerateMachineClassSecretData() map[string][]byte {
	return map[string][]byte{
		machinev1alpha1.KubeVirtNamespace:  b.Shoot.Secret.Data[Namespace],
		machinev1alpha1.KubeVirtKubeConfig: b.Shoot.Secret.Data[Kubeconfig],
	}
}

// getMachineTypeFromString returns a MachineType object from the CloudProfile for a given machine type name.
func (b *KubeVirtBotanist) getMachineTypeFromString(machineType string) (v1beta1.MachineType, error) {
	for _, mType := range b.Shoot.GetMachineTypesFromCloudProfile() {
		if mType.Name == machineType {
			return mType, nil
		}
	}
	return v1beta1.MachineType{}, fmt.Errorf("could not find the machinetype %v in the %v cloudprofile", machineType, b.Shoot.CloudProfile.Name)
}

// GenerateMachineConfig generates the configuration values for the cloud-specific machine class Helm chart. It
// also generates a list of corresponding MachineDeployments. The provided worker groups will be distributed over
// the desired availability zones. It returns the computed list of MachineClasses and MachineDeployments.
func (b *KubeVirtBotanist) GenerateMachineConfig() ([]map[string]interface{}, operation.MachineDeployments, error) {
	var (
		workers = b.Shoot.Info.Spec.Cloud.KubeVirt.Workers

		machineDeployments = operation.MachineDeployments{}
		machineClasses     = []map[string]interface{}{}
	)

	for _, worker := range workers {
		machineType, err := b.getMachineTypeFromString(worker.MachineType)
		if err != nil {
			return nil, nil, err
		}

		machineClassSpec := map[string]interface{}{
			"region":         b.Shoot.Info.Spec.Cloud.Region,
			"cores":          machineType.CPU,
			"memory":         machineType.Memory,
			"imageName":      b.Shoot.Info.Spec.Cloud.KubeVirt.MachineImage.Image,
			"podNetworkCidr": b.Shoot.GetPodNetwork(),
			"tags": map[string]string{
				fmt.Sprintf("kubernetes.io-cluster-%s", b.Shoot.SeedNamespace): "1",
				"kubernetes.io-role-node":                                      "1",
			},
			"secret": map[string]interface{}{
				"cloudConfig": b.Shoot.CloudConfigMap[worker.Name].Downloader.Content,
			},
		}

		var (
			machineClassSpecHash = common.MachineClassHash(machineClassSpec, b.Shoot.KubernetesMajorMinorVersion)
			deploymentName       = fmt.Sprintf("%s-%s-z%d", b.Shoot.SeedNamespace, worker.Name, 1)
			className            = fmt.Sprintf("%s-%s", deploymentName, machineClassSpecHash)
			secretData           = b.GenerateMachineClassSecretData()
		)

		machineDeployments = append(machineDeployments, operation.MachineDeployment{
			Name:           deploymentName,
			ClassName:      className,
			Minimum:        common.DistributeOverZones(0, worker.AutoScalerMin, 1),
			Maximum:        common.DistributeOverZones(0, worker.AutoScalerMax, 1),
			MaxSurge:       common.DistributePositiveIntOrPercent(0, *worker.MaxSurge, 1, worker.AutoScalerMax),
			MaxUnavailable: common.DistributePositiveIntOrPercent(0, *worker.MaxUnavailable, 1, worker.AutoScalerMin),
		})

		machineClassSpec["name"] = className
		machineClassSpec["secret"].(map[string]interface{})[Namespace] = string(secretData[machinev1alpha1.KubeVirtNamespace])
		machineClassSpec["secret"].(map[string]interface{})[Kubeconfig] = string(secretData[machinev1alpha1.KubeVirtKubeConfig])

		machineClasses = append(machineClasses, machineClassSpec)
	}

	return machineClasses, machineDeployments, nil
}

// ListMachineClasses returns two sets of strings whereas the first contains the names of all machine
// classes, and the second the names of all referenced secrets.
func (b *KubeVirtBotanist) ListMachineClasses() (sets.String, sets.String, error) {
	var (
		classNames  = sets.NewString()
		secretNames = sets.NewString()
	)

	existingMachineClasses, err := b.K8sSeedClient.Machine().MachineV1alpha1().KubeVirtMachineClasses(b.Shoot.SeedNamespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	for _, existingMachineClass := range existingMachineClasses.Items {
		if existingMachineClass.Spec.SecretRef == nil {
			return nil, nil, fmt.Errorf("could not find secret reference in class %s", existingMachineClass.Name)
		}

		secretNames.Insert(existingMachineClass.Spec.SecretRef.Name)
		classNames.Insert(existingMachineClass.Name)
	}

	return classNames, secretNames, nil
}

// CleanupMachineClasses deletes all machine classes which are not part of the provided list <existingMachineDeployments>.
func (b *KubeVirtBotanist) CleanupMachineClasses(existingMachineDeployments operation.MachineDeployments) error {
	existingMachineClasses, err := b.K8sSeedClient.Machine().MachineV1alpha1().KubeVirtMachineClasses(b.Shoot.SeedNamespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, existingMachineClass := range existingMachineClasses.Items {
		if !existingMachineDeployments.ContainsClass(existingMachineClass.Name) {
			if err := b.K8sSeedClient.Machine().MachineV1alpha1().KubeVirtMachineClasses(b.Shoot.SeedNamespace).Delete(existingMachineClass.Name, &metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	return nil
}
