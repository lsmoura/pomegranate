package sabnzbd

import (
	"net/url"
	"testing"
)

func TestInjectQuery(t *testing.T) {
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
