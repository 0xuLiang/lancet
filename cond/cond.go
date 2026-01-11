package cond

func Val[T any](cond bool, ifVal, elseVal T) T {
	if cond {
		return ifVal
	} else {
		return elseVal
	}
}

func Fun[T any](cond bool, ifFun, elseFun func() T) T {
	if cond {
		return ifFun()
	} else {
		return elseFun()
	}
}

func FunVoid(cond bool, ifFun, elseFun func()) {
	if cond {
		ifFun()
	} else {
		elseFun()
	}
}

func Or[T comparable](vals ...T) T {
	var zero T
	for _, val := range vals {
		if val != zero {
			return val
		}
	}
	return zero
}
