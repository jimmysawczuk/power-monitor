package tmplfunc

import (
	"os"
)

func Env(s string) string {
	return os.Getenv(s)
}
