FROM python:3.8

RUN apt update
RUN apt install -y openssl

WORKDIR /app
COPY ./device.py ./device.py
COPY ./requirements.txt ./requirements.txt
COPY ./media ./media
RUN pip install -r requirements.txt

CMD ["python", "-u", "device.py" ]