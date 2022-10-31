package traefik_test

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/acouvreur/sablier/plugins/traefik"
)

func TestConfig_BuildRequest(t *testing.T) {
	type fields struct {
		SablierURL      string
		Names           string
		SessionDuration string
		Dynamic         *traefik.DynamicConfiguration
		Blocking        *traefik.BlockingConfiguration
	}
	tests := []struct {
		name    string
		fields  fields
		want    *http.Request
		wantErr bool
	}{
		{
			name: "dynamic session with default values",
			fields: fields{
				SablierURL:      "http://sablier:10000",
				Names:           "nginx , apache",
				SessionDuration: "1m",
				Dynamic:         &traefik.DynamicConfiguration{},
			},
			want:    createRequest("GET", "http://sablier:10000/api/strategies/dynamic?display_name=sablier-middleware&names=nginx&names=apache&session_duration=1m", nil),
			wantErr: false,
		},
		{
			name: "dynamic session with theme values",
			fields: fields{
				SablierURL:      "http://sablier:10000",
				Names:           "nginx , apache",
				SessionDuration: "1m",
				Dynamic: &traefik.DynamicConfiguration{
					Theme: "hacker-terminal",
				},
			},
			want:    createRequest("GET", "http://sablier:10000/api/strategies/dynamic?display_name=sablier-middleware&names=nginx&names=apache&session_duration=1m&theme=hacker-terminal", nil),
			wantErr: false,
		},
		{
			name: "dynamic session with theme and display name values",
			fields: fields{
				SablierURL:      "http://sablier:10000",
				Names:           "nginx , apache",
				SessionDuration: "1m",
				Dynamic: &traefik.DynamicConfiguration{
					Theme:       "hacker-terminal",
					DisplayName: "Hello World!",
				},
			},
			want:    createRequest("GET", "http://sablier:10000/api/strategies/dynamic?display_name=Hello+World%21&names=nginx&names=apache&session_duration=1m&theme=hacker-terminal", nil),
			wantErr: false,
		},
		{
			name: "blocking session with default values",
			fields: fields{
				SablierURL:      "http://sablier:10000",
				Names:           "nginx , apache",
				SessionDuration: "1m",
				Blocking:        &traefik.BlockingConfiguration{},
			},
			want:    createRequest("GET", "http://sablier:10000/api/strategies/blocking?names=nginx&names=apache&session_duration=1m", nil),
			wantErr: false,
		},
		{
			name: "blocking session with timeout value",
			fields: fields{
				SablierURL:      "http://sablier:10000",
				Names:           "nginx , apache",
				SessionDuration: "1m",
				Blocking: &traefik.BlockingConfiguration{
					Timeout: "5m",
				},
			},
			want:    createRequest("GET", "http://sablier:10000/api/strategies/blocking?names=nginx&names=apache&session_duration=1m&timeout=5m", nil),
			wantErr: false,
		},
		{
			name: "both strategies defined",
			fields: fields{
				SablierURL:      "http://sablier:10000",
				Names:           "nginx , apache",
				SessionDuration: "1m",
				Dynamic:         &traefik.DynamicConfiguration{},
				Blocking:        &traefik.BlockingConfiguration{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &traefik.Config{
				SablierURL:      tt.fields.SablierURL,
				Names:           tt.fields.Names,
				SessionDuration: tt.fields.SessionDuration,
				Dynamic:         tt.fields.Dynamic,
				Blocking:        tt.fields.Blocking,
			}

			got, err := c.BuildRequest("sablier-middleware")
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.BuildRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.BuildRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createRequest(method string, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return request
}