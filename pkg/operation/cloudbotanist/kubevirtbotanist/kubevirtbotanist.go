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
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client"

	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/gardener/gardener/pkg/operation"
	"github.com/gardener/gardener/pkg/operation/common"
	"k8s.io/client-go/tools/clientcmd"
	kubevirtv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

var KubeVirtBotanistScheme = runtime.NewScheme()

func init() {
	KubeVirtBotanistScheme = runtime.NewScheme()
	utilruntime.Must(scheme.AddToScheme(KubeVirtBotanistScheme))
	utilruntime.Must(kubevirtv1.AddToScheme(KubeVirtBotanistScheme))
}

// New takes an operation object <o> and creates a new KubeVirtBotanist object.
func New(o *operation.Operation, purpose string) (*KubeVirtBotanist, error) {
	var cloudProvider gardenv1beta1.CloudProvider
	switch purpose {
	case common.CloudPurposeShoot:
		cloudProvider = o.Shoot.CloudProvider
	case common.CloudPurposeSeed:
		cloudProvider = o.Seed.CloudProvider
	}

	if cloudProvider != gardenv1beta1.CloudProviderKubeVirt {
		return nil, errors.New("cannot instantiate an KubeVirt botanist if neither Shoot nor Seed cluster specifies KubeVirt")
	}

	clientConfig, err := clientcmd.NewClientConfigFromBytes(o.Shoot.Secret.Data[KubeConfig])
	if err != nil {
		return nil, fmt.Errorf("failed to create clientconfig from bytes: %v", err)
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get RESTconfig for clientConfig: %v", err)
	}

	c, err := client.New(restConfig, client.Options{
		Scheme: KubeVirtBotanistScheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client from config: %v", err)
	}

	return &KubeVirtBotanist{
		client:            c,
		Operation:         o,
		CloudProviderName: "kubevirt",
	}, nil
}

// GetCloudProviderName returns the Kubernetes cloud provider name for this cloud.
func (b *KubeVirtBotanist) GetCloudProviderName() string {
	return b.CloudProviderName
}
