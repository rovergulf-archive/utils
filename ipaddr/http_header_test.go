package ipaddr_test

import (
	"fmt"
	"github.com/rovergulf/utils"
	"github.com/rovergulf/utils/ipaddr"
	"net/http"
	"testing"
)

type testReq struct {
	r    *http.Request
	name string
	val  string
	want string
}

var xForwardedCases = []testReq{
	{
		name: ipaddr.XForwardedFor,
		want: "192.168.0.1",
		val:  "192.168.0.1",
	},
	{
		name: ipaddr.XForwardedFor,
		want: "192.168.0.2",
		val:  "192.168.0.2",
	},
	{
		name: "_" + ipaddr.XForwardedFor,
		val:  "192.168.0.3",
		want: "",
	},
	{
		name: "X-Forwarded-Fo",
		val:  "192.168.0.4",
		want: "",
	},
}

var cfConnectingCases = []testReq{
	{
		name: ipaddr.CFConnectingIp,
		want: "172.20.0.1",
		val:  "172.20.0.1",
	},
	{
		name: ipaddr.CFConnectingIp,
		want: "172.20.0.2",
		val:  "172.20.0.2",
	},
	{
		name: "_" + ipaddr.CFConnectingIp,
		val:  "172.20.0.3",
		want: "",
	},
	{
		name: "F-Connecting-IP",
		val:  "172.20.0.4",
		want: "",
	},
}

var cfRealCases = []testReq{
	{
		name: ipaddr.CFRealIp,
		want: "10.0.0.1",
		val:  "10.0.0.1",
	},
	{
		name: ipaddr.CFRealIp,
		want: "10.0.0.2",
		val:  "10.0.0.2",
	},
	{
		name: "_" + ipaddr.CFRealIp,
		val:  "10.0.0.3",
		want: "",
	},
	{
		name: "X-Forwarded-Fo",
		val:  "10.0.0.4",
		want: "",
	},
}

//func TestGetRequestIPAddress(t *testing.T) {
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := ipaddr.GetRequestIPAddress(tt.args.r); got != tt.want {
//				t.Errorf("GetRequestIPAddress() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func (t *testReq) prepareRequest(ind int) {
	t.r = new(http.Request)
	t.r.Header = make(map[string][]string)
	t.r.Header.Set(t.name, t.val)
	t.r.Header.Set("User-Agent", fmt.Sprintf("Test-Case_%d", ind))
	t.name = fmt.Sprintf("%s_%d", t.name, ind)
}

func TestHttpCloudflareConnectingIP(t *testing.T) {
	for i, tt := range cfConnectingCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareRequest(i)
			if got := ipaddr.HttpCloudflareConnectingIP(tt.r); got != tt.want {
				t.Errorf("[%s] HttpCloudflareConnectingIP() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHttpCloudflareRealIP(t *testing.T) {
	for i, tt := range cfRealCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareRequest(i)
			if got := ipaddr.HttpCloudflareRealIP(tt.r); got != tt.want {
				t.Errorf("HttpCloudflareRealIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpForwardedFor(t *testing.T) {
	for i, tt := range xForwardedCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareRequest(i)
			if got := ipaddr.HttpForwardedFor(tt.r); got != tt.want {
				t.Errorf("HttpRequestFootprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequestFootprint(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{r: &http.Request{
				Method: http.MethodGet,
				Header: make(map[string][]string),
			}},
			want: utils.GenerateHashFromString("32.150.20.0:TestHttpRequestFootprint"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.r.Header.Set("User-Agent", "TestHttpRequestFootprint")
			tt.args.r.Header.Set(ipaddr.CFConnectingIp, "32.150.20.0")
			if got := ipaddr.HttpRequestFootprint(tt.args.r); got != tt.want {
				t.Errorf("HttpRequestFootprint() = %v, want %v", got, tt.want)
			}
		})
	}
}
