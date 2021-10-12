package downloader

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDownloaderWithOptions(t *testing.T) {
	testClient := &http.Client{}

	type want struct {
		threads   *int
		timeout   *time.Duration
		userAgent *string
		client    HttpClient
	}

	tests := []struct {
		name string
		opts []OptFunc
		want want
	}{
		{
			name: "with timeout",
			opts: []OptFunc{
				WithTimeout(time.Second * 10),
			},
			want: want{timeout: durationPtr(time.Second * 10)},
		},
		{
			name: "with threads",
			opts: []OptFunc{
				WithThreads(10),
			},
			want: want{threads: intPtr(10)},
		},
		{
			name: "with http client",
			opts: []OptFunc{
				WithHttpClient(testClient),
			},
			want: want{client: testClient},
		},
		{
			name: "with user-agent",
			opts: []OptFunc{
				WithUserAgent("test agent"),
			},
			want: want{userAgent: stringPtr("test agent")},
		},
		{
			name: "with all options",
			opts: []OptFunc{
				WithTimeout(time.Second * 8),
				WithThreads(16),
				WithHttpClient(testClient),
				WithUserAgent("user agent"),
			},
			want: want{
				timeout:   durationPtr(time.Second * 8),
				threads:   intPtr(16),
				client:    testClient,
				userAgent: stringPtr("user agent"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewDownloader(nil, tt.opts...)
			if tt.want.threads != nil {
				assert.Equal(t, *tt.want.threads, got.threads)
			}
			if tt.want.timeout != nil {
				assert.Equal(t, *tt.want.timeout, got.timeout)
			}
			if tt.want.client != nil {
				assert.Same(t, tt.want.client, got.client)
			}
			if tt.want.userAgent != nil {
				assert.Equal(t, *tt.want.userAgent, got.userAgent)
			}
		})
	}
}

func TestNewDownloaderDefaultOptions(t *testing.T) {
	got, _ := NewDownloader(nil)
	assert.Same(t, http.DefaultClient, got.client)
	assert.Equal(t, DefaultThreadOpt, got.threads)
	assert.Equal(t, DefaultTimeoutOpt, got.timeout)
	assert.Equal(t, DefaultUserAgent(), got.userAgent)
}

func TestDownloader_Download(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientMock := NewMockHttpClient(ctrl)
	clientMock.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://example.com", req.URL.String())

		resp := &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("this is the answer")),
		}

		return resp, nil
	})
	saverMock := NewMockSaver(ctrl)
	saverMock.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(link string, r io.Reader) error {
		data, err := ioutil.ReadAll(r)
		require.NoError(t, err)

		assert.Equal(t, "index.html", link)
		assert.Equal(t, "this is the answer", string(data))

		return nil
	})

	tests := []struct {
		name   string
		client HttpClient
		saver  Saver
	}{
		{
			name:   "one simple test",
			client: clientMock,
			saver:  saverMock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDownloader(
				tt.saver,
				WithHttpClient(tt.client),
			)
			require.NoError(t, err)

			d.Download(context.Background(), []string{"https://example.com"})
		})
	}
}

func TestDownloader_Download_CheckFileName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		link     string
		fileName string
	}{
		{
			name:     "only domain, expects index.html",
			link:     "https://example.com",
			fileName: "index.html",
		},
		{
			name:     "only domain with a slash at the end, expects index.html",
			link:     "https://example.com/",
			fileName: "index.html",
		},
		{
			name:     "long path with a slash at the end, expects index.html",
			link:     "https://example.com/abracadabra/abra-abra-cadabra/",
			fileName: "index.html",
		},
		{
			name:     "link with filename",
			link:     "https://example.com/file.zip",
			fileName: "file.zip",
		},
		{
			name:     "link with filename and with directories in the path",
			link:     "https://example.com/one/two/file.zip",
			fileName: "file.zip",
		},
		{
			name:     "link with filename and with directories in the path",
			link:     "https://example.com/one/two/file.zip",
			fileName: "file.zip",
		},
		{
			name:     "url encoded",
			link:     "https://example.com/one/two/%D0%BA%D0%BE%D0%B7%D0%B0%20%D0%B5%D0%B3%D0%BE%D0%B7%D0%B0",
			fileName: "коза егоза",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewMockHttpClient(ctrl)
			clientMock.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				require.Equal(t, tt.link, req.URL.String())

				resp := &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("this is the answer")),
				}

				return resp, nil
			})
			saverMock := NewMockSaver(ctrl)
			saverMock.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(link string, r io.Reader) error {
				assert.Equal(t, tt.fileName, link)

				return nil
			})
			d, err := NewDownloader(
				saverMock,
				WithHttpClient(clientMock),
			)
			require.NoError(t, err)

			d.Download(context.Background(), []string{tt.link})
		})
	}
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
