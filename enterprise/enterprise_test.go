package enterprise_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/influxdata/chronograf"
	"github.com/influxdata/chronograf/enterprise"
	"github.com/influxdata/chronograf/log"
)

func Test_Enterprise_FetchesDataNodes(t *testing.T) {
	t.Parallel()

	ctrl := &ControlClient{
		Cluster: &enterprise.Cluster{},
	}
	cl := &enterprise.Client{
		Ctrl: ctrl,
	}

	bg := context.Background()
	err := cl.Connect(bg, &chronograf.Source{})

	if err != nil {
		t.Fatal("Unexpected error while creating enterprise client. err:", err)
	}

	if ctrl.ShowClustersCalled != true {
		t.Fatal("Expected request to meta node but none was issued")
	}
}

func Test_Enterprise_IssuesQueries(t *testing.T) {
	t.Parallel()

	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		called = true
		if r.URL.Path != "/query" {
			t.Fatal("Expected request to '/query' but was", r.URL.Path)
		}
		rw.Write([]byte(`{}`))
	}))
	defer ts.Close()

	cl := &enterprise.Client{
		Ctrl:   NewMockControlClient(ts.URL),
		Logger: log.New(log.DebugLevel),
	}

	err := cl.Connect(context.Background(), &chronograf.Source{})
	if err != nil {
		t.Fatal("Unexpected error initializing client: err:", err)
	}

	_, err = cl.Query(context.Background(), chronograf.Query{Command: "show shards", DB: "_internal", RP: "autogen"})

	if err != nil {
		t.Fatal("Unexpected error while querying data node: err:", err)
	}

	if called == false {
		t.Fatal("Expected request to data node but none was received")
	}
}

func Test_Enterprise_AdvancesDataNodes(t *testing.T) {
	m1 := NewMockTimeSeries("http://host-1.example.com:8086")
	m2 := NewMockTimeSeries("http://host-2.example.com:8086")
	cl := enterprise.NewClientWithTimeSeries(log.New(log.DebugLevel), chronograf.TimeSeries(m1), chronograf.TimeSeries(m2))

	err := cl.Connect(context.Background(), &chronograf.Source{})
	if err != nil {
		t.Error("Unexpected error while initializing client: err:", err)
	}

	_, err = cl.Query(context.Background(), chronograf.Query{Command: "show shards", DB: "_internal", RP: "autogen"})
	if err != nil {
		t.Fatal("Unexpected error while issuing query: err:", err)
	}

	_, err = cl.Query(context.Background(), chronograf.Query{Command: "show shards", DB: "_internal", RP: "autogen"})
	if err != nil {
		t.Fatal("Unexpected error while issuing query: err:", err)
	}

	if m1.QueryCtr != 1 || m2.QueryCtr != 1 {
		t.Fatalf("Expected m1.Query to be called once but was %d. Expected m2.Query to be called once but was %d\n", m1.QueryCtr, m2.QueryCtr)
	}
}

func Test_Enterprise_NewClientWithURL(t *testing.T) {
	t.Parallel()

	urls := []struct {
		url       string
		tls       bool
		shouldErr bool
	}{
		{"http://localhost:8086", false, false},
		{"https://localhost:8086", false, false},

		{"http://localhost:8086", true, false},
		{"https://localhost:8086", true, false},

		{"localhost:8086", false, false},
		{"localhost:8086", true, false},

		{":http", false, true},
	}

	for _, testURL := range urls {
		_, err := enterprise.NewClientWithURL(testURL.url, testURL.tls, log.New(log.DebugLevel))
		if err != nil && !testURL.shouldErr {
			t.Errorf("Unexpected error creating Client with URL %s and TLS preference %t. err: %s", testURL.url, testURL.tls, err.Error())
		} else if err == nil && testURL.shouldErr {
			t.Errorf("Expected error creating Client with URL %s and TLS preference %t", testURL.url, testURL.tls)
		}
	}
}

func Test_Enterprise_ComplainsIfNotOpened(t *testing.T) {
	m1 := NewMockTimeSeries("http://host-1.example.com:8086")
	cl := enterprise.NewClientWithTimeSeries(log.New(log.DebugLevel), chronograf.TimeSeries(m1))
	_, err := cl.Query(context.Background(), chronograf.Query{Command: "show shards", DB: "_internal", RP: "autogen"})
	if err != chronograf.ErrUninitialized {
		t.Error("Expected ErrUnitialized, but was this err:", err)
	}
}