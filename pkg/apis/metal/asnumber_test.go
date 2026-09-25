// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	"encoding/json"
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
		{name: "dotted max", input: "65535.65535", want: ASNumber(^uint32(0))},
		{name: "invalid format", input: "abc", wantErr: true},
		{name: "too many dots", input: "1.2.3", wantErr: true},
		{name: "high part too large", input: "65536.0", wantErr: true},
		{name: "low part too large", input: "0.65536", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
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
		{name: "unquoted dotted rejected", json: `{"asn":65161.5018}`, wantErr: true},
		{name: "invalid string", json: `{"asn":"abc"}`, wantErr: true},
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
