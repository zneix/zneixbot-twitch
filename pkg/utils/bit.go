package utils

// HasBits checks if sum contains bit by performing a bitwise AND operation between values
func HasBits(sum uint64, bit uint64) bool {
	return (sum & bit) == bit
}
