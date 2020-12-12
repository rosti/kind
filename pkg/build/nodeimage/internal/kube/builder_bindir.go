/*
Copyright 2020 The Kubernetes Authors.

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

package kube

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/log"
)

// bindirBuilder implements Bits for pre-built Kubernetes binaries placed in a local dir
type bindirBuilder struct {
	binaryDir string
}

var _ Builder = &bindirBuilder{}

// NewBindirBuilder returns a new Bits backed by a pre-built Kubernetes binaries,
// given kubeRoot, the path to the kubernetes binary directory
func NewBindirBuilder(logger log.Logger, kubeRoot, arch string) (Builder, error) {
	return &bindirBuilder{
		binaryDir: kubeRoot,
	}, nil
}

// Build implements Bits.Build
func (b *bindirBuilder) Build() (Bits, error) {
	versionFile := filepath.Join(b.binaryDir, "kube-apiserver.docker_tag")
	rawVersionBytes, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return nil, errors.Wrapf(err, "could not obtain kubernetes version from %s", versionFile)
	}

	return &bits{
		binaryPaths: []string{
			filepath.Join(b.binaryDir, "kubeadm"),
			filepath.Join(b.binaryDir, "kubelet"),
			filepath.Join(b.binaryDir, "kubectl"),
		},
		imagePaths: []string{
			filepath.Join(b.binaryDir, "kube-apiserver.tar"),
			filepath.Join(b.binaryDir, "kube-controller-manager.tar"),
			filepath.Join(b.binaryDir, "kube-scheduler.tar"),
			filepath.Join(b.binaryDir, "kube-proxy.tar"),
		},
		version: strings.TrimSpace(string(rawVersionBytes)),
	}, nil
}
