import time
import argparse
import os
import signal

STOP_LOOP = False

def on_exit(sig, frame):  
    print ("Recebido sinal de encerramento")
    global STOP_LOOP
    STOP_LOOP = True

def main():
    default_blocos = int(os.environ.get("ALOC_MEMO_BLOCK", 20))
    default_tamanho = int(os.environ.get("ALOC_MEMO_SIZE", 100))
    default_intervalo = float(os.environ.get("ALOC_MEMO_INTERVAl", 1.0))

    parser = argparse.ArgumentParser(description="Stress de memória simples em Python.")
    parser.add_argument("--blocos", type=int, default=default_blocos, help="Número de blocos a alocar (cada um com 100MB)")
    parser.add_argument("--tamanho", type=int, default=default_tamanho, help="Tamanho de cada bloco em MB (padrão: 100MB)")
    parser.add_argument("--intervalo", type=float, default=default_intervalo, help="Intervalo entre alocações em segundos")

    args = parser.parse_args()

    MB_POR_BLOCO = args.tamanho
    NUM_BLOCOS = args.blocos
    INTERVALO = args.intervalo

    blocos = []

    signal.signal(signal.SIGINT, on_exit)
    signal.signal(signal.SIGTERM, on_exit)

    print(f"Iniciando stress de memória: {MB_POR_BLOCO}MB × {NUM_BLOCOS} blocos")
    for i in range(NUM_BLOCOS):
        bloco = bytearray(MB_POR_BLOCO * 1024 * 1024)
        bloco[:] = b'\xAA' * len(bloco)
        blocos.append(bloco)
        print(f"Bloco {i+1}/{NUM_BLOCOS} alocado — Total: {(i+1)*MB_POR_BLOCO} MB")
        if STOP_LOOP:
            break
        time.sleep(INTERVALO)

    while not STOP_LOOP:
        time.sleep(10)

if __name__ == "__main__":
    main()
