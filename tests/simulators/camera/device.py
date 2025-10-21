#!/usr/bin/env python3
import argparse
import os
import signal
import time
import random
import json
import paho.mqtt.client as mqtt
import boto3
from datetime import datetime


def make_ack_payload(device_id, timestamp_ref):
    now = int(datetime.now().timestamp() * 1000000)
    return json.dumps({
        "deviceId": device_id,
        "correlation": timestamp_ref,
        "success": True,
        "message": "",
        "timestamp": now
    })

def make_payload(device_id, filepath ):
    now = int(datetime.now().timestamp() * 1000000)
    return json.dumps({
        "deviceId": device_id,
        "filename": filepath,
        "timestamp": now
    })


STOP_LOOP = False

def on_exit(sig, frame):  
    print ("Recebido sinal de encerramento")
    global STOP_LOOP
    STOP_LOOP = True

def make_on_connect(trigger_topic):
    def on_connect(client, userdata, flags, rc):
        if rc == 0:
            print("✅ MQTT conectado com sucesso")
            client.subscribe(trigger_topic)  # Assina o tópico de comando
            print(f"📡 Aguardando mensagens ${trigger_topic} em ...")
        else:
            print("❌ Falha na conexão MQTT, rc=", rc)

    return on_connect

def on_disconnect(client, userdata, rc):
    print("Desconectado do broker (rc=", rc, ")")

def make_on_message(trigger_topic, data_topic, result_topic, device_id, s3_client, s3_bucket):
    def on_message(client, userdata, msg):
        try:
            payload_str = msg.payload.decode()
            print(f"📥 Mensagem recebida em {msg.topic}: {payload_str}")
            data = json.loads(payload_str)
            timestamp_ref = data.get("timestamp", 0)
            if timestamp_ref == 0:
                print("❌ Mensagem sem timestamp")
                return
        
            if msg.topic == trigger_topic:                      
                payload = make_ack_payload(device_id, timestamp_ref)
                print(f"🚀 Publicando Ack em {result_topic}: {payload}")
                client.publish(result_topic, payload, qos=0, retain=0)            
        except Exception as e:
            print("⚠️ Erro ao processar mensagem:", e)

        try:
            print(f"🚀 Fazendo upload da imagem... ")
            now = int(datetime.now().timestamp() * 1000000)
            image_path = f"{device_id}-{now}.png"
            s3_client.upload_file("media/image_test.png", s3_bucket, image_path)
            print(f"✅ Imagem enviada para {image_path}. Publicando metadata em {data_topic}...")
            payload = make_payload(device_id, image_path)
            client.publish(data_topic, payload, qos=0, retain=0)
            print(f"✅ Metadata enviada para {data_topic}")
        except Exception as e:
            print(f"❌ Erro inesperado: {e}")

    return on_message

def simule(host, port, data_topic, config_topic, result_topic, device_id, interval, s3_endpoint, s3_access_key, s3_secret_key, s3_region, s3_bucket):
    s3_client = boto3.client(
        "s3",
        endpoint_url=s3_endpoint,
        aws_access_key_id=s3_access_key,
        aws_secret_access_key=s3_secret_key,
        region_name=s3_region 
    )

    trigger_topic = f"{config_topic}/{device_id}"
    client = mqtt.Client()
    client.on_connect = make_on_connect(trigger_topic)
    client.on_disconnect = on_disconnect
    client.on_message = make_on_message(
        trigger_topic,
        f"{data_topic}/{device_id}",
        f"{result_topic}/{device_id}",
        device_id,
        s3_client,
        s3_bucket,
    )

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
    default_mqtt_result_topic = os.environ.get("DEVICE_PY_RESULT_TOPIC", "device/results")
    default_mqtt_config_topic = os.environ.get("DEVICE_PY_CONFIG_TOPIC", "device/config")
    default_device_id = os.environ.get("DEVICE_PY_DEVICE_ID", "cc0001")
    default_interval = float(os.environ.get("DEVICE_PY_DATA_INTERVAL", 1000))
    default_s3_endpoint = os.environ.get("DEVICE_PY_AWS_S3_ENDPOINT")
    default_s3_access_key = os.environ.get("DEVICE_PY_AWS_ACCESS_KEY_ID")
    default_s3_secret_key = os.environ.get("DEVICE_PY_AWS_SECRET_KEY")
    default_s3_region = os.environ.get("DEVICE_PY_AWS_DEFAULT_REGION")
    default_s3_bucket = os.environ.get("DEVICE_PY_AWS_BUCKET_NAME", default_device_id)

    print(default_s3_bucket)

    parser = argparse.ArgumentParser(description="Simulador de dispositivo Python.")
    parser.add_argument("--host", type=str, default=default_mqtt_host, help="Host do servidor MQTT (padrão: localhost)")
    parser.add_argument("--port", type=int, default=default_mqtt_port, help="Porta do servidor MQTT (padrão: 1883)")
    parser.add_argument("--data_topic", type=str, default=default_mqtt_data_topic, help="Tópico para publicação de dados (padrão: device/data)")
    parser.add_argument("--result_topic", type=str, default=default_mqtt_result_topic, help="Tópico para publicação de resultados (padrão: device/data)")
    parser.add_argument("--config_topic", type=str, default=default_mqtt_config_topic, help="Tópico para recebimento de atuações (padrão: device/config)")
    parser.add_argument("--device_id", type=str, default=default_device_id, help="ID do dispositivo (padrão: device/data)")
    parser.add_argument("--interval", type=int, default=default_interval, help="Intervalo entre mensagens (ms) (padrão: 1000)")
    parser.add_argument("--s3_endpoint", type=str, default=default_s3_endpoint, help="Endpoint do S3")
    parser.add_argument("--s3_access_key", type=str, default=default_s3_access_key, help="Chave de acesso ao S3")
    parser.add_argument("--s3_secret_key", type=str, default=default_s3_secret_key, help="Secret de acesso ao S3")
    parser.add_argument("--s3_region", type=str, default=default_s3_region, help="Região do S3")
    parser.add_argument("--s3_bucket", type=str, default=default_s3_bucket, help="Bucket do S3")

    args = parser.parse_args()

    simule(
        args.host,
        args.port,
        args.data_topic,
        args.config_topic,
        args.result_topic,
        args.device_id,
        args.interval,
        args.s3_endpoint,
        args.s3_access_key,
        args.s3_secret_key,
        args.s3_region, 
        args.s3_bucket,
    )
