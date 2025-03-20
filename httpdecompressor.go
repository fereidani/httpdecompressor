package httpdecompressor

import (
	"io"
	"net/http"
	"os"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
)

const ACCEPT_ENCODING = "gzip, deflate, br, zstd, snappy, zlib, lz4"

type UnsupportedContentEncodingError struct {
	ContentEncoding string
}

func (e *UnsupportedContentEncodingError) Error() string {
	return "Unsupported content encoding: " + e.ContentEncoding
}

func ReadAll(response *http.Response) ([]byte, error) {
	reader, err := Reader(response)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func ReadIntoWriter(response *http.Response, writer io.Writer) error {
	reader, err := Reader(response)
	if err != nil {
		return err
	}
	defer reader.Close()
	var buf [1024 * 64]byte
	_, err = io.CopyBuffer(writer, reader, buf[:])
	return err
}

func ReadIntoFile(response *http.Response, filename string, perm os.FileMode) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	return ReadIntoWriter(response, file)
}

func ReaderFromReader(body io.ReadCloser, contentEncoding string) (io.ReadCloser, error) {
	switch contentEncoding {
	case "gzip":
		reader, err := gzip.NewReader(body)
		if err != nil {
			return nil, err
		}
		return reader, nil
	case "deflate":
		reader := flate.NewReader(body)
		return reader, nil
	case "zlib":
		reader, err := zlib.NewReader(body)
		if err != nil {
			return nil, err
		}
		return reader, nil
	case "zstd":
		reader, err := zstd.NewReader(body)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(reader), nil
	case "snappy":
		reader := snappy.NewReader(body)
		return io.NopCloser(reader), nil
	case "br":
		reader := brotli.NewReader(body)
		return io.NopCloser(reader), nil
	case "lz4":
		reader := lz4.NewReader(body)
		return io.NopCloser(reader), nil
	default:
		if contentEncoding != "" && contentEncoding != "identity" {
			return nil, &UnsupportedContentEncodingError{contentEncoding}
		}
		return body, nil
	}
}

func Reader(response *http.Response) (io.ReadCloser, error) {
	return ReaderFromReader(response.Body, response.Header.Get("Content-Encoding"))
}
