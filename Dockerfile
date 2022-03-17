FROM golang:1.16

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /neat

CMD [ "/neat" ]
