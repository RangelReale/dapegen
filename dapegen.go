package dapegen

import (
	"errors"
	"github.com/RangelReale/epochdate"
	"time"
)

// Grouping enum
type GROUP int

const (
	DAY GROUP = iota
	WEEK
	MONTH
)

// Generator object
type Generator struct {
	StartDate      epochdate.Date
	EndDate        epochdate.Date
	Group          GROUP
	FirstDate      epochdate.Date
	LastDate       epochdate.Date
	CurrentDate    epochdate.Date
	FirstDayOfWeek time.Weekday

	currentEndDate epochdate.Date
	isfirst        bool
	isforward      bool
}

// Creates a forward generator
func NewGenerator(startdate epochdate.Date, enddate epochdate.Date, group GROUP) (*Generator, error) {
	d := &Generator{
		StartDate:      startdate,
		EndDate:        enddate,
		Group:          group,
		FirstDayOfWeek: time.Monday,
		isfirst:        true,
		isforward:      true,
	}
	if err := d.initialize(); err != nil {
		return nil, err
	}
	return d, nil
}

// Creates a backwards generator
func NewGeneratorBackwards(startdate epochdate.Date, enddate epochdate.Date, group GROUP) (*Generator, error) {
	d := &Generator{
		StartDate:      startdate,
		EndDate:        enddate,
		Group:          group,
		FirstDayOfWeek: time.Monday,
		isfirst:        true,
		isforward:      false,
	}
	if err := d.initialize(); err != nil {
		return nil, err
	}
	return d, nil
}

// checks the passed parameters
func (d *Generator) initialize() error {
	if d.isforward {
		if d.StartDate.After(d.EndDate) {
			return errors.New("Start date must be before end date")
		}
	} else {
		if d.StartDate.Before(d.EndDate) {
			return errors.New("Start date must be after end date")
		}
	}
	return nil
}

// returns the first day from the date, depending on the group
func (d Generator) FirstFromDate(date epochdate.Date) epochdate.Date {
	switch d.Group {
	case WEEK:
		// first date is the previous monday of the week
		wd := date.UTC().Weekday()
		if wd > d.FirstDayOfWeek {
			return date - epochdate.Date(wd) + epochdate.Date(d.FirstDayOfWeek)
		} else if wd < d.FirstDayOfWeek {
			return date - 6 + epochdate.Date(wd)
		}
	case MONTH:
		// first day of the month
		year, month, _ := date.Date()
		ret, _ := epochdate.NewFromDate(year, month, 1)
		return ret
	}
	return date
}

// Generates the next period, returns false if finished
func (d *Generator) Next() bool {
	ret, _ := d.nextUntilInternal(nil, true)
	return ret
}

// Generates the next period until the passed date, returns if have next or finished
func (d *Generator) NextUntil(date epochdate.Date) (havenext bool, isfinished bool) {
	f := d.FirstFromDate(date)
	havenext, isfinished = d.nextUntilInternal(&f, false)
	return
}

// Returns if is finished
func (d *Generator) IsFinished() bool {
	return d.isFinishedInternal(d.currentEndDate)
}

// Initialize internal values on first pass
func (d *Generator) initializeGeneration() {
	d.CurrentDate = d.FirstFromDate(d.StartDate)
	d.currentEndDate = d.FirstFromDate(d.EndDate)

	d.FirstDate = d.CurrentDate
	d.LastDate = d.currentEndDate
}

// Calculate the next period until the passed date, or the end date if nil.
// includelast determines if the last item is included.
func (d *Generator) nextUntilInternal(date *epochdate.Date, includelast bool) (havenext bool, isfinished bool) {
	if d.isfirst {
		// initialize on first run
		d.initializeGeneration()
	}

	if date == nil {
		// use default end date
		date = &d.currentEndDate
	}

	if d.isFinishedInternal(*date) {
		d.isfirst = false
		return false, true
	}

	if !d.isfirst {
		switch d.Group {
		case WEEK:
			// advance 7 days
			if d.isforward {
				d.CurrentDate += 7
			} else {
				d.CurrentDate -= 7
			}
		case MONTH:
			// advance day 1 of month
			year, month, _ := d.CurrentDate.Date()
			if d.isforward {
				d.CurrentDate, _ = epochdate.NewFromDate(year, month+1, 1)
			} else {
				d.CurrentDate, _ = epochdate.NewFromDate(year, month-1, 1)
			}
		default:
			// advance 1 day
			if d.isforward {
				d.CurrentDate++
			} else {
				d.CurrentDate--
			}
		}
	} else {
		d.isfirst = false
	}
	if !includelast && d.isFinishedInternal(*date) {
		return false, false
	}

	return true, false

}

// checks if iteration finished on passed date
func (d *Generator) isFinishedInternal(date epochdate.Date) bool {
	if d.CurrentDate == date {
		return true
	}

	if d.isforward {
		return !d.CurrentDate.Before(date)
	} else {
		return !d.CurrentDate.After(date)
	}
}
