package strings

import (
	"strconv"
	"strings"
)

type IPS []string

func (ip IPS) Len() int { return len(ip) }

func (ip IPS) Swap(i, j int) { ip[i], ip[j] = ip[j], ip[i] }

func (ip IPS) Less(i, j int) bool {
	address1 := strings.Split(ip[i], ".")
	address2 := strings.Split(ip[j], ".")

	var result bool
	for k := 0; k < 4; k++ {
		if address1[k] != address2[k] {
			num1, _ := strconv.Atoi(address1[k])
			num2, _ := strconv.Atoi(address2[k])
			result = num1 < num2
			break
		}
	}

	return result
}
