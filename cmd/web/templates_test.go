package main

import (
	"testing"
	"time"

	"github.com/arynkh/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	//create a slice of anonymous structs to hold the test case name, input to our humanDate() func and expect output (the want field)
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2025 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "GET",
			tm:   time.Date(2025, 3, 17, 10, 15, 0, 0, time.FixedZone("GET", 1*60*60)),
			want: "17 Mar 2025 at 09:15",
		},
	}

	for _, tt := range tests {
		//use t.Run to run a subtest for each test case. The first param is the name of the test and the second is an anonymous function containing the actual test
		t.Run(tt.name, func(t *testing.T) {
			//call the humanDate function with the time.Time value from the test case
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
