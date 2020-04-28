package chart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func GetChangedValues(original map[string]interface{}, modified map[string]interface{}) []string {
	cmds := getChangedValues(original, modified, "", nil)
	sort.Strings(cmds)
	return cmds
}

func getChangedValues(original map[string]interface{}, modified map[string]interface{}, prefix string, cmds []string) []string {
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
			cmds = append(cmds, getChangedValues(oVal, val, curKey, cmds)...)
		case []interface{}:
			if !reflect.DeepEqual(v, original[k]) {
				if len(val) == 0 {
					cmds = append(cmds, fmt.Sprintf("%s=null", curKey))
					continue
				}

				if isSimpleArray(val) {
					s, err := PrintArray(val)
					if err != nil {
						panic(err)
					}
					cmds = append(cmds, fmt.Sprintf("%s=%s", curKey, s))
					continue
				}

				for i, element := range val {
					em, ok := element.(map[string]interface{})
					if !ok {
						panic("element is not a map")
					}
					cmds = append(cmds, getChangedValues(map[string]interface{}{}, em, fmt.Sprintf("%s[%d]", curKey, i), cmds)...)
				}
			}
		case string:
			cmds = append(cmds, fmt.Sprintf("%s=%v", curKey, escapeValue(val)))
		default:
			if !reflect.DeepEqual(original[k], val) {
				data, err := json.Marshal(val)
				if err != nil {
					panic(err)
				}
				cmds = append(cmds, fmt.Sprintf("%s=%v", curKey, string(data)))
			}
		}
	}
	return cmds
}

// kubernetes.io/role becomes "kubernetes\.io/role"
func escapeKey(s string) string {
	if !strings.ContainsRune(s, '.') {
		return s
	}
	return `"` + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `.`, `\.`) + `"`
}

// "value1,value2" becomes value1\,value2
func escapeValue(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `,`, `\,`)
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
		case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, nil, json.Number:
			if i > 0 {
				buf.WriteString(", ")
			}
			err := json.NewEncoder(&buf).Encode(v)
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("[%d] holds a complex type %v", i, reflect.TypeOf(a[i]))
		}
	}
	buf.WriteRune('}')
	return buf.String(), nil
}
