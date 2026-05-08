# Kubernetes 入口

`deployments/k8s` 是蓝图中的 Kubernetes 目录名。当前实际清单维护在 `deployments/kubernetes`，本目录通过 `kustomization.yaml` 引用同一组资源，避免两套 YAML 漂移。

```bash
kubectl apply -k deployments/k8s
```
