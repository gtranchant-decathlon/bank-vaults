// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/banzaicloud/bank-vaults/internal/injector"
	"github.com/banzaicloud/bank-vaults/pkg/sdk/vault"
)

func getDataFromVault(data map[string]string, vaultClient *vault.Client, vaultConfig VaultConfig, logger logrus.FieldLogger) (map[string]string, error) {
	vaultData := make(map[string]string, len(data))

	for key, value := range data {
		value = strings.Map(removePunctuation, value)
		data[key] = value
	}

	inject := func(key, value string) {
		vaultData[key] = value
	}

	config := injector.Config{
		TransitKeyID: vaultConfig.TransitKeyID,
		TransitPath:  vaultConfig.TransitPath,
	}
	secretInjector := injector.NewSecretInjector(config, vaultClient, nil, logger)

	return vaultData, secretInjector.InjectSecretsFromVault(data, inject)
}

func hasVaultPrefix(value string) bool {
	return strings.HasPrefix(value, "vault:") || strings.HasPrefix(value, ">>vault:")
}