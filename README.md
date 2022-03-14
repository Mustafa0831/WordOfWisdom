# WordOfWisdom
# About

Design and implement "Word of Wisdom" TCP server:

- TCP server should be protected from DDoS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from "word of wisdom" book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.


## Building
### Server

``` shell
$ make docker/server-image
```

### Client

``` shell
$ make docker/client-image
```

## Usage
### Server

``` shell
$ make start/server # starts server container
```

### Client

``` shell
$ make start/client # starts client container and opens its shell
$ client -addr go-server:1111 # connects to the server, performs POW puzzle solving and receives a quote.
```

#### Arguments
```
Usage of client:
  -addr string
         (default "0.0.0.0:1111")
  -count uint
        Consumers count (default 1)
  -print
        Print the output
  -print_err
        Print network errors
```