# Kubernetes Sequential Scheduler Plugin

A Kubernetes scheduler plugin that respects pod scheduling order based on creation timestamps (FIFO).

## Overview

This plugin implements the `QueueSort` extension point to ensure pods are scheduled in the order they were created. Older pods (by creation timestamp) are prioritized in the scheduling queue.

## Features

- FIFO scheduling based on pod creation timestamp
- Tie-breaking using pod UID when timestamps are equal
- Minimal overhead and simple implementation
- Compatible with Kubernetes 1.29+

## Building

Build the Docker image:

```bash
docker build -t sequential-scheduler:latest .
```

For development, you can build locally:

```bash
go mod download
go build -o scheduler .
```

## Deployment

1. Apply the RBAC configuration:

```bash
kubectl apply -f manifests/rbac.yaml
```

2. Apply the scheduler configuration:

```bash
kubectl apply -f manifests/scheduler-config.yaml
```

3. Deploy the scheduler:

```bash
kubectl apply -f manifests/deployment.yaml
```

4. Verify the scheduler is running:

```bash
kubectl get pods -n kube-system -l component=sequential-scheduler
```

## Using the Sequential Scheduler

To use this scheduler for your pods, specify the scheduler name in the pod spec:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  schedulerName: sequential-scheduler
  containers:
  - name: my-container
    image: nginx
```

## How It Works

The plugin implements the `QueueSort` interface, which determines the order in which pods are popped from the scheduling queue. The `Less` function compares two pods:

1. First by creation timestamp (older pods first)
2. If timestamps are equal, by UID (lexicographically)

This ensures deterministic, sequential scheduling behavior.

## Configuration

The scheduler configuration is in `manifests/scheduler-config.yaml`. The key configuration:

```yaml
profiles:
- schedulerName: sequential-scheduler
  plugins:
    queueSort:
      enabled:
      - name: SequentialScheduling
      disabled:
      - name: "*"
```

This disables the default queue sort plugin and enables the sequential scheduling plugin.

## Development

To modify the plugin behavior, edit `plugin.go`. The main logic is in the `Less` method:

```go
func (s *SequentialScheduling) Less(pInfo1, pInfo2 *framework.QueuedPodInfo) bool {
    // Returns true if pod1 should be scheduled before pod2
}
```

## Troubleshooting

Check scheduler logs:

```bash
kubectl logs -n kube-system -l component=sequential-scheduler
```

Verify scheduler configuration:

```bash
kubectl get configmap sequential-scheduler-config -n kube-system -o yaml
```

## License

This project is open source and available under the MIT License.
