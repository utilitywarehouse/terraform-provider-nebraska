package nebraska

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

var (
	testUserAgent = "test-user-agent"
)

func testClientServer(handler func(w http.ResponseWriter, r *http.Request)) (*Client, *httptest.Server) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	c := New(s.URL, testUserAgent)
	return c, s
}

type testSchema struct {
	Name       string   `json:"name"`
	Parameters []string `json:"parameters"`
}

func TestClientRequest(t *testing.T) {
	expectedReqBody := &testSchema{
		Name:       "foo",
		Parameters: []string{"one", "two", "three"},
	}
	expectedRespBody := &testSchema{
		Name:       "bar",
		Parameters: []string{"apple", "orange", "banana"},
	}
	c, s := testClientServer(func(w http.ResponseWriter, r *http.Request) {
		assert.DeepEqual(t, r.Header.Get("User-Agent"), testUserAgent)

		var ts testSchema

		err := json.NewDecoder(r.Body).Decode(&ts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		assert.DeepEqual(t, &ts, expectedReqBody)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedRespBody)
	})
	defer s.Close()
	req, err := c.newRequest(http.MethodPost, "/", expectedReqBody)
	if err != nil {
		t.Fatal(err)
	}
	data := &testSchema{}
	if err := c.do(req, data); err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, data, expectedRespBody)
}
