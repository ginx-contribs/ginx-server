package captcha

import "testing"

func TestGenCaptcha(t *testing.T) {
	captcha := GenCaptcha(8)
	t.Log(captcha)
}

func testGenCaptchaConflicts(l int, n int) int64 {
	dict := make(map[string]int)
	for range n {
		dict[GenCaptcha(l)]++
	}
	var conflict int64
	for _, v := range dict {
		if v > 1 {
			conflict++
		}
	}
	return conflict
}

func TestGenCaptcha_Conflict(t *testing.T) {
	samples := []struct {
		l, n int
	}{
		{4, 10_000},
		{6, 10_000},
		{8, 10_000},
		{4, 100_000},
		{6, 100_000},
		{8, 100_000},
		{4, 1_000_000},
		{6, 1_000_000},
		{8, 1_000_000},
	}

	for _, sample := range samples {
		conflicts := testGenCaptchaConflicts(sample.l, sample.n)
		t.Logf("len: %d\tn: %10d\tconflicts: %10d", sample.l, sample.n, conflicts)
	}
}
