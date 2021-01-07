package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type size struct {
	key string
	val int
}

func (r *size) String() string {
	return fmt.Sprintf("%s=%d", r.key, r.val)
}

type sizeList []size

func (p sizeList) String() string {
	ret := make([]string, len(p))
	for i, x := range p {
		ret[i] = x.String()
	}

	return strings.Join(ret, ",")
}

func (p sizeList) Len() int           { return len(p) }
func (p sizeList) Less(i, j int) bool { return p[i].val < p[j].val }
func (p sizeList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func getPasswordSize(cmd *cobra.Command) (int, error) {
	allSizes := availableSizes()

	key, _ := cmd.Flags().GetString(optSize)
	for _, x := range allSizes {
		if strings.EqualFold(x.key, key) {
			return x.val, nil
		}
	}

	return -1, fmt.Errorf("invalid size <%s>, allowed values are [%s]", key, allSizes)
}

func availableSizes() sizeList {
	ret := make(sizeList, 6)
	ret[0] = size{"S", 8}
	ret[1] = size{"M", 12}
	ret[2] = size{"L", 16}
	ret[3] = size{"XL", 32}
	ret[4] = size{"XXL", 64}
	ret[5] = size{"XXXL", 128}

	sort.Sort(ret)

	return ret
}
