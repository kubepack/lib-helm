package chart

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
)

// https://helm.sh/docs/intro/using_helm/#the-format-and-limitations-of---set
func TestGetChangedValues(t *testing.T) {
	type args struct {
		original map[string]interface{}
		modified map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "single key of simple value",
			args: args{
				original: map[string]interface{}{
					"name": "v1",
				},
				modified: map[string]interface{}{
					"name": "v2",
				},
			},
			want: []string{
				"name=v2",
			},
		},
		{
			name: "multi key of simple value",
			args: args{
				original: map[string]interface{}{
					"a": "b1",
					"c": "d1",
				},
				modified: map[string]interface{}{
					"a": "b2",
					"c": "d2",
				},
			},
			want: []string{
				"a=b2",
				"c=d2",
			},
		},
		{
			name: "nested object",
			args: args{
				original: map[string]interface{}{
					"outer": map[string]interface{}{},
				},
				modified: map[string]interface{}{
					"outer": map[string]interface{}{
						"inner": "value",
					},
				},
			},
			want: []string{
				"outer.inner=value",
			},
		},
		{
			name: "simple array",
			args: args{
				original: map[string]interface{}{
					"name": []interface{}{},
				},
				modified: map[string]interface{}{
					"name": []interface{}{"a", "b", "c"},
				},
			},
			want: []string{
				"name={a, b, c}",
			},
		},
		{
			name: "nested array",
			args: args{
				original: map[string]interface{}{
					"servers": []interface{}{},
				},
				modified: map[string]interface{}{
					"servers": []interface{}{
						map[string]interface{}{
							"port": []interface{}{
								80,
							},
							"host": []interface{}{
								"example",
							},
							"name": "nginx",
						},
					},
				},
			},
			want: []string{
				"servers[0].port={80}",
				"servers[0].host={example}",
				"servers[0].name=nginx",
			},
		},
		{
			name: "array of single key object",
			args: args{
				original: map[string]interface{}{
					"servers": []interface{}{},
				},
				modified: map[string]interface{}{
					"servers": []interface{}{
						map[string]interface{}{
							"port": 80,
						},
						map[string]interface{}{
							"port": 443,
						},
					},
				},
			},
			want: []string{
				"servers[0].port=80",
				"servers[1].port=443",
			},
		},
		{
			name: "array of multi key object",
			args: args{
				original: map[string]interface{}{
					"servers": []interface{}{},
				},
				modified: map[string]interface{}{
					"servers": []interface{}{
						map[string]interface{}{
							"port": 80,
							"name": "http",
						},
						map[string]interface{}{
							"port": 443,
							"name": "https",
						},
					},
				},
			},
			want: []string{
				"servers[0].port=80",
				"servers[0].name=http",
				"servers[1].port=443",
				"servers[1].name=https",
			},
		},
		{
			name: "escape value",
			args: args{
				original: map[string]interface{}{
					"name": "",
				},
				modified: map[string]interface{}{
					"name": "value1,value2",
				},
			},
			want: []string{
				`name=value1\,value2`,
			},
		},
		{
			name: "escape key",
			args: args{
				original: map[string]interface{}{
					"nodeSelector": map[string]interface{}{},
				},
				modified: map[string]interface{}{
					"nodeSelector": map[string]interface{}{
						"kubernetes.io/role": "master",
					},
				},
			},
			want: []string{
				`nodeSelector."kubernetes\.io/role"=master`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChangedValues(tt.args.original, tt.args.modified)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChangedValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !sets.NewString(got...).Equal(sets.NewString(tt.want...)) {
				t.Errorf("GetChangedValues() got = %v, want %v", got, tt.want)
			}
		})
	}
}
