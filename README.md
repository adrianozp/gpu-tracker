# gpu-tracker
K8s operator for monitoring nodes

This operators lists all the nodes that matches the label from the labelConfig helm chart values. Defaults to: node-type: gpu-node.

# Setting up
Build and push the docker image:

`docker build --push -t <your-registry>/gpu-node-operator:<your-version> .`

Update the charts/gpu-tracker/values.yaml file with your image:
```yaml
image:
  repository: <your-registry>/gpu-node-operator
  tag: <your-version>
```

Install the helm chart:
```
helm install gpu-tracker ./charts/gpu-tracker
```

# Testing

Create a CRD, you can use the example:

```sh
kubectl apply -f example/gpu-instance.yaml
```

Check the GPU nodes listed:

```sh
kubectl get gputrackers.suse.tests.dev suse-gpu-tracker -o jsonpath='{.gpu_nodes}
```

You can also set the label on a node by using:
```sh
kubectl label nodes <node-name> node-type=gpu-node --overwrite
```