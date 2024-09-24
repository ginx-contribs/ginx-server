package email

import (
	"github.com/matcornic/hermes/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGeneratedAuthCode(t *testing.T) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      "ginx",
			Copyright: " ",
		},
	}

	html, err := h.GenerateHTML(TmplConfirmCode("sign up", "jack", "Y2635H6", time.Minute*5))
	assert.Nil(t, err)
	err = os.WriteFile("testdata/auth_code.html", []byte(html), 0666)
	assert.Nil(t, err)
}
