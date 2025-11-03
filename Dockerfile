FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /mcp cmd/mcp/main.go

ENV LM_STUDIO_ENDPOINT="http://localhost:1234/v1/chat/completions"

EXPOSE 3333

CMD [ "/mcp" ]
