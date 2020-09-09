FROM python:3.8-slim
ENV PYTHONUNBUFFERED True
WORKDIR /code
RUN pip install grpcio protobuf
COPY *.py ./
CMD [ "python", "greeter_server.py" ]
