package curve

import "math"

var (
	StdCurve *Curve
)

type Curve struct {
	Iop float64 `config:"iop"`
	C   float64 `config:"c"`
	K   float64 `config:"k"`
	NC  []float64
}

func Init(iOP float64, c, k float64) {
	var NC []float64
	for i := 0; i <= 20; i++ {
		NC = append(NC, math.Pow(float64(i), c))
	}
	StdCurve = &Curve{
		Iop: iOP,
		C:   c,
		K:   k * 1000,
		NC:  NC,
	}
}

/*
			k					 I	  								  ∆n				C··(C-k+1)    ∆n
	t = --------		fi ≈ (-------)^c   ≈  (N + ∆n)^c   ≈  N^c(1+ -----)^c  ≈ N^c(∑ -------------(----)^k )
		 fi - 1					Iop									   N				 	k!         N
*/
//return ms
func (c *Curve) GetTime(amplitude float64) float64 {
	if amplitude < c.Iop {
		return -1
	}
	proportion := amplitude / c.Iop
	N := math.Trunc(proportion)
	deltaN := proportion - N
	fi := c.NC[int(N)] * (1 + c.C*deltaN/N + c.C*(c.C-1)/2*math.Pow(deltaN/N, 2))
	return c.K / (fi - 1)
}
