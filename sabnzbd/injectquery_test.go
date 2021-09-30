package sabnzbd

import (
	"net/url"
	"testing"
)

func TestInjectQuery(t *testing.T) {
	type Foo struct {
		StrParam  string `query_name:"str"`
		IntParam  int32  `query_name:"test,omitempty"`
		Bar       string
		Something string `query_name:",omitempty"`
	}

	u, _ := url.Parse("http://example.com")

	params := Foo{
		StrParam: "luke",
		IntParam: 42,
		Bar:      "bar",
	}
	query := u.Query()

	if err := InjectQuery(query, params); err != nil {
		t.Fatalf("Unexpected error in InjectQuery: %s", err)
	}
	u.RawQuery = query.Encode()

	generatedUrl := u.String()
	const expectedUrl = "http://example.com?Bar=bar&str=luke&test=42"

	if generatedUrl != expectedUrl {
		t.Fatalf("Generated url (%s) is not the expected url (%s)", generatedUrl, expectedUrl)
	}
}

func TestInjectInUrl(t *testing.T) {
	type Foo struct {
		StrParam string `query_name:"str"`
		IntParam int32  `query_name:"test"`
		Bar      string
	}

	u, _ := url.Parse("http://example.com")
	params := Foo{
		StrParam: "luke",
		IntParam: 42,
		Bar:      "bar",
	}

	if err := InjectInUrl(u, params); err != nil {
		t.Fatalf("Unexpected error in InjectInUrl: %s", err)
	}

	generatedUrl := u.String()
	const expectedUrl = "http://example.com?Bar=bar&str=luke&test=42"

	if generatedUrl != expectedUrl {
		t.Fatalf("Generated url (%s) is not the expected url (%s)", generatedUrl, expectedUrl)
	}
}

func TestInjectInUrlWithArray(t *testing.T) {
	type Foo struct {
		StrArray []string `query_name:"str"`
		IntArray []int32  `query_name:"test"`
	}

	u, _ := url.Parse("http://example.com")
	params := Foo{
		StrArray: []string{"luke", "leia"},
		IntArray: []int32{42, 88},
	}

	if err := InjectInUrl(u, params); err != nil {
		t.Fatalf("Unexpected error in InjectInUrl: %s", err)
	}

	generatedUrl := u.String()
	const expectedUrl = "http://example.com?str=luke%2Cleia&test=42%2C88"

	if generatedUrl != expectedUrl {
		t.Fatalf("Generated url (%s) is not the expected url (%s)", generatedUrl, expectedUrl)
	}
}
