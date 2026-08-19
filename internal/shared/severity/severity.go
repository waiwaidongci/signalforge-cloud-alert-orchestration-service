package severity

type Severity string

const (
	Critical Severity = "critical"
	High     Severity = "high"
	Medium   Severity = "medium"
	Low      Severity = "low"
	Info     Severity = "info"
)

func (s Severity) String() string {
	return string(s)
}

func Valid(raw string) bool {
	switch Severity(raw) {
	case Critical, High, Medium, Low, Info:
		return true
	default:
		return false
	}
}

func Order(s Severity) int {
	switch s {
	case Critical:
		return 5
	case High:
		return 4
	case Medium:
		return 3
	case Low:
		return 2
	case Info:
		return 1
	default:
		return 0
	}
}

func HigherOrEqual(a, b Severity) bool {
	return Order(a) >= Order(b)
}

func All() []string {
	return []string{Critical.String(), High.String(), Medium.String(), Low.String(), Info.String()}
}
