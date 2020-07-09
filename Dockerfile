FROM golang:1.14.4
WORKDIR /app/goapp/
COPY app.go .
RUN go build -o main .

FROM python:3.7.7-slim-buster

RUN apt-get -qq update && apt-get -qqy install \
    apt-utils \
    alien \
    libaio1 \
    mongo-tools \
    && pip install --upgrade pip

COPY . /app

RUN cd /app \
    && ./install.sh --acceptlicenses --nousage --notestextras \
    && ln -s /root/.pipelinewise /app/.pipelinewise

RUN mkdir /app/config-data


COPY --from=0 /app/goapp/main .
EXPOSE 8080
CMD ["./main"]
