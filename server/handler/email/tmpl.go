package email

import (
	"fmt"
	"github.com/matcornic/hermes/v2"
	"time"
)

// TmplConfirmCode template for confirm code
func TmplConfirmCode(usage, name, code string, ttl time.Duration) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				fmt.Sprintf("Welcome, you are applying for a verification code to %s.", usage),
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf("Please use your code as quickly, it will be expired in %s.", ttl.String()),
					InviteCode:   code,
				},
			},
			Outros: []string{
				"If it is not your own operation, please ignore this email. If you need help, access https://github.com/ginx-contribs for more information.",
			},
		},
	}
}
