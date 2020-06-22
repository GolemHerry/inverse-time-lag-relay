package fft

import "math"

type FFT struct {
	value     complex128
	amplitude float64
	phase     float64
}

func Calculate(sample []float64) *FFT {
	sum1, sum2, length := 0.0, 0.0, len(sample)
	for i, v := range sample {
		sum1 += v * math.Sin(float64(2*(i+1))*math.Pi/float64(length))
		sum2 += v * math.Cos(float64(2*(i+1))*math.Pi/float64(length))
	}
	return &FFT{
		value:     complex(sum1*2/float64(length), sum2*2/float64(length)),
		amplitude: math.Hypot(sum1*2/float64(length), sum2*2/float64(length)),
		phase:     _getPhase(sum1*2/float64(length), sum2*2/float64(length)),
	}
}

func (f *FFT) Phase() float64 {
	return f.phase
}

func (f *FFT) Amplitude() float64 {
	return f.amplitude
}

func _getPhase(x, y float64) float64 {
	res := math.Atan(math.Abs(y / x))
	if x < 0 {
		res += math.Pi / 2
	}
	if y < 0 {
		res = -res
	}
	return res
}
