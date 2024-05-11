FROM golang:latest AS buld 

COPY . .

RUN go build -o /server ./