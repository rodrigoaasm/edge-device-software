#!/bin/bash
IMAGE_DIR=".images"

if [ ! -d "$IMAGE_DIR" ]; then
  echo "❌ Erro: Diretório '$IMAGE_DIR' não encontrado."
  echo "Certifique-se de que o diretório existe e contém seus arquivos .tar."
  exit 1
fi

shopt -s nullglob
files=("$IMAGE_DIR"/*.tar)

# Verifica se nenhum arquivo .tar foi encontrado
if [ ${#files[@]} -eq 0 ]; then
    echo "ℹ️  Nenhum arquivo .tar foi encontrado em '$IMAGE_DIR'."
    exit 0
fi

echo "Iniciando importação de ${#files[@]} imagem(ns) de '$IMAGE_DIR' para o K3s..."
echo "--------------------------------------------------------"

# 4. Faz o loop e importa cada arquivo .tar
for tar_file in "${files[@]}"; do
  echo "🔧 Importando $tar_file..."
  
  # Comando para importar a imagem no K3s (containerd)
  k3s ctr images import "$tar_file"
  
  # Verifica o status do último comando
  if [ $? -eq 0 ]; then
    echo "✅ Importação de $tar_file concluída!"
  else
    echo "⚠️  Falha ao importar $tar_file."
  fi
  echo "--------------------------------------------------------"
done

echo "🎉 Processo finalizado. Todas as imagens foram processadas."