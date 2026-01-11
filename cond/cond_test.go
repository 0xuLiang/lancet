package cond

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVal(t *testing.T) {
	// Test with string
	assert.Equal(t, "trueVal", Val(true, "trueVal", "falseVal"))
	assert.Equal(t, "falseVal", Val(false, "trueVal", "falseVal"))

	// Test with int
	assert.Equal(t, 1, Val(true, 1, 2))
	assert.Equal(t, 2, Val(false, 1, 2))

	// Test with float
	assert.Equal(t, 1.1, Val(true, 1.1, 2.2))
	assert.Equal(t, 2.2, Val(false, 1.1, 2.2))

	// Test with bool
	assert.Equal(t, true, Val(true, true, false))
	assert.Equal(t, false, Val(false, true, false))
}

func TestFun(t *testing.T) {
	// Test with string
	trueFun := func() string { return "trueFun" }
	falseFun := func() string { return "falseFun" }
	assert.Equal(t, "trueFun", Fun(true, trueFun, falseFun))
	assert.Equal(t, "falseFun", Fun(false, trueFun, falseFun))

	// Test with int
	trueFunInt := func() int { return 1 }
	falseFunInt := func() int { return 2 }
	assert.Equal(t, 1, Fun(true, trueFunInt, falseFunInt))
	assert.Equal(t, 2, Fun(false, trueFunInt, falseFunInt))

	// Test with float
	trueFunFloat := func() float64 { return 1.1 }
	falseFunFloat := func() float64 { return 2.2 }
	assert.Equal(t, 1.1, Fun(true, trueFunFloat, falseFunFloat))
	assert.Equal(t, 2.2, Fun(false, trueFunFloat, falseFunFloat))

	// Test with bool
	trueFunBool := func() bool { return true }
	falseFunBool := func() bool { return false }
	assert.Equal(t, true, Fun(true, trueFunBool, falseFunBool))
	assert.Equal(t, false, Fun(false, trueFunBool, falseFunBool))
}

func TestFunVoid(t *testing.T) {
	var result string
	trueFun := func() { result = "trueFun" }
	falseFun := func() { result = "falseFun" }

	// Test when condition is true
	FunVoid(true, trueFun, falseFun)
	assert.Equal(t, "trueFun", result)

	// Test when condition is false
	FunVoid(false, trueFun, falseFun)
	assert.Equal(t, "falseFun", result)
}

func TestOr(t *testing.T) {
	// Test with int
	assert.Equal(t, 1, Or(0, 0, 0, 1, 0))
	assert.Equal(t, 0, Or(0, 0, 0, 0, 0))

	// Test with float
	assert.Equal(t, 1.1, Or(0.0, 0.0, 0.0, 1.1, 0.0))
	assert.Equal(t, 0.0, Or(0.0, 0.0, 0.0, 0.0, 0.0))

	// Test with string
	assert.Equal(t, "nonZero", Or("", "", "", "nonZero", ""))
	assert.Equal(t, "", Or("", "", "", "", ""))

	// Test with bool
	assert.Equal(t, true, Or(false, false, false, true, false))
	assert.Equal(t, false, Or(false, false, false, false, false))
}
