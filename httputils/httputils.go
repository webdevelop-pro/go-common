package httputils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Request struct {
	Host, Port, Scheme string
	Method, Path       string
	Body               []byte
	Headers            map[string]string
}

// CreateDefaultRequest default json request
func CreateDefaultRequest(ctx context.Context, req Request) (*http.Request, error) {
	if req.Port != "" {
		req.Host = net.JoinHostPort(req.Host, req.Port)
	}
	if req.Scheme == "" {
		req.Scheme = "http"
	}

	res, err := http.NewRequestWithContext(
		ctx,
		req.Method,
		fmt.Sprintf("%s://%s%s", req.Scheme, req.Host, req.Path),
		bytes.NewBuffer((req.Body)),
	)
	if err != nil {
		return res, errors.Wrapf(err, "cannot create new request")
	}

	// if no content type set
	if ok := res.Header.Get("Content-Type"); ok == "" {
		res.Header.Add("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		res.Header.Add(key, value)
	}

	return res, nil
}

func CreateRequestWithFiles(req Request, body map[string]any, files map[string]string) (*http.Request, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	if req.Port != "" {
		req.Host = net.JoinHostPort(req.Host, req.Port)
	}

	if req.Body != nil {
		return nil, errors.New("req body should be empty, use body parameter")
	}

	values := map[string]io.Reader{}
	for k, v := range body {
		if vs, ok := v.(string); ok {
			values[k] = strings.NewReader(vs)
		}
	}

	for k, v := range files {
		f, err := os.Open(v)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot open file %s", v)
		}
		values[k] = f
	}

	for key, r := range values {
		var fw io.Writer
		var err error
		x, ok := r.(io.Closer)
		if !ok {
			continue
		}
		// upload a file
		if _, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, files[key]); err != nil {
				w.Close()
				x.Close()
				return nil, errors.Wrapf(err, "cannot CreateFormFile %s", key)
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				w.Close()
				x.Close()
				return nil, errors.Wrapf(err, "cannot CreateFormField %s", key)
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return nil, errors.Wrapf(err, "cannot io.Copy %s", key)
		}

		x.Close()
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	res, err := http.NewRequest(
		req.Method,
		fmt.Sprintf("%s://%s%s", req.Scheme, req.Host, req.Path),
		buf,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create new request")
	}

	res = res.WithContext(context.Background())

	// Don't forget to set the content type, this will contain the boundary.
	res.Header.Set("Content-Type", w.FormDataContentType())
	// Set up content length
	res.Header.Set("Content-Length", strconv.FormatInt(res.ContentLength, 10))

	return res, nil
}

func SendRequest(req *http.Request) ([]byte, *http.Header, int, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, resp.StatusCode, errors.Wrapf(err, "cannot read response body")
	}

	return bodyBytes, &resp.Header, resp.StatusCode, nil
}
