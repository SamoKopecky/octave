# octave

A discord music bot written in Go using Lavalink.

## How to run

For local development first run lavalink with 

```sh
docker compose up lavalink
```

Then setup the `.env` file by using the template in `.env.example`. Then run the discord app with 

```sh
go run .
```

## How to deploy

To deploy also setup the `.env` file using `.env.example`. Then run

```sh
docker compose up
```
