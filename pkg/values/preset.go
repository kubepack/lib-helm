package values

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/hub/resourceeditors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	chartsapi "x-helm.dev/apimachinery/apis/charts/v1alpha1"
)

func MergePresetValues(kc client.Client, chrt *chart.Chart, ref chartsapi.ChartPresetFlatRef) (map[string]interface{}, error) {
	var valOpts Options

	if ref.PresetName != "" {
		rid := &kmapi.ResourceID{
			Group: ref.Group,
			Name:  ref.Resource,
			Kind:  ref.Kind,
		}
		rid, err := kmapi.ExtractResourceID(kc.RESTMapper(), *rid)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to detect resource ID for %#v", *rid)
		}
		ed, ok := resourceeditors.LoadByGVR(kc, rid.GroupVersionResource())
		if !ok {
			return nil, errors.Errorf("failed to detect ResourceEditor for %#v", rid.GroupVersionResource())
		}

		for _, variant := range ed.Spec.Variants {
			if variant.Name != ref.PresetName {
				continue
			}
			if variant.Selector == nil {
				continue
			}

			sel, err := metav1.LabelSelectorAsSelector(variant.Selector)
			if err != nil {
				return nil, err
			}
			var list chartsapi.ClusterChartPresetList
			err = kc.List(context.TODO(), &list, client.MatchingLabelsSelector{Selector: sel})
			if err != nil {
				return nil, err
			}
			ccps := list.Items
			sort.Slice(ccps, func(i, j int) bool {
				return ccps[i].Name < ccps[j].Name
			})
			for _, ccp := range ccps {
				if ref.Namespace == "" {
					valOpts.ValueBytes = append(valOpts.ValueBytes, ccp.Spec.Values.Raw)
				} else {
					var cp chartsapi.ChartPreset
					err := kc.Get(context.TODO(), client.ObjectKey{Namespace: ref.Namespace, Name: ccp.Name}, &cp)
					if err == nil {
						valOpts.ValueBytes = append(valOpts.ValueBytes, cp.Spec.Values.Raw)
					} else if apierrors.IsNotFound(err) {
						valOpts.ValueBytes = append(valOpts.ValueBytes, ccp.Spec.Values.Raw)
					} else {
						return nil, err
					}
				}
			}
		}
	}
	return valOpts.MergeValues(chrt)
}
