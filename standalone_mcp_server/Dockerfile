FROM golang:1.25

WORKDIR /app

# Copy the go.mod and go.sum files.
# Using a wildcard (go.*) allows this to work even if go.sum does not exist.
COPY go.* ./

# Download dependencies. Go will create or update go.sum as needed inside the container.
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the application.
RUN go build -o /mcp cmd/server/main.go

ENV LM_STUDIO_ENDPOINT="http://host.docker.internal:1234/v1/chat/completions"

EXPOSE 3333

CMD [ "/mcp" ]
