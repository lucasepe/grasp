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
			[]string{"mail.google.com", "pinco.pallo@gmail.com"},
			6, "|Lv2\\N",
		},
		{
			[]string{"twitter.com", "monster.hunter@email.com"},
			9, "@z21o<K&>",
		},
		{
			[]string{"yahoo.com", "dev.to.dev@yahoo.com"},
			12, "ARbCD5}&7Mxu",
		},
		{
			[]string{"facebook.com", "mev@github.com"},
			20, "w(W4C+@SOH[FI5*aryjt",
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
			[]string{"mail.google.com", "pinco.pallo@gmail.com"},
			6, "dvTEfL",
		},
		{
			[]string{"twitter.com", "monster.hunter@email.com"},
			9, "q0AlZ9cGy",
		},
		{
			[]string{"yahoo.com", "dev.to.dev@yahoo.com"},
			12, "GldprPXSAKxt",
		},
		{
			[]string{"facebook.com", "mev@github.com"},
			20, "u1BKgqY60zFixbQXVWDj",
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
			[]string{"mail.google.com", "pinco.pallo@gmail.com"},
			6, "dvTEfL", true,
		},
		{
			[]string{"twitter.com", "monster.hunter@email.com"},
			9, "q0AlZ9cGy", true,
		},
		{
			[]string{"yahoo.com", "dev.to.dev@yahoo.com"},
			12, "GlpIPSASKt1C", false,
		},
		{
			[]string{"facebook.com", "mev@github.com"},
			20, "VKqiuB1bDzx0KglYs0WH", false,
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
