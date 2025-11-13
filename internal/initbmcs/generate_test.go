// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

package initbmcs

import (
	"reflect"
	"testing"

	"bootstrap/internal/inventory"
)

func TestParseChassisSpec(t *testing.T) {
	got := ParseChassisSpec("x9000c1=02:23:28:01, x9000c3=02:23:28:03")
	want := map[string]string{
		"x9000c1": "02:23:28:01",
		"x9000c3": "02:23:28:03",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseChassisSpec mismatch: got=%v want=%v", got, want)
	}
}

func TestGenerateSingleChassisDeterministic(t *testing.T) {
	chassis := map[string]string{"x9000c1": "02:23:28:01"}
	bmcs := Generate(chassis, 4, 2, 1, "192.168.100")

	want := []inventory.Entry{
		{Xname: "x9000c1s0b0", MAC: "02:23:28:01:30:00", IP: "192.168.100.1"},
		{Xname: "x9000c1s0b1", MAC: "02:23:28:01:30:10", IP: "192.168.100.2"},
	}
	if !reflect.DeepEqual(bmcs, want) {
		t.Fatalf("Generate result mismatch:\n got: %#v\nwant: %#v", bmcs, want)
	}
}
