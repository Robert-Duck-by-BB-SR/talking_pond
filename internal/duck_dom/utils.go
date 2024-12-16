package duckdom

import "fmt"

// in case we will need it later in our lifes
func conv_perc_to_string(dyn_value int, perc float64) string{
	return fmt.Sprint(int(float64(dyn_value) * perc))
}
