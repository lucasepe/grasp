package csprng

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSource(t *testing.T) {
	cases := []struct {
		secrets []string
		pulls   int
		max     int64
		want    string
	}{
		{
			[]string{"reddit.com", "my", "great", "secret", "easy", "to", "remember"},
			12, 1000,
			"WzIzNSwzMjcsODI2LDY3NywxOTcsNTY0LDgyNiw1MTAsNDc2LDI0MSw4NTIsNzI1XQ==",
		},
		{
			[]string{"facebook.com", "i", "like", "grasp"},
			18, 1000,
			"WzM0Nyw1NjAsMywzNjIsMTUzLDIzNyw5MzQsNTU5LDM0NiwyODQsNjEzLDQ3MywyMzcsNTE4LDMxOCw0MTQsMjc2LDcxNV0=",
		},
		{
			[]string{"google.it", "pinco.pallo@gmail.com", "another", "one", "bites", "the", "dust!"},
			10, 500,
			"WzE0MywyODQsMjc3LDQ1NCwyNDMsNDE3LDMzNiwxMCwyMzQsMzIyXQ==",
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			src, err := New(tt.secrets)
			if err != nil {
				t.Fatal(err)
			}

			r := rand.New(src)

			b := make([]int64, tt.pulls)
			for i := 0; i < tt.pulls; i++ {
				b[i] = r.Int63n(tt.max)
			}

			data, err := json.Marshal(&b)
			if err != nil {
				t.Fatal(err)
			}

			got := base64.StdEncoding.EncodeToString(data)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}

}
