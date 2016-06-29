# Webmock-proxy - Record external API interactions, Provide fast and safe testing.
Webmock-proxy is "Man in the middle" proxy -- allow you to external API test without HTTP/S connection. Has two functions "Record" and "Mock".

## Installation

```bash
$ git clone git@github.com:wantedly/webmock-proxy.git
$ cd webmock-proxy
$ make deps
$ make
```

## Usage
### Record HTTP/S interactions
Run webmock-proxy "Record" mode.

```
$ WEBMOCK_PROXY_RECORD=true bin/webmock-proxy
```

All HTTP/S connection is recorded and output simple JSON file.
Connection is used transparently.

### Mock API server
Run webmock-proxy "Mock" mode.

```
$ bin/webmock-proxy
```

Webmock-proxy reply HTTP/S connection using cache file.
In case of not exist cache file, webmock-proxy return status code *418*.

### Test using webmock-proxy
Read sample code [Go](./example/go) and [Ruby](./example/ruby).

## License
This project is releases under the [MIT license](http://opensource.org/licenses/MIT).
