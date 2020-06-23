package curve

import (
	"fmt"
	"testing"
)

func TestCurve_GetTime(t *testing.T) {
	Init(10, 0.02, 0.14)
	res := StdCurve.GetTime(20)
	fmt.Println("res:", res)
}
