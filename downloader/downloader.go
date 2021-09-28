package downloader

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	DefaultThreadOpt  = 4
	DefaultTimeoutOpt = time.Second * 30
)

// HttpClient is a client that does http requests.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Saver stores downloaded files.
type Saver interface {
	Save(link string, r io.Reader) error
}

// Downloader downloads files.
type Downloader struct {
	client  HttpClient
	saver   Saver
	threads int
	dir     string
	timeout time.Duration
}

// NewDownloader create new instance of Downloader.
func NewDownloader(saver Saver, opts ...OptFunc) (*Downloader, error) {
	d := &Downloader{
		client:  http.DefaultClient,
		saver:   saver,
		threads: DefaultThreadOpt,
		timeout: DefaultTimeoutOpt,
	}

	for _, opt := range opts {
		opt(d)
	}

	return d, nil
}

// Download takes links, downloads and saves them in several streams.
func (d *Downloader) Download(ctx context.Context, links []string) {
	sem := make(chan struct{}, d.threads)

	var wg sync.WaitGroup
	for _, link := range links {
		select {
		case <-ctx.Done():
			log.Printf("context done: %s", ctx.Err())
			break
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(link string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			d.downloadLink(ctx, link)
		}(link)
	}

	wg.Wait()
}

func (d *Downloader) downloadLink(ctx context.Context, link string) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		log.Printf("error making request for link %s: %s", link, err)
		return
	}

	resp, err := d.client.Do(req)
	if err != nil {
		log.Printf("error downloading link %s: %s", link, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("response status code is not 200 OK: %d", resp.StatusCode)
		return
	}

	fileName := d.fileNameForLink(link)

	err = d.saver.Save(fileName, resp.Body)
	if err != nil {
		log.Printf("error saving file for link %s: %s", link, err)
	} else {
		log.Printf("successfuly downloaded file by link %s", link)
	}
}

func (d *Downloader) fileNameForLink(link string) string {
	u, _ := url.Parse(link)
	if u.Path == "" || u.Path == "/" || strings.HasSuffix(u.Path, "/") {
		return "index.html"
	}

	parts := strings.Split(u.Path, "/")

	return parts[len(parts)-1]
}
