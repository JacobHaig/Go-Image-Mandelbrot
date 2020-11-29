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

// Normalize is cool
func Normalize(value, minIn, maxIn, minOut, maxOut float64) float64 {
	return (minOut + (((value - minIn) * (maxOut - minOut)) / (maxIn - minIn)))
}
