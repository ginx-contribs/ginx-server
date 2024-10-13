package captcha

import (
	"math/rand/v2"
	"strings"
)

// GenCaptcha generates captcha with given length, the recommended n should be greater than 6.
func GenCaptcha(n int) string {
	var w strings.Builder
	for range n {
		if rand.Int()%2 == 1 {
			w.WriteByte('0' + byte(rand.IntN(10)))
		} else {
			w.WriteByte('A' + byte(rand.IntN(26)))
		}
	}
	return w.String()
}
