package currency

import "fmt"

// Format formats cent representation of currency into dollars with cents
// Ex. 1 -> 0.01$
// Ex 50 -> 0.50$
// Ex 1155 -> 11.55$
func Format(value int) string {
	cents := value % 100
	dollars := value / 100
	return fmt.Sprintf("%d.%02.f$", dollars, float64(cents))
}
