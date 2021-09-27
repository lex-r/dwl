package downloader

import "time"

// OptFunc is a function for configuring Downloader.
type OptFunc func(d *Downloader)

// WithHttpClient sets given HttpClient in Downloader.
func WithHttpClient(c HttpClient) OptFunc {
	return func(d *Downloader) {
		d.client = c
	}
}

// WithTimeout sets timeout in Downloader.
func WithTimeout(t time.Duration) OptFunc {
	return func(d *Downloader) {
		d.timeout = t
	}
}

// WithThreads sets threads in Downloader.
func WithThreads(t int) OptFunc {
	return func(d *Downloader) {
		d.threads = t
	}
}
