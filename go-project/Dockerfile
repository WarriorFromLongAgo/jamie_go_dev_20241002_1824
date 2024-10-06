FROM golang:1.22.5-windowsservercore-ltsc2022

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main.exe .

EXPOSE 8888

CMD ["./main.exe"]