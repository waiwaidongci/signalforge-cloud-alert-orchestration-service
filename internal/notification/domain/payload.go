package domain

type Payload map[string]any

func NewPayload() Payload { return make(Payload) }

func (p Payload) Clone() Payload {
	if p == nil {
		return NewPayload()
	}
	clone := make(Payload, len(p))
	for key, value := range p {
		clone[key] = value
	}
	return clone
}
