FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY batch/ ./batch
COPY item_meta.yaml .

RUN mkdir -p ./data/raw ./data/processed

CMD ["python", "batch/main.py"]