FROM python:3.12.1

WORKDIR /app

COPY ./requirements.txt /app

RUN pip install --no-cache-dir -r requirements.txt

COPY . /app

EXPOSE 8080

# CMD ["gunicorn", "--bind=0.0.0.0", "--timeout", "600", "main:app"]
CMD ["python", "app.py"]