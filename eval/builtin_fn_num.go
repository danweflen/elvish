package eval

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/elves/elvish/eval/types"
)

// Numerical operations.

func init() {
	addToBuiltinFns([]*BuiltinFn{
		// Comparison
		{"<",
			wrapNumCompare(func(a, b float64) bool { return a < b })},
		{"<=",
			wrapNumCompare(func(a, b float64) bool { return a <= b })},
		{"==",
			wrapNumCompare(func(a, b float64) bool { return a == b })},
		{"!=",
			wrapNumCompare(func(a, b float64) bool { return a != b })},
		{">",
			wrapNumCompare(func(a, b float64) bool { return a > b })},
		{">=",
			wrapNumCompare(func(a, b float64) bool { return a >= b })},

		// Arithmetics
		{"+", plus},
		{"-", minus},
		{"*", times},
		{"/", slash},
		{"^", pow},
		{"%", mod},

		// Random
		{"rand", randFn},
		{"randint", randint},
	})
}

func wrapNumCompare(cmp func(a, b float64) bool) BuiltinFnImpl {
	return func(ec *Frame, args []interface{}, opts map[string]interface{}) {
		TakeNoOpt(opts)
		floats := make([]float64, len(args))
		for i, a := range args {
			f, err := toFloat(a)
			maybeThrow(err)
			floats[i] = f
		}
		result := true
		for i := 0; i < len(floats)-1; i++ {
			if !cmp(floats[i], floats[i+1]) {
				result = false
				break
			}
		}
		ec.OutputChan() <- types.Bool(result)
	}
}

func plus(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var nums []float64
	ScanArgsVariadic(args, &nums)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	sum := 0.0
	for _, f := range nums {
		sum += f
	}
	out <- floatToString(sum)
}

func minus(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var (
		sum  float64
		nums []float64
	)
	ScanArgsVariadic(args, &sum, &nums)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	if len(nums) == 0 {
		// Unary -
		sum = -sum
	} else {
		for _, f := range nums {
			sum -= f
		}
	}
	out <- floatToString(sum)
}

func times(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var nums []float64
	ScanArgsVariadic(args, &nums)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	prod := 1.0
	for _, f := range nums {
		prod *= f
	}
	out <- floatToString(prod)
}

func slash(ec *Frame, args []interface{}, opts map[string]interface{}) {
	TakeNoOpt(opts)
	if len(args) == 0 {
		// cd /
		cdInner("/", ec)
		return
	}
	// Division
	divide(ec, args, opts)
}

func divide(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var (
		prod float64
		nums []float64
	)
	ScanArgsVariadic(args, &prod, &nums)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	for _, f := range nums {
		prod /= f
	}
	out <- floatToString(prod)
}

func pow(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var b, p float64
	ScanArgs(args, &b, &p)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	out <- floatToString(math.Pow(b, p))
}

func mod(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var a, b int
	ScanArgs(args, &a, &b)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	out <- strconv.Itoa(a % b)
}

func randFn(ec *Frame, args []interface{}, opts map[string]interface{}) {
	TakeNoArg(args)
	TakeNoOpt(opts)

	out := ec.ports[1].Chan
	out <- floatToString(rand.Float64())
}

func randint(ec *Frame, args []interface{}, opts map[string]interface{}) {
	var low, high int
	ScanArgs(args, &low, &high)
	TakeNoOpt(opts)

	if low >= high {
		throw(ErrArgs)
	}
	out := ec.ports[1].Chan
	i := low + rand.Intn(high-low)
	out <- strconv.Itoa(i)
}
