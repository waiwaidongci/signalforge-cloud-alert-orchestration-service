package id

type Allocator struct {
	generator IDGenerator
	issued    map[string]struct{}
}

func NewAllocator(generator IDGenerator) *Allocator {
	return &Allocator{generator: NewIDGenerator(generator)}
}
func (a *Allocator) AllocateID(prefix string) string {
	value := a.generator.Next(prefix)
	a.issued[value] = struct{}{}
	return value
}
func (a *Allocator) Issued(value string) bool { _, ok := a.issued[value]; return ok }
