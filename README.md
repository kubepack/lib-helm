[![Go Report Card](https://goreportcard.com/badge/kubepack.dev/lib-helm)](https://goreportcard.com/report/kubepack.dev/lib-helm)
[![Build Status](https://github.com/kubepack/lib-helm/workflows/CI/badge.svg)](https://github.com/kubepack/lib-helm/actions?workflow=CI)
[![codecov](https://codecov.io/gh/kubepack/lib-helm/branch/master/graph/badge.svg)](https://codecov.io/gh/kubepack/lib-helm)
[![Slack](http://slack.kubernetes.io/badge.svg)](http://slack.kubernetes.io/#kubepack)
[![Twitter](https://img.shields.io/twitter/follow/kubepack.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=Kubepack)

# lib-helm
A Multi-tenant server-side Helm Library

[Helm](https://github.com/helm/helm), the Kubernetes package manager, implementation makes assumption about running from a single workstation. This library adapts the code from Helm project to remove any such assumptions and makes it possible to download & render charts on a multi-tenant server.

## Acknowledgement
Majority of the source code in this project has been taken from Helm project as evident from the license header in source files.
