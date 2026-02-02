# EDGE-DEVICE-SOFTWARE

## Build das imagens

Nesse repositório você encontrará os arquivos necessários para construir as imagens do Edge-Device-Software e importá-las para dentro do cluster k3s. Para fazer build e importação de uma única imagem utilize:

```sh
./build.sh <nome do serviço> <versão>
```
Para realizar o build e a importação de todas as imagens dos serviços padrão do sistema automaticamente, utilize:

```sh
./build-all.sh <versão>
```

## Instalação e desinstalação do ed-software

Para instalar o ed-software no cluster k3s, utilize:

```sh
./up.sh
```

Para derrubar o cluster mantendo os dados do ed-software, utilize:

```sh
./down.sh
```

Para desinstalação completa do ed-software, utilize:

```sh 
./down.sh -v
```

## Comandos do sistema

A mensagem enviada a um dispositivo deve seguir a seguinte estrutura:

```ts
{
  deviceId: string, // Um identificador único para o dispositivo alvo.
  correlationId: string, // Um ID único para correlacionar a resposta a esta mensagem específica.
  commands: Array<{ // Uma lista de comandos a serem executados no dispositivo.
    correlationId: string, // Um ID único para correlacionar a resposta a este comando específico.
    command: string, // O tipo de comando a ser executado (ex: "deploy", "update", "undeploy").
    args: any // Argumentos específicos necessários para o comando.
  }>
}
```

Cada objeto dentro do array commands representa uma única instrução para o dispositivo, permitindo o envio de múltiplos comandos em uma única mensagem.

### Deploy or update:

Este comando é usado para implantar um novo serviço ou atualizar um serviço existente no dispositivo

```ts
{
  correlationId: string, // ID de correlação para esta operação de deploy/update.
  command: "deploy" | "update", // Indica se é uma implantação inicial ("deploy") ou uma atualização ("update").
  args: {
    name: string, // O nome do serviço a ser implantado ou atualizado.
    image: string, // A imagem do contêiner a ser usada para o serviço (ex: "my-app:1.0.0").
    env?: Map<string,string>, // Um mapa de variáveis de ambiente a serem configuradas para o serviço.
    priorityProfile?: string | number, // O nome de profile padrão ou número para a criação de um profile especifico (padrão: single-service).
    port?: number, // A porta principal onde o serviço será exposto e acessível.
    internalPort?: number, // A porta usada para comunicação interna dentro do cluster.
    externalPort?: number // A porta que as chamadas externas ao dispositivo devem usar para acessar o serviço.
    requestMemory?: number // Quantidade memória necessária em MB (padrão: 128).
    limitMemory?: number // Limite de memória em MB (Padrão: 256).
    requestCPU?: number // Quantidade de núcleos de CPU necessário em Milli-núcleos (Padrão: 200).
    limitCPU?: number // Limite de núcleos de CPU em Milli-núcleos (Padrão: 500).
  }
}
```
Ao usar deploy, o sistema criará e iniciará um novo serviço com o name e image especificados, aplicando as env variáveis de ambiente. Se o comando for update, o serviço existente com o name será atualizado com a nova image e env configurações, geralmente resultando em uma reinicialização do serviço.

O campo opcional priorityProfile permite definir o perfil de prioridade de execução do serviço. Se um nome de perfil padrão ("iot-agent" ou "single-service") for fornecido, o serviço utilizará as configurações predefinidas para esse perfil. Alternativamente, um número entre 2 e 1000 pode ser enviado para criar um perfil de prioridade específico para este serviço. Alguns serviços tem seus nome reservados ("operator", "telegraf-agg" e "telegraf-bridge") por serem padrão do sistema, para esses serviços o campo priorityProfile não é aplicável.

As propriedades port, internalPort e externalPort definem como o serviço será acessível: port é a principal porta de comunicação; internalPort é para chamadas de outros serviços dentro do cluster; e externalPort é para requisições que vêm de fora do dispositivo.

### Undeploy Command:

Este comando é usado para remover um serviço existente de um dispositivo.

```ts
{
  correlationId: string, // ID de correlação para esta operação de undeploy.
  command: "undeploy", // O tipo de comando, sempre "undeploy" para esta operação.
  args: {
    name: string // O nome do serviço a ser removido do dispositivo.
    image: string // A imagem do contêiner associada ao serviço a ser removido.
  }
}
```
Ao receber este comando, o dispositivo encerrará e removerá o serviço identificado pelo name e image correspondentes, liberando os recursos que ele estava utilizando.

Ao adicionar novos serviços é possível determinar comandos especifico para eles.

## Comandos Básicos do Kubenetes

```sh
## Consumo de recurso do Cluster
kubectl top nodes

## Consumo de recurso dos Pods
kubectl top pods -A
```

## Ambiente de Preempção

Para facilitar os testes de preempção por memoria o ideal é manter um quantidade de memória para o S.O e para os sistemas básicos do Kubernetes para garantir o funcionamento básico do cluster. Para reserva memoria o arquivo `/etc/systemd/system/k3s.service` deve ser modificado para que o executável do k3s receba os argumentos abaixo.

```sh
--kubelet-arg="system-reserved=memory=1Gi" \
--kubelet-arg="kube-reserved=memory=1Gi"
```