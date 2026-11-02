package repositories

import (
	"errors"
	"reflect"
	"testing"
)

type fakePageRows struct {
	pages   []PageStat
	index   int
	scanErr error
	rowsErr error
}

func (r *fakePageRows) Next() bool {
	if r.index >= len(r.pages) {
		return false
	}
	r.index++
	return true
}

func (r *fakePageRows) Scan(dest ...any) error {
	if r.scanErr != nil {
		return r.scanErr
	}
	page := r.pages[r.index-1]
	*dest[0].(*string) = page.Page
	*dest[1].(*uint64) = page.Count
	return nil
}

func (r *fakePageRows) Err() error { return r.rowsErr }

func TestCollectPageStats(t *testing.T) {
	want := []PageStat{{Page: "/home", Count: 5}, {Page: "/pricing", Count: 2}}
	got, err := collectPageStats(&fakePageRows{pages: want})
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, %v; want %v, nil", got, err, want)
	}
}

func TestCollectPageStatsRejectsIncompleteResults(t *testing.T) {
	iterationErr := errors.New("connection lost during iteration")
	scanErr := errors.New("invalid row")
	tests := []struct {
		name string
		rows fakePageRows
		want error
	}{
		{"iteration error after a row", fakePageRows{pages: []PageStat{{Page: "/home", Count: 5}}, rowsErr: iterationErr}, iterationErr},
		{"scan error", fakePageRows{pages: []PageStat{{Page: "/home", Count: 5}}, scanErr: scanErr}, scanErr},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := collectPageStats(&tc.rows)
			if got != nil || !errors.Is(err, tc.want) {
				t.Fatalf("got %v, %v; want nil, %v", got, err, tc.want)
			}
		})
	}
}
