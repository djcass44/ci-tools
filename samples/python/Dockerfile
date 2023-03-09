ARG BASE_VERSION
FROM harbor.dcas.dev/registry.gitlab.com/av1o/base-images/python:$BASE_VERSION

RUN mkdir -p /home/somebody/app
WORKDIR /home/somebody/app

COPY requirements.txt .
RUN pip install micropipenv && \
    micropipenv install

COPY . .

ENTRYPOINT ["python"]
CMD ["app.py"]