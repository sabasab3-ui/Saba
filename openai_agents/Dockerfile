FROM python:3.11-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app ./app
COPY .env.example ./.env.example
COPY run.sh ./run.sh

EXPOSE 8090
CMD ["python", "-m", "app.main"]
