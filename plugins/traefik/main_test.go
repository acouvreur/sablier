package traefik

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"
)

func TestSablierMiddleware_ServeHTTP(t *testing.T) {
	type fields struct {
		Next   http.Handler
		Config *Config
	}
	type sablier struct {
		headers map[string]string
		body    string
	}
	tests := []struct {
		name     string
		fields   fields
		sablier  sablier
		expected string
		code     int
	}{
		{
			name: "sablier service is ready",
			sablier: sablier{
				headers: map[string]string{
					"X-Sablier-Session-Status": "ready",
				},
				body: "response from sablier",
			},
			fields: fields{
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					httptrace.ContextClientTrace(r.Context()).WroteHeaders()
					fmt.Fprint(w, "response from service")

				}),
				Config: &Config{
					SessionDuration: "1m",
					Dynamic:         &DynamicConfiguration{},
				},
			},
			expected: "response from service",
			code:     200,
		},
		{
			name: "sablier service is not ready",
			sablier: sablier{
				headers: map[string]string{
					"X-Sablier-Session-Status": "not-ready",
				},
				body: "response from sablier",
			},
			fields: fields{
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					httptrace.ContextClientTrace(r.Context()).WroteHeaders()
					fmt.Fprint(w, "response from service")
				}),
				Config: &Config{
					SessionDuration: "1m",
					Dynamic:         &DynamicConfiguration{},
				},
			},
			expected: "response from sablier",
			code:     200,
		},
		{
			name: "sablier service is ready but 503",
			sablier: sablier{
				headers: map[string]string{
					"X-Sablier-Session-Status": "ready",
				},
				body: "response from sablier",
			},
			fields: fields{
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}),
				Config: &Config{
					SessionDuration: "1m",
					Dynamic:         &DynamicConfiguration{},
				},
			},
			expected: "response from sablier",
			code:     200,
		},
		{
			name: "sablier service is ready blocking",
			sablier: sablier{
				headers: map[string]string{
					"X-Sablier-Session-Status": "ready",
				},
				body: "response from sablier",
			},
			fields: fields{
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					httptrace.ContextClientTrace(r.Context()).WroteHeaders()
					fmt.Fprint(w, "response from service")
				}),
				Config: &Config{
					SessionDuration: "1m",
					Blocking:        &BlockingConfiguration{},
				},
			},
			expected: "response from service",
			code:     200,
		},
		{
			name: "sablier service is not ready blocking",
			sablier: sablier{
				headers: map[string]string{
					"X-Sablier-Session-Status": "not-ready",
				},
				body: "response from sablier",
			},
			fields: fields{
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					httptrace.ContextClientTrace(r.Context()).WroteHeaders()
					fmt.Fprint(w, "response from service")
				}),
				Config: &Config{
					SessionDuration: "1m",
					Blocking:        &BlockingConfiguration{},
				},
			},
			expected: "response from sablier",
			// is this correct for blocking? I would expect to get error
			code: 200,
		},
		{
			name: "sablier service is ready blocking but 503",
			sablier: sablier{
				headers: map[string]string{
					"X-Sablier-Session-Status": "ready",
				},
				body: "response from sablier",
			},
			fields: fields{
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}),
				Config: &Config{
					SessionDuration: "1m",
					Blocking:        &BlockingConfiguration{},
				},
			},
			expected: "Found",
			code:     302,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sablierMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for key, value := range tt.sablier.headers {
					w.Header().Add(key, value)
				}
				w.Write([]byte(tt.sablier.body))
			}))
			defer sablierMockServer.Close()

			tt.fields.Config.SablierURL = sablierMockServer.URL

			sm, err := New(context.Background(), tt.fields.Next, tt.fields.Config, "middleware")
			if err != nil {
				panic(err)
			}

			req := httptest.NewRequest(http.MethodGet, "/my-nginx", nil)
			w := httptest.NewRecorder()

			sm.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("expected '%s' got '%v'", tt.expected, string(data))
			}
			if res.StatusCode != tt.code {
				t.Errorf("expected '%d' got '%d'", tt.code, res.StatusCode)

			}
		})
	}
}
