FROM python:3.8-slim
ENV PYTHONUNBUFFERED True
WORKDIR /code
RUN pip install flask
COPY server.py .
CMD [ "python", "server.py" ]
