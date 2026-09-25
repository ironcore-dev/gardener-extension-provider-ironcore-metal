// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ASNumber represents a BGP Autonomous System Number.
// It supports plain integer notation (e.g. 4270202922) and dotted "asdot" notation (e.g. "65161.5018").
// Dotted notation must be quoted in YAML to avoid float coercion (e.g. asNumber: "65161.5018").
// Internally it is stored as a uint32.
type ASNumber uint32

// ASNumberFromString parses an AS number from a string.
// Accepts plain integers ("4270202922") and dotted notation ("65161.5018").
func ASNumberFromString(s string) (ASNumber, error) {
	if num, err := strconv.ParseUint(s, 10, 32); err == nil {
		return ASNumber(num), nil
	}

	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid AS number format (%s)", s)
	}

	high, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid AS number format (%s)", s)
	}

	low, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid AS number format (%s)", s)
	}

	return ASNumber((high << 16) + low), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Accepts integer (4270202922) and quoted string ("65161.5018" or "4270202922").
// Unquoted dotted notation (65161.5018) is not supported because YAML float
// coercion drops trailing zeros (e.g. 1.10 becomes 1.1), causing silent
// data corruption.
func (a *ASNumber) UnmarshalJSON(b []byte) error {
	// Try integer first.
	if err := json.Unmarshal(b, (*uint32)(a)); err == nil {
		return nil
	}

	// Try quoted string.
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("AS number must be an integer or a quoted string: %s", string(b))
	}

	v, err := ASNumberFromString(s)
	if err != nil {
		return err
	}
	*a = v
	return nil
}

// MarshalJSON implements json.Marshaler. Serializes as an integer.
func (a ASNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint32(a))
}

// String returns the decimal string representation of the AS number.
func (a ASNumber) String() string {
	return strconv.FormatUint(uint64(a), 10)
}

// ToUint32 returns the AS number as a uint32.
func (a ASNumber) ToUint32() uint32 {
	return uint32(a)
}

// OpenAPISchemaType is used by the kube-openapi generator.
func (ASNumber) OpenAPISchemaType() []string { return []string{"string"} }

// OpenAPISchemaFormat is used by the kube-openapi generator.
func (ASNumber) OpenAPISchemaFormat() string { return "" }

// OpenAPIV3OneOfTypes is used by the kube-openapi generator to accept integer or string.
func (ASNumber) OpenAPIV3OneOfTypes() []string { return []string{"integer", "string"} }
