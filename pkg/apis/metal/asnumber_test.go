// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestASNumberFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ASNumber
		wantErr bool
	}{
		{name: "plain integer", input: "64512", want: 64512},
		{name: "large plain integer", input: "4270202922", want: 4270202922},
		{name: "dotted notation", input: "65161.5018", want: ASNumber((65161 << 16) + 5018)},
		{name: "dotted zero high", input: "0.64512", want: 64512},
		{name: "dotted zero both", input: "00.00", want: 0},
		{name: "dotted zero low", input: "00.01", want: 1},
		{name: "dotted max", input: "65535.65535", want: ASNumber(^uint32(0))},
		{name: "invalid format", input: "abc", wantErr: true},
		{name: "too many dots", input: "1.2.3", wantErr: true},
		{name: "high part too large", input: "65536.0", wantErr: true},
		{name: "low part too large", input: "0.65536", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "negative low part", input: "0.-1", wantErr: true},
		{name: "overflow", input: "4294967296", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ASNumberFromString(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got %d", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestASNumberUnmarshalJSON(t *testing.T) {
	type wrapper struct {
		ASN ASNumber `json:"asn"`
	}

	tests := []struct {
		name    string
		json    string
		want    ASNumber
		wantErr bool
	}{
		{name: "integer", json: `{"asn":64512}`, want: 64512},
		{name: "large integer", json: `{"asn":4270202922}`, want: 4270202922},
		{name: "quoted plain", json: `{"asn":"64512"}`, want: 64512},
		{name: "quoted dotted", json: `{"asn":"65161.5018"}`, want: ASNumber((65161 << 16) + 5018)},
		{name: "quoted dotted 1.10", json: `{"asn":"1.10"}`, want: ASNumber(65546)},
		{name: "quoted dotted zero", json: `{"asn":"00.00"}`, want: 0},
		{name: "unquoted dotted rejected", json: `{"asn":65161.5018}`, wantErr: true},
		{name: "overflow int", json: `{"asn":4294967296}`, wantErr: true},
		{name: "negative int", json: `{"asn":-1}`, wantErr: true},
		{name: "invalid string", json: `{"asn":"abc"}`, wantErr: true},
		{name: "negative dotted", json: `{"asn":"0.-1"}`, wantErr: true},
		{name: "boolean", json: `{"asn":true}`, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var w wrapper
			err := json.Unmarshal([]byte(tc.json), &w)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got %d", w.ASN)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if w.ASN != tc.want {
				t.Errorf("got %d, want %d", w.ASN, tc.want)
			}
		})
	}
}

func TestASNumberMarshalJSON(t *testing.T) {
	type wrapper struct {
		ASN ASNumber `json:"asn"`
	}

	w := wrapper{ASN: 4270202922}
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `{"asn":4270202922}`
	if string(b) != expected {
		t.Errorf("got %s, want %s", string(b), expected)
	}
}

func TestASNumberString(t *testing.T) {
	asn := ASNumber(4270202922)
	if s := asn.String(); s != "4270202922" {
		t.Errorf("got %s, want 4270202922", s)
	}
}

// TestASNumberHelmRoundTrip verifies that ASNumber values survive the JSON
// round-trip that occurs when Helm values are serialized. Go's
// json.Unmarshal into interface{} turns all numbers into float64.
// The chart template must use printf "%.0f" (not "%d") to render correctly.
func TestASNumberHelmRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string // JSON input as it arrives from the API
		want  string // expected rendered string in the helm template
	}{
		{name: "plain integer", input: `{"asn":65148}`, want: "65148"},
		{name: "large integer", input: `{"asn":4270396314}`, want: "4270396314"},
		{name: "quoted dotted notation", input: `{"asn":"65161.5018"}`, want: "4270396314"},
		{name: "quoted plain string", input: `{"asn":"65148"}`, want: "65148"},
		{name: "max ASN", input: `{"asn":4294967295}`, want: "4294967295"},
	}

	type wrapper struct {
		ASN ASNumber `json:"asn"`
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Unmarshal from API input (simulates shoot spec decoding)
			var w wrapper
			if err := json.Unmarshal([]byte(tc.input), &w); err != nil {
				t.Fatalf("unmarshal input: %v", err)
			}

			// Step 2: Values provider converts to uint32 for helm
			values := map[string]interface{}{"asNumber": w.ASN.ToUint32()}

			// Step 3: Simulate Helm JSON round-trip (marshal then unmarshal into interface{})
			b, err := json.Marshal(values)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var roundTripped map[string]interface{}
			if err := json.Unmarshal(b, &roundTripped); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			// After round-trip, value is float64
			fVal, ok := roundTripped["asNumber"].(float64)
			if !ok {
				t.Fatalf("expected float64 after JSON round-trip, got %T", roundTripped["asNumber"])
			}

			// Step 4: This is what the helm template does: printf "%.0f"
			got := fmt.Sprintf("%.0f", fVal)
			if got != tc.want {
				t.Errorf("printf %%.0f = %s, want %s", got, tc.want)
			}

			// Verify that %d would produce broken output (the bug we fixed).
			// Wrap in []any to prevent static analysis from seeing the float64 type.
			args := []any{fVal}
			badOutput := fmt.Sprintf("%d", args...)
			if badOutput == tc.want {
				t.Errorf("printf %%d unexpectedly produced correct output; float64 handling may have changed")
			}
		})
	}
}
