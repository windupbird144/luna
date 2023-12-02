FROM docker.io/library/golang
WORKDIR /app
COPY . /app
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go build -o luna