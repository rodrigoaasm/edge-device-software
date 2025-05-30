#EDGE-DEVICE-SOFTWARE

## Build das imagens

Nesse repositório você encontrará os arquivos necessários para construir as imagens do 
Edge-Device-Software e importá-las para dentro do cluster k3s. Para fazer build e importação de uma única imagem utilize:

```sh
./build.sh <nome do serviço> <versão>
```
Para realizar o build e importação de todas as imagens automaticamente, utilize:

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




