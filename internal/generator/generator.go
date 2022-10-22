package generator

import (
	"math/rand"
	"strings"

	"github.com/lucasepe/grasp/internal/csprng"
)

const (
	// Letters is the list of letters.
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// Symbols is the list of symbols.
	Symbols = "!@#$%^&*()_+-={}|[]\\:<>?,./"
)

// Generator is the stateful generator which can be used to customize the list
// of letters, digits, and/or symbols.
type Generator struct {
	letters    string
	digits     string
	symbols    string
	randomizer *rand.Rand
}

// Option is used as input to the NewGenerator function.
type Option func(*Generator)

// WithLetters is a function option used
// to specify a custom set of letters
func WithLetters(letters string) Option {
	return func(gen *Generator) {
		gen.letters = letters
	}
}

// WithDigits is a function option used
// to specify a custom set of digits
func WithDigits(digits string) Option {
	return func(gen *Generator) {
		gen.digits = digits
	}
}

// WithSymbols is a function option used
// to specify a custom set of symbols
func WithSymbols(symbols string) Option {
	return func(gen *Generator) {
		gen.symbols = symbols
	}
}

// NewGenerator creates a new Generator from the specified configuration. If no
// input is given, all the default values are used. This function is safe for
// concurrent use.
func NewGenerator(secrets []string, opts ...Option) (*Generator, error) {
	src, err := csprng.New(secrets)
	if err != nil {
		return nil, err
	}

	res := &Generator{
		letters:    Letters,
		digits:     Digits,
		symbols:    Symbols,
		randomizer: rand.New(src),
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *Generator as the argument
		opt(res)
	}

	return res, nil
}

// Generate generates a password with the given requirements. length is the
// total number of characters in the password. numDigits is the number of digits
// to include in the result. numSymbols is the number of symbols to include in
// the result. noUpper excludes uppercase letters from the results. allowRepeat
// allows characters to repeat.
//
// The algorithm is fast, but it's not designed to be performant; it favors
// entropy over speed. This function is safe for concurrent use.
func (g *Generator) Generate(length int, noDigits, noSymbols, allowRepeat bool) (string, error) {

	chars := g.letters
	if !noDigits {
		chars = chars + g.digits
	}

	if !noSymbols {
		chars = chars + g.symbols
	}

	var result string

	// Characters
	for i := 0; i < length; i++ {
		ch := randomLetter(g.randomizer, chars)

		if !allowRepeat && strings.Contains(result, ch) {
			i--
			continue
		}

		result = randomInsert(g.randomizer, result, ch)
	}

	return result, nil
}
