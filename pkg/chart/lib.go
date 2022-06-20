package chart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/pkg/errors"
	kj "gomodules.xyz/encoding/json"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	"sigs.k8s.io/yaml"
)

func GetValuesDiff(original, modified map[string]interface{}) (map[string]interface{}, error) {
	return getValuesDiff(original, modified, "", nil)
}

func getValuesDiff(original, modified map[string]interface{}, prefix string, diff map[string]interface{}) (map[string]interface{}, error) {
	if diff == nil {
		diff = map[string]interface{}{}
	}

	for k, v := range modified {
		curKey := ""
		if prefix == "" {
			curKey = escapeKey(k)
		} else {
			curKey = prefix + "." + escapeKey(k)
		}

		switch val := v.(type) {
		case map[string]interface{}:
			oVal, ok := original[k].(map[string]interface{})
			if !ok {
				oVal = map[string]interface{}{}
			}

			d2, err := getValuesDiff(oVal, val, curKey, nil)
			if err != nil {
				return nil, err
			}
			if len(d2) > 0 {
				diff[k] = d2
			}
		case []interface{}, string, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, json.Number, nil:
			if !reflect.DeepEqual(original[k], val) {
				diff[k] = val
			}
		default:
			return nil, fmt.Errorf("unknown type %v with value %v", reflect.TypeOf(v), v)
		}
	}

	// https://github.com/kubepack/lib-helm/blob/32de2acacbfb84f57d4a66c6d896360eb664399c/pkg/values/options.go#L133
	for k, v := range original {
		if _, found := modified[k]; !found {
			curKey := ""
			if prefix == "" {
				curKey = escapeKey(k)
			} else {
				curKey = prefix + "." + escapeKey(k)
			}

			// TODO: how does Helm merge --values remove keys?
			// diff[k] = nil
			return nil, fmt.Errorf("key %s is missing in the modified values, original values %v", curKey, v)
		}
	}
	return diff, nil
}

func GetValuesDiffYAML(orig, od interface{}) ([]byte, error) {
	origMap, err := toJson(orig)
	if err != nil {
		return nil, err
	}
	modMap, err := toJson(od)
	if err != nil {
		return nil, err
	}

	diff, err := GetValuesDiff(origMap, modMap)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(diff)
}

func toJson(v interface{}) (map[string]interface{}, error) {
	data, err := kj.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out map[string]interface{}
	err = kj.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func GetChangedValues(original, modified map[string]interface{}) ([]string, error) {
	cmds, err := getChangedValues(original, modified, "", nil)
	if err != nil {
		return nil, err
	}
	sort.Strings(cmds)
	return cmds, nil
}

func getChangedValues(original, modified map[string]interface{}, prefix string, cmds []string) ([]string, error) {
	for k, v := range modified {
		curKey := ""
		if prefix == "" {
			curKey = escapeKey(k)
		} else {
			curKey = prefix + "." + escapeKey(k)
		}

		switch val := v.(type) {
		case map[string]interface{}:
			oVal, ok := original[k].(map[string]interface{})
			if !ok {
				oVal = map[string]interface{}{}
			}
			next, err := getChangedValues(oVal, val, curKey, nil)
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, next...)
		case []interface{}:
			if !reflect.DeepEqual(v, original[k]) {
				if len(val) == 0 {
					cmds = append(cmds, fmt.Sprintf("%s=null", curKey))
					continue
				}

				if isSimpleArray(val) {
					s, err := PrintArray(val)
					if err != nil {
						return nil, fmt.Errorf("failed to print simple array %v, reason: %v", v, err)
					}
					cmds = append(cmds, fmt.Sprintf("%s=%s", curKey, s))
					continue
				}

				for i, element := range val {
					em, ok := element.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("%s[%d] element is not a map", curKey, i)
					}
					next, err := getChangedValues(map[string]interface{}{}, em, fmt.Sprintf("%s[%d]", curKey, i), nil)
					if err != nil {
						return nil, err
					}
					cmds = append(cmds, next...)
				}
			}
		case string:
			if !reflect.DeepEqual(original[k], val) {
				cmds = append(cmds, fmt.Sprintf("%s=%v", curKey, escapeValue(val)))
			}
		case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, json.Number:
			if !reflect.DeepEqual(original[k], val) {
				cmds = append(cmds, fmt.Sprintf("%s=%v", curKey, val))
			}
		case nil:
			if !reflect.DeepEqual(original[k], val) {
				cmds = append(cmds, fmt.Sprintf("%s=null", curKey))
			}
		default:
			return nil, fmt.Errorf("unknown type %v with value %v", reflect.TypeOf(v), v)
		}
	}

	for k := range original {
		if _, found := modified[k]; !found {
			curKey := ""
			if prefix == "" {
				curKey = escapeKey(k)
			} else {
				curKey = prefix + "." + escapeKey(k)
			}

			cmds = append(cmds, fmt.Sprintf("%s=null", curKey))
		}
	}
	return cmds, nil
}

// kubernetes.io/role becomes "kubernetes\.io/role"
func escapeKey(s string) string {
	return shellescape.Quote(strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `.`, `\.`))
}

// "value1,value2" becomes value1\,value2
func escapeValue(s string) string {
	return shellescape.Quote(strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `,`, `\,`))
}

func isSimpleArray(a []interface{}) bool {
	for i := range a {
		switch a[i].(type) {
		case string, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, nil, json.Number:
		default:
			return false
		}
	}
	return true
}

func PrintArray(a []interface{}) (string, error) {
	var buf bytes.Buffer
	buf.WriteRune('{')
	for i := range a {
		switch v := a[i].(type) {
		case string:
			if i > 0 {
				buf.WriteString(", ")
			}
			_, err := fmt.Fprint(&buf, escapeValue(v))
			if err != nil {
				return "", err
			}
		case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, json.Number:
			if i > 0 {
				buf.WriteString(", ")
			}
			_, err := fmt.Fprint(&buf, v)
			if err != nil {
				return "", err
			}
		case nil:
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString("null")
		default:
			return "", fmt.Errorf("[%d] holds a complex type %v", i, reflect.TypeOf(a[i]))
		}
	}
	buf.WriteRune('}')
	return buf.String(), nil
}

// helm.sh/helm/v3/pkg/action/install.go
const notesFileSuffix = "NOTES.txt"

// RenderResources renders the templates in a chart
func RenderResources(ch *chart.Chart, caps *chartutil.Capabilities, values chartutil.Values) ([]*release.Hook, []releaseutil.Manifest, error) {
	hs := []*release.Hook{}
	b := bytes.NewBuffer(nil)

	if ch.Metadata.KubeVersion != "" {
		if !chartutil.IsCompatibleRange(ch.Metadata.KubeVersion, caps.KubeVersion.String()) {
			return hs, nil, errors.Errorf("chart requires kubeVersion: %s which is incompatible with Kubernetes %s", ch.Metadata.KubeVersion, caps.KubeVersion.String())
		}
	}

	files, err := engine.Render(ch, values)
	if err != nil {
		return hs, nil, err
	}

	for k := range files {
		if strings.HasSuffix(k, notesFileSuffix) {
			delete(files, k)
		}
	}

	// Sort hooks, manifests, and partials. Only hooks and manifests are returned,
	// as partials are not used after renderer.Render. Empty manifests are also
	// removed here.
	hs, manifests, err := releaseutil.SortManifests(files, caps.APIVersions, releaseutil.InstallOrder)
	if err != nil {
		// By catching parse errors here, we can prevent bogus releases from going
		// to Kubernetes.
		//
		// We return the files as a big blob of data to help the user debug parser
		// errors.
		for name, content := range files {
			if strings.TrimSpace(content) == "" {
				continue
			}
			fmt.Fprintf(b, "---\n# Source: %s\n%s\n", name, content)
		}
		return hs, manifests, err
	}

	return hs, manifests, nil
}

func IsEvent(events []release.HookEvent, x release.HookEvent) bool {
	for _, event := range events {
		if event == x {
			return true
		}
	}
	return false
}

// IsChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func IsChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
