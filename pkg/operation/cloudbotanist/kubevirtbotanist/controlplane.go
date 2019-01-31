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

	"github.com/gardener/gardener/pkg/operation/common"
)

// GenerateCloudProviderConfig generates the KubeVirt cloud provider config.
// See this for more details:
// https://github.com/kubeforge/cloud-provider-kubevirt/tree/master/pkg/cloudprovider/kubevirt
func (b *KubeVirtBotanist) GenerateCloudProviderConfig() (string, error) {
	return string(b.Shoot.Secret.Data[KubeConfig]), nil
}

// RefreshCloudProviderConfig refreshes the cloud provider credentials in the existing cloud
// provider config.
func (b *KubeVirtBotanist) RefreshCloudProviderConfig(currentConfig map[string]string) map[string]string {
	var (
		existing  = currentConfig[common.CloudProviderConfigMapKey]
		updated   = existing
		separator = "="
	)
	updated = common.ReplaceCloudProviderConfigKey(updated, separator, "kubeConfig", string(b.Shoot.Secret.Data[KubeConfig]))
	return map[string]string{
		common.CloudProviderConfigMapKey: updated,
	}
}

// GenerateKubeAPIServerServiceConfig generates the cloud provider specific values which are required to render the
// Service manifest of the kube-apiserver-service properly.
func (b *KubeVirtBotanist) GenerateKubeAPIServerServiceConfig() (map[string]interface{}, error) {
	return nil, nil
}

// GenerateKubeAPIServerExposeConfig defines the cloud provider specific values which configure how the kube-apiserver
// is exposed to the public.
func (b *KubeVirtBotanist) GenerateKubeAPIServerExposeConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"advertiseAddress": b.APIServerAddress,
		"additionalParameters": []string{
			fmt.Sprintf("--external-hostname=%s", b.APIServerAddress),
		},
	}, nil
}

// GenerateKubeAPIServerConfig generates the cloud provider specific values which are required to render the
// Deployment manifest of the kube-apiserver properly.
func (b *KubeVirtBotanist) GenerateKubeAPIServerConfig() (map[string]interface{}, error) {
	return nil, nil
}

// GenerateCloudControllerManagerConfig generates the cloud provider specific values which are required to
// render the Deployment manifest of the cloud-controller-manager properly.
func (b *KubeVirtBotanist) GenerateCloudControllerManagerConfig() (map[string]interface{}, string, error) {
	return nil, common.CloudControllerManagerDeploymentName, nil
}

// GenerateCSIConfig generates the configuration for CSI charts
func (b *KubeVirtBotanist) GenerateCSIConfig() (map[string]interface{}, error) {
	return nil, nil
}

// GenerateKubeControllerManagerConfig generates the cloud provider specific values which are required to
// render the Deployment manifest of the kube-controller-manager properly.
func (b *KubeVirtBotanist) GenerateKubeControllerManagerConfig() (map[string]interface{}, error) {
	return nil, nil
}

// GenerateKubeSchedulerConfig generates the cloud provider specific values which are required to render the
// Deployment manifest of the kube-scheduler properly.
func (b *KubeVirtBotanist) GenerateKubeSchedulerConfig() (map[string]interface{}, error) {
	return nil, nil
}

// GenerateEtcdBackupConfig returns the etcd backup configuration for the etcd Helm chart.
func (b *KubeVirtBotanist) GenerateEtcdBackupConfig() (map[string][]byte, map[string]interface{}, error) {
	return make(map[string][]byte), make(map[string]interface{}), nil
}

// DeployCloudSpecificControlPlane does currently nothing for OpenStack.
func (b *KubeVirtBotanist) DeployCloudSpecificControlPlane() error {
	return nil
}
