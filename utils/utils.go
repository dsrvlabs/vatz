package utils

import (
	"strconv"
)

func MakeUniqueValue(pName, pAddr string, pPort int) string {
	return pName + pAddr + strconv.Itoa(pPort)
}
