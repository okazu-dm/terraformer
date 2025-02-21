// Copyright 2019 The Terraformer Authors.
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

package cloudflare

import (
	"errors"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

const cloudflareProviderVersion = "~> 1.16"

type CloudflareProvider struct {
	terraform_utils.Provider
}

func (p *CloudflareProvider) Init(args []string) error {
	return nil
}

func (p *CloudflareProvider) GetName() string {
	return "cloudflare"
}

func (p *CloudflareProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			"cloudflare": map[string]interface{}{
				"version": cloudflareProviderVersion,
			},
		},
	}
}

func (CloudflareProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}

func (p *CloudflareProvider) GetSupportedService() map[string]terraform_utils.ServiceGenerator {
	return map[string]terraform_utils.ServiceGenerator{
		"dns":      &DNSGenerator{},
		"firewall": &FirewallGenerator{},
		//"access":   &AccessGenerator{},
	}
}

func (p *CloudflareProvider) InitService(serviceName string) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New("cloudflare: " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetProviderName(p.GetName())

	return nil
}
