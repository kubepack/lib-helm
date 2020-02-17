package helm

import (
	"fmt"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func ListReleases(getter genericclioptions.RESTClientGetter, namespace string) ([]*release.Release, error) {
	cfg := new(action.Configuration)
	if err := cfg.Init(getter, namespace, "", debug); err != nil {
		return nil, err
	}

	client := action.NewList(cfg)

	return client.Run()
}

func UninstallRelease(getter genericclioptions.RESTClientGetter, name string, namespace string) error {
	cfg := new(action.Configuration)
	if err := cfg.Init(getter, namespace, "", debug); err != nil {
		return err
	}
	_, err := cfg.GetCapabilities()
	if err != nil {
		return err
	}

	client := action.NewUninstall(cfg)
	_, err = client.Run(name)
	return err
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	_ = log.Output(2, fmt.Sprintf(format, v...))
}
