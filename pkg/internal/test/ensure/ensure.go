/*
Copyright The Helm Authors.

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

package ensure

import (
	"io/ioutil"
	"os"
	"testing"

	"helm.sh/helm/v3/pkg/helmpath/xdg"
)

// HelmHome sets up a Helm Home in a temp dir.
func HelmHome(t *testing.T) func() {
	t.Helper()
	base := TempDir(t)
	os.Setenv(xdg.CacheHomeEnvVar, base)
	os.Setenv(xdg.ConfigHomeEnvVar, base)
	os.Setenv(xdg.DataHomeEnvVar, base)
	return func() {
		os.RemoveAll(base)
	}
}

// TempDir ensures a scratch test directory for unit testing purposes.
func TempDir(t *testing.T) string {
	t.Helper()
	d, err := ioutil.TempDir("", "helm")
	if err != nil {
		t.Fatal(err)
	}
	return d
}
