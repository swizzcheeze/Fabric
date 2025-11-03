FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /mcp cmd/mcp/main.go

EXPOSE 3333

CMD [ "/mcp" ]
