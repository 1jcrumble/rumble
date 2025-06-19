package main

import "testing"

func TestMaskPAN(t *testing.T) {
	cases := []struct {
		in  int64
		out string
	}{
		{4080230386144446, "**** **** **** 4446"},
		{1234, "**** **** **** 1234"},
		{12, "****"},
	}
	for _, c := range cases {
		got := maskPAN(c.in)
		if got != c.out {
			t.Errorf("maskPAN(%d) == %q, want %q", c.in, got, c.out)
		}
	}
}

func TestSortTransactionsDesc_Timestamps(t *testing.T) {
	orig := DefaultMockTransactions()
	sorted := SortTransactionsDesc(orig)

	// expected timestamps in descending order
	want := []string{
		"2025-06-15T12:10:30+00:00",
		"2025-06-12T18:05:30+00:00",
		"2025-06-11T20:30:00+00:00",
		"2025-06-10T13:20:00+00:00",
		"2025-06-08T08:00:00+00:00",
		"2025-06-07T11:15:05+00:00",
		"2025-06-05T16:48:10+00:00",
		"2025-06-03T10:45:24+00:00",
		"2025-06-02T14:17:42+00:00",
		"2025-06-01T09:21:15+00:00",
	}

	for i := range want {
		got := sorted[i].PostedTimeStamp
		if got != want[i] {
			t.Errorf("index %d: got timestamp %q, want %q", i, got, want[i])
		}
	}
}
