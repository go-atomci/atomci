package settings

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestTryLoginRegistry(t *testing.T) {
	type args struct {
		basicUrl string
		username string
		password string
		insecure bool
		authHead string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "basic auth pass", args: args{
			basicUrl: "www.unitest.com",
			username: "abc",
			password: "def",
			insecure: false,
			authHead: "Basic realm=Nexus Docker Registry",
		}, wantErr: false},
		{name: "basic upper auth pass", args: args{
			basicUrl: "www.unitest.com",
			username: "abc",
			password: "def",
			insecure: false,
			authHead: "BASIC realm=Nexus Docker Registry",
		}, wantErr: false},
		{name: "bearer auth pass", args: args{
			basicUrl: "www.unitest.com",
			username: "abc",
			password: "def",
			insecure: false,
			authHead: `Bearer realm=https://beaer.unitest.com,service="registry.unitest.com"`,
		}, wantErr: false},
		{name: "bearer complex auth pass", args: args{
			basicUrl: "www.unitest.com",
			username: "abc",
			password: "def",
			insecure: false,
			authHead: `Bearer realm=https://beaer.unitest.com/auth?tartget=https://abc.def.com,service="registry.unitest.com"`,
		}, wantErr: false},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		httpmock.RegisterResponder("GET", "https://"+tt.args.basicUrl, httpmock.NewStringResponder(200, "ok"))
		httpmock.RegisterResponder("GET", fmt.Sprintf("https://%s/v2/", tt.args.basicUrl),
			func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("Authorization") != "" {
					return httpmock.NewStringResponse(200, ""), nil
				}
				resp := httpmock.NewStringResponse(401, "")
				resp.Header.Add("Www-Authenticate", tt.args.authHead)
				return resp, nil
			},
		)
		httpmock.RegisterResponder("GET", `=~https://beaer\.unitest\.com.*`, httpmock.NewStringResponder(200, "ok"))
		t.Run(tt.name, func(t *testing.T) {
			if err := TryLoginRegistry(tt.args.basicUrl, tt.args.username, tt.args.password, tt.args.insecure); (err != nil) != tt.wantErr {
				t.Errorf("TryLoginRegistry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
