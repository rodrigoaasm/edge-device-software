FROM python:3.8

WORKDIR /app
COPY ./device.py ./device.py
RUN pip install paho-mqtt

CMD ["python", "-u", "device.py" ]