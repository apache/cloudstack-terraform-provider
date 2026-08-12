//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestParseKubernetesClusterConfig(t *testing.T) {
	tests := []struct {
		name       string
		configData string
		expected   kubernetesClusterCredentials
		expectErr  string
	}{
		{
			name: "kubeadm admin.conf as returned by CloudStack",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
preferences: {}
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://10.1.1.1:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
users:
- name: kubernetes-admin
  user:
    client-certificate-data: %s
    client-key-data: %s
`, testBase64("ca-certificate"), testBase64("client-certificate"), testBase64("client-key")),
			expected: kubernetesClusterCredentials{
				Endpoint:             "https://10.1.1.1:6443",
				ClusterCACertificate: "ca-certificate",
				ClientCertificate:    "client-certificate",
				ClientKey:            "client-key",
			},
		},
		{
			// Guards against blindly taking clusters[0] / users[0].
			name:       "current context selects the second cluster and user",
			configData: testKubeConfigTwoClusters("second@second"),
			expected: kubernetesClusterCredentials{
				Endpoint:             "https://10.2.2.2:6443",
				ClusterCACertificate: "second-ca-certificate",
				ClientCertificate:    "second-client-certificate",
				ClientKey:            "second-client-key",
			},
		},
		{
			name:       "absent current context falls back to the first entry",
			configData: testKubeConfigTwoClusters(""),
			expected: kubernetesClusterCredentials{
				Endpoint:             "https://10.1.1.1:6443",
				ClusterCACertificate: "first-ca-certificate",
				ClientCertificate:    "first-client-certificate",
				ClientKey:            "first-client-key",
			},
		},
		{
			name:       "unknown current context falls back to the first entry",
			configData: testKubeConfigTwoClusters("missing@missing"),
			expected: kubernetesClusterCredentials{
				Endpoint:             "https://10.1.1.1:6443",
				ClusterCACertificate: "first-ca-certificate",
				ClientCertificate:    "first-client-certificate",
				ClientKey:            "first-client-key",
			},
		},
		{
			// The context resolves, but names a cluster and user that the
			// kubeconfig does not define. With only one entry of each,
			// there's only one reasonable candidate, so this still falls
			// back rather than erroring.
			name: "current context naming a missing cluster falls back to the only entry",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: third@third
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://10.1.1.1:6443
  name: first
contexts:
- context:
    cluster: third
    user: third
  name: third@third
users:
- name: first
  user:
    client-certificate-data: %s
    client-key-data: %s
`, testBase64("first-ca-certificate"),
				testBase64("first-client-certificate"), testBase64("first-client-key")),
			expected: kubernetesClusterCredentials{
				Endpoint:             "https://10.1.1.1:6443",
				ClusterCACertificate: "first-ca-certificate",
				ClientCertificate:    "first-client-certificate",
				ClientKey:            "first-client-key",
			},
		},
		{
			// Here the context resolves and names a cluster absent from a
			// list of more than one: guessing which of the two is intended
			// could silently pair the wrong cluster/user together, so this
			// must be a hard error instead of a silent fallback. Clusters are
			// resolved before users, so no users section is needed to reach
			// this error.
			name: "current context naming an undefined cluster among several is an error",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: third@third
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://10.1.1.1:6443
  name: first
- cluster:
    certificate-authority-data: %s
    server: https://10.2.2.2:6443
  name: second
contexts:
- context:
    cluster: third
    user: third
  name: third@third
`, testBase64("first-ca-certificate"), testBase64("second-ca-certificate")),
			expectErr: "which is not defined",
		},
		{
			// The same hard error, but reached through the user branch
			// instead of the cluster branch, so a future edit that swaps
			// clusterName/userName or config.Clusters/config.Users between
			// the two findKubeConfigEntry call sites would still be caught.
			name: "current context naming an undefined user among several is an error",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: kubernetes@kubernetes
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://10.1.1.1:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: third
  name: kubernetes@kubernetes
users:
- name: first
  user:
    client-certificate-data: %s
    client-key-data: %s
- name: second
  user:
    client-certificate-data: %s
    client-key-data: %s
`, testBase64("ca-certificate"),
				testBase64("first-client-certificate"), testBase64("first-client-key"),
				testBase64("second-client-certificate"), testBase64("second-client-key")),
			expectErr: "which is not defined",
		},
		{
			// A kubeconfig without client certificates must not fail the data
			// source, so that config_data stays usable.
			name: "no users leaves the client credentials empty",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://10.1.1.1:6443
  name: kubernetes
`, testBase64("ca-certificate")),
			expected: kubernetesClusterCredentials{
				Endpoint:             "https://10.1.1.1:6443",
				ClusterCACertificate: "ca-certificate",
			},
		},
		{
			name: "no clusters leaves the endpoint empty",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
users:
- name: kubernetes-admin
  user:
    client-certificate-data: %s
    client-key-data: %s
`, testBase64("client-certificate"), testBase64("client-key")),
			expected: kubernetesClusterCredentials{
				ClientCertificate: "client-certificate",
				ClientKey:         "client-key",
			},
		},
		{
			// The empty-config guard lives in the read function, not here.
			name:       "empty document yields no credentials and no error",
			configData: "",
			expected:   kubernetesClusterCredentials{},
		},
		{
			name:       "invalid yaml",
			configData: "\tthis is not a kubeconfig",
			expectErr:  "Invalid kubeconfig",
		},
		{
			name: "certificate authority data is not base64",
			configData: `apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: "@@@ not base64 @@@"
    server: https://10.1.1.1:6443
  name: kubernetes
`,
			expectErr: "Failed to decode certificate-authority-data",
		},
		{
			name: "client key data is not base64",
			configData: fmt.Sprintf(`apiVersion: v1
kind: Config
users:
- name: kubernetes-admin
  user:
    client-certificate-data: %s
    client-key-data: "@@@ not base64 @@@"
`, testBase64("client-certificate")),
			expectErr: "Failed to decode client-key-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, err := parseKubernetesClusterConfig(tt.configData)

			if tt.expectErr != "" {
				if err == nil {
					t.Fatalf("Expected an error containing %q, got none", tt.expectErr)
				}
				if !strings.Contains(err.Error(), tt.expectErr) {
					t.Fatalf("Expected an error containing %q, got: %s", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if *credentials != tt.expected {
				t.Errorf("Expected %+v, got %+v", tt.expected, *credentials)
			}
		})
	}
}

func testBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// testKubeConfigTwoClusters renders a kubeconfig holding two clusters and two
// users, so that context resolution can be told apart from taking the first
// entry. currentContext selects which one context resolution should pick; an
// empty string omits the current-context key entirely. Every %s verb appears
// in the same order as its argument below, so the substitution can be checked
// by reading both top to bottom in lockstep.
func testKubeConfigTwoClusters(currentContext string) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: %s
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://10.1.1.1:6443
  name: first
- cluster:
    certificate-authority-data: %s
    server: https://10.2.2.2:6443
  name: second
contexts:
- context:
    cluster: first
    user: first
  name: first@first
- context:
    cluster: second
    user: second
  name: second@second
users:
- name: first
  user:
    client-certificate-data: %s
    client-key-data: %s
- name: second
  user:
    client-certificate-data: %s
    client-key-data: %s
`,
		currentContext,
		testBase64("first-ca-certificate"),
		testBase64("second-ca-certificate"),
		testBase64("first-client-certificate"),
		testBase64("first-client-key"),
		testBase64("second-client-certificate"),
		testBase64("second-client-key"))
}
