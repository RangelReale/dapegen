package main

import (
	"fmt"
	"github.com/RangelReale/dapegen"
	"github.com/RangelReale/epochdate"
)

func main() {
	d1, err := epochdate.Parse(epochdate.RFC3339, "2013-04-11") // thursday
	if err != nil {
		panic(err)
	}
	d2, err := epochdate.Parse(epochdate.RFC3339, "2013-07-17") // sunday
	if err != nil {
		panic(err)
	}
	du, err := epochdate.Parse(epochdate.RFC3339, "2013-06-04") // tuesday
	if err != nil {
		panic(err)
	}

	p, err := dapegen.NewGenerator(d1, d2, dapegen.DAY)
	//p.FirstDayOfWeek = time.Wednesday
	if err != nil {
		panic(err)
	}

	// print first days until du
	// loop 4 times to test if nextuntil is really stopping
	for i := 0; i < 4; i++ {
		for pn, _ := p.NextUntil(du); pn; pn, _ = p.NextUntil(du) {
			fmt.Printf("%s\n", p.CurrentDate.String())
		}
		fmt.Printf("@@ %s\n", p.CurrentDate.String())
	}

	// print remaining days
	for p.Next() {
		fmt.Printf("!! %s\n", p.CurrentDate.String())
	}
}
