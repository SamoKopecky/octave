FROM golang:1.18.3-bullseye

COPY . .

ENV GOPATH=""
ENV LAVALINK_HOST="lavalink_octave:2334"
ENV LAVALINK_PASSPHRASE=youshallnotpass

RUN go build

CMD ["./octave"]