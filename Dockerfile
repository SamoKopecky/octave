FROM golang:1.21.1 as build

WORKDIR /usr/src

COPY . .

ENV GOPATH=""

RUN go build -v

FROM golang:1.21.1 as prod

WORKDIR /usr/app

COPY .env .
COPY --from=build /usr/src/octave .

ENTRYPOINT ["./octave"]