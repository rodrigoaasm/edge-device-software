import time
import argparse

def main():
    parser = argparse.ArgumentParser(description="Stress de memória simples em Python.")
    parser.add_argument("--blocos", type=int, default=20, help="Número de blocos a alocar (cada um com 100MB)")
    parser.add_argument("--tamanho", type=int, default=100, help="Tamanho de cada bloco em MB (padrão: 100MB)")
    parser.add_argument("--intervalo", type=float, default=1.0, help="Intervalo entre alocações em segundos")

    args = parser.parse_args()

    MB_POR_BLOCO = args.tamanho
    NUM_BLOCOS = args.blocos
    INTERVALO = args.intervalo

    blocos = []

    print(f"Iniciando stress de memória: {MB_POR_BLOCO}MB × {NUM_BLOCOS} blocos")

    try:
        for i in range(NUM_BLOCOS):
            bloco = bytearray(MB_POR_BLOCO * 1024 * 1024)
            bloco[:] = b'\xAA' * len(bloco)
            blocos.append(bloco)
            print(f"Bloco {i+1}/{NUM_BLOCOS} alocado — Total: {(i+1)*MB_POR_BLOCO} MB")
            time.sleep(INTERVALO)

        print("Alocação concluída. Pressione Ctrl+C para liberar a memória.")
        while True:
            time.sleep(10)

    except KeyboardInterrupt:
        print("\nEncerrando e liberando memória.")
        blocos.clear()

if __name__ == "__main__":
    main()
