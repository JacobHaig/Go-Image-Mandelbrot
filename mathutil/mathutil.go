package mathutil

// Clamp simply return the clamped value given a certain range.
// A clamped value is a value that has to be between two values
// If it is below or above the max value it will return the max value
func Clamp(value, min, max int16) int16 {
	if value > max {
		return max
	} else if value < min {
		return min
	}
	return value
}

// Normalize function takes a value between a range of numbers and normalize
// it between a new range of numbers. For instance, if the range is 10..20 with
// a value of 15, and the new range is 0..100, the new value will be 50.
func Normalize(value, minIn, maxIn, minOut, maxOut float64) float64 {
	return (minOut + (((value - minIn) * (maxOut - minOut)) / (maxIn - minIn)))
}
