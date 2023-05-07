package openapi

import (
	"os"
	"testing"

	"github.com/apicat/apicat/commom/spec"
)

// func TestDecode(t *testing.T) {
// 	fs := map[string]string{
// 		"swagger": "../../testdata/swagger.json",
// 		"openapi3.0": "../../testdata/openapi3.0.yaml",
// 		"openapi3.1": "../../testdata/openapi3.1.yaml",
// 	}
// 	for k, v := range fs {
// 		raw, err := os.ReadFile(v)
// 		if err != nil {
// 			t.Fatal(k, err)
// 		}
// 		if x, err := Decode(raw); err != nil {
// 			t.Fatal(k, err)
// 		} else {
// 			d, _ := x.ToJSON(spec.JSONOption{Indent: " "})
// 			fmt.Println(string(d))
// 		}
// 	}
// }

func TestEncode(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/spec.json")
	if err != nil {
		t.FailNow()
	}
	p, err := spec.ParseJSON(raw)
	if err != nil {
		t.FailNow()
	}

	for _, v := range []string{"3.0.0", "3.0.1", "2.0"} {
		_, err := Encode(p, v)
		if err != nil {
			t.Fatal(v, err)
		}
		// fmt.Println(string(p))
	}

}