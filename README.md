# EDGE-DEVICE-SOFTWARE

## Dependencies

- k3s (kubernetes)
- docker

## Building Images

In this repository, you will find the necessary files to build the Edge-Device-Software images and import them into the k3s cluster. To build and import a single image, use:

>Note: `sudo` privileges might be required.

```sh
./build.sh <service-name> <version>
```

To automatically build and import all standard system service images, use:

>Note: `sudo` privileges might be required.

```sh
./build-all.sh <version>
```

---

## Installing and Uninstalling ed-software

To install ed-software in the k3s cluster, run:

```sh
./up.sh
```

To shut down the cluster while preserving ed-software data, use:

```sh
./down.sh
```

For a complete removal of ed-software, including all data, run:

```sh
./down.sh -v
```

---

## System Commands

Messages sent to a device must follow the structure below:

```ts
{
  deviceId: string, // Unique identifier of the target device.
  correlationId: string, // Unique ID to correlate the response to this specific message.
  commands: Array<{ // List of commands to be executed on the device.
    correlationId: string, // Unique ID to correlate the response to this specific command.
    command: string, // Type of command (e.g., "deploy", "update", "undeploy").
    args: any // Command-specific arguments.
  }>
}
```

Each object inside the `commands` array represents a single instruction for the device, allowing multiple commands to be sent in a single message.

---

### Deploy or Update

This command is used to deploy a new service or update an existing service on the device.

```ts
{
  correlationId: string, // Correlation ID for this deploy/update operation.
  command: "deploy" | "update", // Indicates initial deployment ("deploy") or update ("update").
  args: {
    name: string, // Name of the service to be deployed or updated.
    image: string, // Container image to be used (e.g., "my-app:1.0.0").
    env?: Map<string,string>, // Map of environment variables.
    priorityProfile?: string | number, // Default profile name or numeric value (default: single-service).
    port?: number, // Main exposed service port.
    internalPort?: number, // Port for internal cluster communication.
    externalPort?: number // Port for external access.
    requestMemory?: number // Requested memory in MB (default: 128).
    limitMemory?: number // Memory limit in MB (default: 256).
    requestCPU?: number // Requested CPU in millicores (default: 200).
    limitCPU?: number // CPU limit in millicores (default: 500).
  }
}
```

When using `deploy`, the system creates and starts a new service with the specified `name` and `image`, applying the environment variables. When using `update`, the existing service is updated with the new `image` and `env` settings, usually resulting in a service restart.

The optional `priorityProfile` field defines the execution priority profile. If a default profile name (`"iot-agent"` or `"single-service"`) is provided, predefined settings are applied. Alternatively, a value between 2 and 1000 can be used to create a custom priority profile.

Some services have reserved names (`"operator"`, `"telegraf-agg"`, and `"telegraf-bridge"`) because they are system defaults. For these services, the `priorityProfile` field is not applicable.

The `port`, `internalPort`, and `externalPort` properties define how the service is exposed:

* `port`: main communication port;
* `internalPort`: for inter-service communication inside the cluster;
* `externalPort`: for requests coming from outside the device.

---

### Undeploy Command

This command is used to remove an existing service from a device.

```ts
{
  correlationId: string, // Correlation ID for this undeploy operation.
  command: "undeploy", // Command type, always "undeploy".
  args: {
    name: string, // Name of the service to be removed.
    image: string // Container image associated with the service.
  }
}
```

When this command is received, the device stops and removes the service identified by `name` and `image`, freeing the resources it was using.

When adding new services, it is possible to define specific commands for them.

---

## Basic Kubernetes Commands

```sh
## Cluster resource usage
kubectl top nodes

## Pod resource usage
kubectl top pods -A
```

---

## Preemption Environment

To facilitate memory preemption testing, it is recommended to reserve part of the system memory for the operating system and basic Kubernetes services, ensuring minimum cluster stability.

To reserve memory, edit the file `/etc/systemd/system/k3s.service` and add the following arguments to the k3s executable:

```sh
--kubelet-arg="system-reserved=memory=1Gi" \
--kubelet-arg="kube-reserved=memory=1Gi"
```
