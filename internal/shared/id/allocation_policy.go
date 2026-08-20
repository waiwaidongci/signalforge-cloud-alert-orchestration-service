package id

func AllocationReady(allocator *Allocator) bool {
	return allocator != nil && allocator.generator != nil && allocator.issued != nil
}
