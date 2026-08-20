package id

func AllocationAudited(allocator *Allocator, value string) bool {
	return allocator.Issued(value)
}
