# HTTP Decompressor

HTTP Decompressor is a Go library that wraps HTTP responses to transparently decompress response bodies based on the Content-Encoding header. It is designed to be used alongside http.Client, simplifying the handling and reading of compressed responses. The package supports multiple encoding formats including gzip, deflate, zlib, zstd, snappy, lz4 and brotli.

## Features

- Automatically detects and decompresses supported encoding formats.
- Returns the original stream if no or "identity" encoding is specified.
- Provides a custom error for unsupported encoding types.

## Supported Encodings

- gzip
- deflate
- zlib
- zstd
- snappy
- brotli

## Installation

Use go modules to install:

    go get github.com/fereidani/httpdecompressor

## Usage

Below is an example of how to use the package to decompress an HTTP response:

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "github.com/fereidani/httpdecompressor"
)

func main() {
    req, err := http.NewRequest("GET", "http://example.com", nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Accept-Encoding", httpdecompressor.ACCEPT_ENCODING)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // Decompress response body based on Content-Encoding
    reader, err := httpdecompressor.Reader(resp)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    // Read and print the decompressed content
    body, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(body))
}
```

### ReadAll Example

```go
// Retrieve and print decompressed HTTP response content.
resp, err := http.DefaultClient.Do(req)
if err != nil {
    // handle error
}
defer resp.Body.Close()
data, err := httpdecompressor.ReadAll(resp)
if err != nil {
    // handle error
}
fmt.Println(string(data))
```

### ReadIntoWriter Example

```go
// Write the decompressed response directly into a file on disk.
resp, err := http.DefaultClient.Do(req)
if err != nil {
    // handle error
}
defer resp.Body.Close()
file, err := os.OpenFile("image.jpg", os.O_CREATE|os.O_WRONLY, 0666)
if err != nil {
    // handle error
}
if err := httpdecompressor.ReadIntoWriter(resp, &file); err != nil {
    // handle error
}
```

### ReaderFromReader Example

```go
reader, err := httpdecompressor.ReaderFromReader(customBodyReader, "gzip")
if err != nil {
    // handle error
}
defer reader.Close()
data, err := io.ReadAll(reader)
if err != nil {
    // handle error
}
fmt.Println(string(data))
```

## License

This project is published with MIT license.

## Acknowledgments

This package utilizes the following libraries for decompression:

- [Andybalholm Brotli](https://github.com/andybalholm/brotli) for brotli
- [Klauspost Compress](https://github.com/klauspost/compress) for gzip, deflate, zlib, zstd, and snappy
- [Pierrec LZ4](github.com/pierrec/lz4) for lz4
