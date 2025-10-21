FROM python:3.8

WORKDIR /app
COPY ./device.py ./device.py
COPY ./requirements.txt ./requirements.txt
RUN pip install -r requirements.txt

CMD ["python", "-u", "device.py" ]