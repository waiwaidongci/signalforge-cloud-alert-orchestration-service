package id

type Allocator struct {
	generator IDGenerator
	issued    map[string]struct{}
}

func NewAllocator(generator IDGenerator) *Allocator {
	return &Allocator{
		generator: NewIDGenerator(generator),
		issued:    make(map[string]struct{}),
	}
}
func (a *Allocator) AllocateID(prefix string) string {
	if a.generator == nil {
		a.generator = NewIDGenerator(nil)
	}
	if a.issued == nil {
		a.issued = make(map[string]struct{})
	}
	value := a.generator.Next(prefix)
	a.issued[value] = struct{}{}
	return value
}
func (a *Allocator) Issued(value string) bool { _, ok := a.issued[value]; return ok }
