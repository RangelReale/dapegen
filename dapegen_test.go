package dapegen

import (
	"github.com/RangelReale/epochdate"
	"testing"
	"time"
)

func TestGeneratorDay(t *testing.T) {
	d1, err := epochdate.Parse(epochdate.RFC3339, "2013-08-02")
	if err != nil {
		t.Fatal(err)
	}
	d2, err := epochdate.Parse(epochdate.RFC3339, "2013-08-10")
	if err != nil {
		t.Fatal(err)
	}
	d3, err := epochdate.Parse(epochdate.RFC3339, "2013-08-05")
	if err != nil {
		t.Fatal(err)
	}

	// creates a new generator by day
	p, err := NewGenerator(d1, d2, DAY)
	if err != nil {
		t.Fatal(err)
	}

	ct := 0
	for p.Next() {
		if ct == 3 && p.CurrentDate != d3 {
			t.Fatalf("Invalid date generated (%s - should be %s)", p.CurrentDate.String(), d3.String())
		}
		ct++
	}

	if ct != 9 {
		t.Fatalf("Should have generated 9 items, but generated %d", ct)
	}
}

func TestGeneratorWeek(t *testing.T) {
	d1, err := epochdate.Parse(epochdate.RFC3339, "2013-08-01") // thursday
	if err != nil {
		t.Fatal(err)
	}
	d2, err := epochdate.Parse(epochdate.RFC3339, "2013-08-25") // sunday
	if err != nil {
		t.Fatal(err)
	}

	d3, err := epochdate.Parse(epochdate.RFC3339, "2013-08-14")
	if err != nil {
		t.Fatal(err)
	}

	// creates a new generator by week
	p, err := NewGenerator(d1, d2, WEEK)
	p.FirstDayOfWeek = time.Wednesday
	if err != nil {
		t.Fatal(err)
	}

	ct := 0
	for p.Next() {
		if ct == 2 && p.CurrentDate != d3 {
			t.Fatalf("Invalid date generated (%s - should be %s)", p.CurrentDate.String(), d3.String())
		}
		ct++
	}

	if ct != 4 {
		t.Fatalf("Should have generated 4 items, but generated %d", ct)
	}
}

func TestGeneratorMonth(t *testing.T) {
	d1, err := epochdate.Parse(epochdate.RFC3339, "2013-05-12")
	if err != nil {
		t.Fatal(err)
	}
	d2, err := epochdate.Parse(epochdate.RFC3339, "2013-08-25")
	if err != nil {
		t.Fatal(err)
	}

	d3, err := epochdate.Parse(epochdate.RFC3339, "2013-07-01")
	if err != nil {
		t.Fatal(err)
	}

	// creates a new generator by month
	p, err := NewGenerator(d1, d2, MONTH)
	if err != nil {
		t.Fatal(err)
	}

	ct := 0
	for p.Next() {
		if ct == 2 && p.CurrentDate != d3 {
			t.Fatalf("Invalid date generated (%s - should be %s)", p.CurrentDate.String(), d3.String())
		}
		ct++
	}

	if ct != 4 {
		t.Fatalf("Should have generated 4 items, but generated %d", ct)
	}
}

func TestGeneratorUntil(t *testing.T) {
	d1, err := epochdate.Parse(epochdate.RFC3339, "2013-04-11") // thursday
	if err != nil {
		t.Fatal(err)
	}
	d2, err := epochdate.Parse(epochdate.RFC3339, "2013-07-17") // sunday
	if err != nil {
		t.Fatal(err)
	}
	du, err := epochdate.Parse(epochdate.RFC3339, "2013-06-04") // tuesday
	if err != nil {
		t.Fatal(err)
	}

	// creates a new generator by day
	p, err := NewGenerator(d1, d2, DAY)
	if err != nil {
		t.Fatal(err)
	}

	ct := 0

	// loop 4 times to test if nextuntil is really stopping
	for i := 0; i < 4; i++ {
		for p.NextUntil(du) {
			ct++
		}
	}

	if ct != 54 {
		t.Fatalf("Should have generated 54 items, but generated %d", ct)
	}

	if p.CurrentDate != du {
		t.Fatalf("Invalid current date (%s - should be %s)", p.CurrentDate.String(), du.String())
	}

	ct++ // count stopped date

	for p.Next() {
		ct++
	}

	if ct != 98 {
		t.Fatalf("Should have generated 98 items, but generated %d", ct)
	}
}
