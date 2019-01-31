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

// DeployInfrastructure kicks off a Terraform job which deploys the infrastructure.
func (b *KubeVirtBotanist) DeployInfrastructure() error {
	return nil
}

// DestroyInfrastructure kicks off a Terraform job which destroys the infrastructure.
func (b *KubeVirtBotanist) DestroyInfrastructure() error {
	return nil
}

// generateTerraformInfraVariablesEnvironment generates the environment containing the credentials which
// are required to validate/apply/destroy the Terraform configuration. These environment must contain
// Terraform variables which are prefixed with TF_VAR_.
func (b *KubeVirtBotanist) generateTerraformInfraVariablesEnvironment() []map[string]interface{} {
	return nil
}

// generateTerraformInfraConfig creates the Terraform variables and the Terraform config (for the infrastructure)
// and returns them (these values will be stored as a ConfigMap and a Secret in the Garden cluster.
func (b *KubeVirtBotanist) generateTerraformInfraConfig(createRouter bool, routerID string) map[string]interface{} {
	return nil
}

// DeployBackupInfrastructure kicks off a Terraform job which creates the infrastructure resources for backup.
func (b *KubeVirtBotanist) DeployBackupInfrastructure() error {
	return nil
}

// DestroyBackupInfrastructure kicks off a Terraform job which destroys the infrastructure for backup.
func (b *KubeVirtBotanist) DestroyBackupInfrastructure() error {
	return nil
}

// generateTerraformBackupVariablesEnvironment generates the environment containing the credentials which
// are required to validate/apply/destroy the Terraform configuration. These environment must contain
// Terraform variables which are prefixed with TF_VAR_.
func (b *KubeVirtBotanist) generateTerraformBackupVariablesEnvironment() []map[string]interface{} {
	return nil
}

// generateTerraformBackupConfig creates the Terraform variables and the Terraform config (for the backup)
// and returns them.
func (b *KubeVirtBotanist) generateTerraformBackupConfig() map[string]interface{} {
	return nil
}
