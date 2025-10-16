#!/usr/bin/env python3
import argparse
import os
import signal
import time
import random
import json
import paho.mqtt.client as mqtt
from datetime import datetime

TRIGGER_TOPIC = None
DATA_TOPIC = None
DEVICE_ID = None

def make_payload(device_id, timestamp_ref):
    now = int(datetime.now().timestamp() * 1000000)
    return json.dumps({
        "deviceId": device_id,
        "correlation": timestamp_ref,
        "success": True,
        "message": "",
        "timestamp": now
    })

STOP_LOOP = False

def on_exit(sig, frame):  
    print ("Recebido sinal de encerramento")
    global STOP_LOOP
    STOP_LOOP = True

def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("✅ MQTT conectado com sucesso")
        client.subscribe(TRIGGER_TOPIC)  # Assina o tópico de comando
        print(f"📡 Aguardando mensagens ${TRIGGER_TOPIC} em ...")
    else:
        print("❌ Falha na conexão MQTT, rc=", rc)

def on_disconnect(client, userdata, rc):
    print("Desconectado do broker (rc=", rc, ")")

def on_message(client, userdata, msg):
    global TRIGGER_TOPIC
    global DATA_TOPIC
    global DEVICE_ID
    try:
        payload_str = msg.payload.decode()
        print(f"📥 Mensagem recebida em {msg.topic}: {payload_str}")
        data = json.loads(payload_str)
        timestamp_ref = data.get("timestamp", 0)
        if timestamp_ref == 0:
            print("❌ Mensagem sem timestamp")
            return
       
        if msg.topic == TRIGGER_TOPIC:                      
            payload = make_payload(DEVICE_ID, timestamp_ref)
            print(f"🚀 Publicando em {DATA_TOPIC}: {payload}")
            client.publish(DATA_TOPIC, payload, qos=0, retain=0)
    except Exception as e:
        print("⚠️ Erro ao processar mensagem:", e)

def simule(host, port, data_topic, config_topic, device_id, interval):
    global TRIGGER_TOPIC
    global DATA_TOPIC
    global DEVICE_ID
    TRIGGER_TOPIC = f"{config_topic}/{device_id}"
    DATA_TOPIC = f"{data_topic}/{device_id}"
    DEVICE_ID= device_id
    
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect
    client.on_message = on_message

    try:
        print(f"📡 Conectando ao broker MQTT em {host}:{port}...")
        client.connect(host, port, keepalive=60)
        print(f"📡 Conexão estabelecida com o broker MQTT em {host}:{port}...")
    except Exception as e:
        print("Erro ao conectar no broker:", e)
        return
    
    signal.signal(signal.SIGINT, on_exit)
    signal.signal(signal.SIGTERM, on_exit)
    
    # Loop do cliente MQTT (escuta as mensagens)
    while not STOP_LOOP:
        client.loop(timeout=1.0)

    client.disconnect()

if __name__ == "__main__":
    default_mqtt_host = os.environ.get("DEVICE_PY_BROKER_HOST", "localhost")
    default_mqtt_port = int(os.environ.get("DEVICE_PY_BROKER_PORT", 31883))
    default_mqtt_data_topic = os.environ.get("DEVICE_PY_DATA_TOPIC", "device/data")
    default_mqtt_config_topic = os.environ.get("DEVICE_PY_CONFIG_TOPIC", "device/config")
    default_device_id = os.environ.get("DEVICE_PY_DEVICE_ID", "cc0001")
    default_interval = float(os.environ.get("DEVICE_PY_DATA_INTERVAL", 1000))
    
    parser = argparse.ArgumentParser(description="Simulador de dispositivo Python.")
    parser.add_argument("--host", type=str, default=default_mqtt_host, help="Host do servidor MQTT (padrão: localhost)")
    parser.add_argument("--port", type=int, default=default_mqtt_port, help="Porta do servidor MQTT (padrão: 1883)")
    parser.add_argument("--data_topic", type=str, default=default_mqtt_data_topic, help="Tópico para publicação de dados (padrão: device/data)")
    parser.add_argument("--config_topic", type=str, default=default_mqtt_config_topic, help="Tópico para recebimento de atuações (padrão: device/config)")
    parser.add_argument("--device_id", type=str, default=default_device_id, help="ID do dispositivo (padrão: device/data)")
    parser.add_argument("--interval", type=int, default=default_interval, help="Intervalo entre mensagens (ms) (padrão: 1000)")

    args = parser.parse_args()

    simule(args.host, args.port, args.data_topic, args.config_topic, args.device_id, args.interval)
