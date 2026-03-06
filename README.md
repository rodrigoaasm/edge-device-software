# EDGE-DEVICE-SOFTWARE

## Dependencies

The following tools must be installed on the device:

* **k3s** (Kubernetes distribution)
* **Docker**
* **openssl**

---

## Building Images

This repository contains all required files to build the edge-device-doftware container images and import them into the k3s cluster.

### Build a Single Image

> Note: `sudo` privileges may be required.

```sh
./build.sh <service-name> <version>
```

* `<service-name>`: Name of the service directory.
* `<version>`: Image tag version (e.g., `1.0.0`).

### Build All Standard System Images

> Note: `sudo` privileges may be required.

```sh
./build-all.sh <version>
```

This command builds and imports all default system service images using the specified version.

---

## Installing and Uninstalling

TLS secrets are generated for internal communication. For this process, two OpenSSL configuration files are required:
- ca.cnf
- cert.cnf

These files must be placed inside the `.certs` directory:

**CA Example:**

```cnf
[ req ]
prompt = no
distinguished_name = req_distinguished_name
x509_extensions = v3_ca

[ req_distinguished_name ]
# Valores do Subject (DN)
C = BR
ST = MG
L = Itajuba
O = UNIFEI
OU = ED-SYSTEM
CN = ED-SYSTEM-CA
emailAddress = rodrigoasmaia@gmail.com

[ v3_ca ]
# Configurações para um Certificado de Autoridade (CA)
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints = CA:true
keyUsage = critical, cRLSign, keyCertSign
```` 

**Cert Example:**

```cnf
prompt = no
distinguished_name = req_distinguished_name
req_extensions = v3_req

[ req_distinguished_name ]
C = BR
ST = MG
L = Itajuba
O = UNIFEI
OU = ED-SYSTEM
CN = ED-SYSTEM
emailAddress = rodrigoasmaia@gmail.com

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = <dns>
IP.1 = <ip>
```

### Install

To deploy ed-software into the k3s cluster:

```sh
./up.sh
```

### Stop (Preserve Data)

To stop the cluster while preserving all persistent data:

```sh
./down.sh
```

### Complete Removal (Including Data)

To completely remove ed-software and all associated data:

```sh
./down.sh -v
```

Use this option carefully, as it deletes persistent volumes.

---

# System Commands

Messages sent to a device must follow the structure below:

```ts
{
  deviceId: string,          // Unique identifier of the target device
  correlationId: string,     // Correlates the response to this message
  commands: Array<{
    correlationId: string,   // Correlates the response to this specific command
    command: string,         // Command type ("deploy", "update", "undeploy")
    args: any                // Command-specific arguments
  }>
}
```

Each item inside `commands` represents a single instruction. Multiple commands can be sent in a single message.

---

## Deploy or Update

This command deploys a new service or updates an existing one.

```ts
{
  correlationId: string,
  command: "deploy" | "update",
  args: {
    name: string,                  // Service name
    image: string,                 // Container image (e.g., "my-app:1.0.0")
    env?: Map<string, string>,     // Environment variables
    priorityProfile?: string | number,
    port?: number,                 // Main service port
    internalPort?: number,         // Cluster internal communication port
    externalPort?: number,         // External access port
    requestMemory?: number,        // Memory request (MB, default: 128)
    limitMemory?: number,          // Memory limit (MB, default: 256)
    requestCPU?: number,           // CPU request (millicores, default: 200)
    limitCPU?: number              // CPU limit (millicores, default: 500)
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

## Undeploy

Removes an existing service from the device.

```ts
{
  correlationId: string,
  command: "undeploy",
  args: {
    name: string,
    image: string
  }
}
```

When this command is received, the device stops and removes the service identified by `name` and `image`, freeing the resources it was using.

When adding new services, it is possible to define specific commands for them.

---

# Basic Kubernetes Commands

### Cluster Resource Usage

```sh
kubectl top nodes
```

### Pod Resource Usage (All Namespaces)

```sh
kubectl top pods -A
```

---

# Preemption Environment

For memory preemption testing, it is recommended to reserve part of the system memory for the OS and essential Kubernetes components to ensure cluster stability.

Edit:

```
/etc/systemd/system/k3s.service
```

Add the following arguments to the k3s executable:

```sh
--kubelet-arg="system-reserved=memory=1Gi" \
--kubelet-arg="kube-reserved=memory=1Gi"
```

After modifying the file:

```sh
sudo systemctl daemon-reload
sudo systemctl restart k3s
```