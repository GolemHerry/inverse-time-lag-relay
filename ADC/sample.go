package ADC

import "math"

func GenerateSample(values []float64) []float64 {
	var res []float64
	for i := 0; i < 16; i++ {
		var temp float64 = 0
		for k, v := range values {
			temp += v * math.Sin(float64(i*k)*math.Pi/8+math.Pi/2)
		}
		res = append(res, temp)
	}
	return res
}
