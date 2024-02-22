package apieastmoney

import (
	"strings"

	"example.com/stocker-back/internal/infra"
)

type APIServiceEastmoney struct {
	logger infra.Logger
}

func NewAPIServiceEastmoney(logger infra.Logger) *APIServiceEastmoney {
	return &APIServiceEastmoney{
		logger: logger,
	}
}

// sliceStringByChar slice input string by startChar and endChar if they are valid.
func sliceStringByChar(input, startChar, endChar string) string {
	startIndex := strings.Index(input, startChar)
	if startIndex == -1 {
		return ""
	}

	endIndex := strings.LastIndex(input, endChar)
	if endIndex == -1 {
		return ""
	}

	return input[startIndex+1 : endIndex]
}
