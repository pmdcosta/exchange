package recommendation

import (
	"fmt"

	"golang.org/x/text/currency"
)

// Naive makes a recommendation as to whether this is a good time to exchange money or not
// it is very simple and basically just compares the current rate against the average of the week
func Naive(current float64, currency currency.Unit, data map[string]map[currency.Unit]float64) (bool, string) {
	// calculate the average exchange value
	var average float64
	for _, c := range data {
		average += c[currency]
	}
	average = average / float64(len(data))

	// compare the value to the average
	if current > average {
		return true, fmt.Sprintf("the current exchange value is above the weekly average of%9f", average)
	}
	return false, fmt.Sprintf("the current exchange value is below the weekly average of%9f", average)
}
