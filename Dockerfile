FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
EXPOSE 4040

ENV MODE=release

ENV GIN_MODE=release

ENV DB_HOST=postgres.blusnake.net
ENV DB_PORT=35432
ENV DB_USER=project-persona
ENV DB_PASS=jZFnGNY7yc6QYb2H
ENV DB_NAME=project-persona

COPY . .
RUN go build -v -o /usr/local/bin/app ./

CMD ["app"]