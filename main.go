package randomFunctionFirer

import (
	"errors"
	"math/rand"
	"sort"
)

type FirerFunction struct {
	F                func()
	ExactProbability float64
	PartProbability  float64
}

func CreateFirerFunctionWithExactProbability(f func(), probability float64) FirerFunction {
	return FirerFunction{F: f, ExactProbability: probability}
}

func CreateFirerFunctionWithPartProbability(f func(), probability float64) FirerFunction {
	return FirerFunction{F: f, PartProbability: probability}
}

type RandomFunctionFirer struct {
	funcs  []FirerFunction
	limits []float64
}

func CreateFunctionFirer() RandomFunctionFirer {
	return RandomFunctionFirer{funcs: make([]FirerFunction, 0), limits: make([]float64, 0)}
}

func (fr *RandomFunctionFirer) CalculateLimits() {
	var totalProbability float64 = 0

	newLimits := make([]float64, len(fr.funcs))
	totalParts := float64(0)
	indexOfPartProbFuncs := -1

	for index, function := range fr.funcs {
		if function.ExactProbability == 0 && function.PartProbability != 0 {
			totalParts += function.PartProbability

			if indexOfPartProbFuncs < 0 {
				indexOfPartProbFuncs = index
			}
		} else {
			totalProbability += function.ExactProbability
			newLimits[index] = totalProbability
		}
	}

	if indexOfPartProbFuncs > -1 {
		var leftProbability float64 = float64(1) - totalProbability
		partOfLeftProbability := leftProbability / totalParts

		for i := indexOfPartProbFuncs; i < len(fr.funcs); i += 1 {
			totalProbability += partOfLeftProbability * fr.funcs[i].PartProbability
			newLimits[i] = totalProbability
		}
	}

	fr.limits = newLimits
}

func (fr *RandomFunctionFirer) AddFunction(f FirerFunction) error {
	fr.funcs = append(fr.funcs, f)
	var totalProbability float64 = 0
	for _, function := range fr.funcs {
		totalProbability += function.ExactProbability
	}

	if totalProbability > float64(1) {
		fr.funcs = fr.funcs[:len(fr.funcs)-1]
		return errors.New("sum of probabilities is greater that 1")
	}

	sort.Slice(fr.funcs, func(a, b int) bool {
		if fr.funcs[a].ExactProbability == 0 && fr.funcs[a].PartProbability == 0 {
			return false
		}

		if fr.funcs[b].ExactProbability == 0 && fr.funcs[b].PartProbability == 0 {
			return false
		}

		if fr.funcs[a].PartProbability != 0 {
			return true
		}

		if fr.funcs[b].PartProbability != 0 {
			return true
		}

		return false
	})

	fr.CalculateLimits()

	return nil
}

func (fr *RandomFunctionFirer) FireFunction() {
	val := rand.Float64()

	for index, limit := range fr.limits {
		if limit > val {
			fr.funcs[index].F()
			return
		}
	}
}
