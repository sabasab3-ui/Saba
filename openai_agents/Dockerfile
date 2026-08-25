FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app ./app
COPY tests ./tests
COPY .env.example .
COPY run.sh .
RUN chmod +x run.sh

EXPOSE 8090

CMD ["./run.sh"]
