package values

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
)

// https://helm.sh/docs/intro/using_helm/#the-format-and-limitations-of---set
func TestGetChangedValues(t *testing.T) {
	type args struct {
		original map[string]any
		modified map[string]any
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
				original: map[string]any{
					"name": "v1",
				},
				modified: map[string]any{
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
				original: map[string]any{
					"a": "b1",
					"c": "d1",
				},
				modified: map[string]any{
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
				original: map[string]any{
					"outer": map[string]any{},
				},
				modified: map[string]any{
					"outer": map[string]any{
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
				original: map[string]any{
					"name": []any{},
				},
				modified: map[string]any{
					"name": []any{"a", "b", "c"},
				},
			},
			want: []string{
				"name={a, b, c}",
			},
		},
		{
			name: "nested array",
			args: args{
				original: map[string]any{
					"servers": []any{},
				},
				modified: map[string]any{
					"servers": []any{
						map[string]any{
							"port": []any{
								80,
							},
							"host": []any{
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
				original: map[string]any{
					"servers": []any{},
				},
				modified: map[string]any{
					"servers": []any{
						map[string]any{
							"port": 80,
						},
						map[string]any{
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
				original: map[string]any{
					"servers": []any{},
				},
				modified: map[string]any{
					"servers": []any{
						map[string]any{
							"port": 80,
							"name": "http",
						},
						map[string]any{
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
				original: map[string]any{
					"name": "",
				},
				modified: map[string]any{
					"name": "value1,value2",
				},
			},
			want: []string{
				`name='value1\,value2'`,
			},
		},
		{
			name: "escape key",
			args: args{
				original: map[string]any{
					"nodeSelector": map[string]any{},
				},
				modified: map[string]any{
					"nodeSelector": map[string]any{
						"kubernetes.io/role": "master",
					},
				},
			},
			want: []string{
				`nodeSelector.'kubernetes\.io/role'=master`,
			},
		},
		{
			name: "delete values",
			args: args{
				original: map[string]any{
					"annotations": map[string]any{},
				},
				modified: map[string]any{
					"annotations": map[string]any{
						"helm.sh/hook": nil,
					},
				},
			},
			want: []string{
				`annotations.'helm\.sh/hook'=null`,
			},
			wantErr: false,
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

func TestGetValuesMapDiff(t *testing.T) {
	type args struct {
		original map[string]any
		modified map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "delete values",
			args: args{
				original: map[string]any{
					"annotations": map[string]any{},
				},
				modified: map[string]any{
					"annotations": map[string]any{
						"helm.sh/hook": nil,
					},
				},
			},
			want: map[string]any{
				"annotations": map[string]any{
					"helm.sh/hook": nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValuesMapDiff(tt.args.original, tt.args.modified)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValuesMapDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValuesMapDiff() got = %v, want %v", got, tt.want)
			}
		})
	}
}
