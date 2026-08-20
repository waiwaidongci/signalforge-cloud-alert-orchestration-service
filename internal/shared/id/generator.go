package id

type IDGenerator interface{ Next(string) string }
type prefixGenerator struct{}

func (prefixGenerator) Next(prefix string) string      { return Prefix(prefix) }
func NewIDGenerator(candidate IDGenerator) IDGenerator { return candidate }
