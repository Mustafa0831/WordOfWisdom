# WordOfWisdom



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
$ wow-client -addr go-wow-server:1024 # connects to the server, performs POW puzzle solving and receives a quote.
```

#### Arguments
```
Usage of wow-client:
  -addr string
         (default "0.0.0.0:1024")
  -count uint
        Consumers count (default 1)
  -print
        Print the output
  -print_err
        Print network errors
```