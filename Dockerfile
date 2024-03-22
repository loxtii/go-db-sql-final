FROM golang:1.22.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go tracker.db ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /parcel_app

CMD [ "/parcel_app" ]