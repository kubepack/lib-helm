/*
Copyright 2021 The Flux authors

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

package repository

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "with slash",
			url:     "http://example.com/",
			want:    "http://example.com/",
			wantErr: false,
		},
		{
			name:    "without slash",
			url:     "http://example.com",
			want:    "http://example.com/",
			wantErr: false,
		},
		{
			name:    "double slash",
			url:     "http://example.com//",
			want:    "http://example.com/",
			wantErr: false,
		},
		{
			name:    "empty",
			url:     "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "oci with slash",
			url:     "oci://example.com/",
			want:    "oci://example.com",
			wantErr: false,
		},
		{
			name:    "oci double slash",
			url:     "oci://example.com//",
			want:    "oci://example.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
