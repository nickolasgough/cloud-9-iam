package utils

import (
	"fmt"
	"strings"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/constants"
)

func BuildClientURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	return fmt.Sprintf("%s%s", constants.CLIENT_BASE_URL, path)
}
