FROM golang:1.22.5

WORKDIR /app

# Set Go proxy to a Chinese mirror
ENV GOPROXY=https://goproxy.cn

COPY . .

RUN go mod download

# Adjust the build command to point to the correct directory
RUN go build -o main.exe ./main

EXPOSE 8888

CMD ["./main.exe"]