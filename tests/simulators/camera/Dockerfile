FROM python:3.8

RUN apt update
RUN apt install -y openssl

WORKDIR /app
COPY ./requirements.txt ./requirements.txt
RUN pip install -r requirements.txt

COPY ./device.py ./device.py
COPY ./media ./media

CMD ["python", "-u", "device.py" ]