package generator

import (
	"fmt"
	"testing"
)

func TestGeneratorDisallowsRepeat(t *testing.T) {
	cases := []struct {
		secrets []string
		size    int
		want    string
	}{
		{
			[]string{"mail.google.com", "pinco.pallo@gmail.com", "this", "is", "sparta!"},
			6, "%{sW:y",
		},
		{
			[]string{"twitter.com", "monster.hunter@email.com", "you", "never", "know$$"},
			9, "%|<T.Ba({",
		},
		{
			[]string{"yahoo.com", "dev.to.dev@yahoo.com", "scarlet", "is", "!hot!"},
			12, "V[Nb#xaI*46\\",
		},
		{
			[]string{"facebook.com", "mev@github.com", "look", "here"},
			20, "!He*Fj-zTwOWGVBQ3Z8&",
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Initialize the Mersenne Twister PRNG
			gen, err := NewGenerator(tt.secrets)
			if err != nil {
				t.Fatal(err)
			}

			got, err := gen.Generate(tt.size, false, false, false)
			if err != nil {
				t.Fatal(err)
			}

			if !unique(got) {
				t.Errorf("string [%v] has not unique chars", got)
			}

			if got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestGeneratorNoSymbolsDisallowsRepeat(t *testing.T) {
	cases := []struct {
		secrets []string
		size    int
		want    string
	}{
		{
			[]string{"mail.google.com", "pinco.pallo@gmail.com", "this", "is", "sparta!"},
			6, "Lm5U3z",
		},
		{
			[]string{"twitter.com", "monster.hunter@email.com", "you", "never", "know$$"},
			9, "CjyKd4VOv",
		},
		{
			[]string{"yahoo.com", "dev.to.dev@yahoo.com", "scarlet", "is", "!hot!"},
			12, "IBHwjoVFYUzS",
		},
		{
			[]string{"facebook.com", "mev@github.com", "look", "here"},
			20, "kwpIhel83U1VZLjTdK5F",
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Initialize the Mersenne Twister PRNG
			gen, err := NewGenerator(tt.secrets)
			if err != nil {
				t.Fatal(err)
			}

			got, err := gen.Generate(tt.size, false, true, false)
			if err != nil {
				t.Fatal(err)
			}

			if !unique(got) {
				t.Errorf("string [%v] has not unique chars", got)
			}

			if got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestGeneratorNoSymbolsAllowsRepeat(t *testing.T) {
	cases := []struct {
		secrets        []string
		size           int
		want           string
		differentChars bool
	}{
		{
			[]string{"mail.google.com", "pinco.pallo@gmail.com", "this", "is", "sparta!"},
			6, "Lm5U3z", true,
		},
		{
			[]string{"twitter.com", "monster.hunter@email.com", "you", "never", "know$$"},
			9, "Cmjyy4VOd", false,
		},
		{
			[]string{"yahoo.com", "dev.to.dev@yahoo.com", "scarlet", "is", "!hot!"},
			12, "pHfHUwbzVOJQ", false,
		},
		{
			[]string{"facebook.com", "mev@github.com", "look", "here"},
			20, "ckwlpkselXU1ZvLLjTKP", false,
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Initialize the Mersenne Twister PRNG
			gen, err := NewGenerator(tt.secrets)
			if err != nil {
				t.Fatal(err)
			}

			got, err := gen.Generate(tt.size, false, true, true)
			if err != nil {
				t.Fatal(err)
			}

			if u := unique(got); u != tt.differentChars {
				t.Errorf("[%v] has unique chars, wants different [%t]", got, tt.differentChars)
			}

			if got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func unique(arr string) bool {
	m := make(map[rune]bool)
	for _, i := range arr {
		_, ok := m[i]
		if ok {
			return false
		}

		m[i] = true
	}

	return true
}
