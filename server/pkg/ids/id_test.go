package ids

import "testing"

func testId(n int, f func() string) map[string]int {
	m := make(map[string]int)
	for range n {
		m[f()]++
	}

	var repeat map[string]int
	for k, v := range m {
		if v > 1 {
			repeat[k] = v
		}
	}
	return repeat
}

func TestULIdRepeat(t *testing.T) {
	n := 100_0000
	ids := testId(n, ULID)
	for id, r := range ids {
		t.Logf("ULID %s conflict %d times in %d times", id, r, n)
	}
}

func TestUUIdRepeat(t *testing.T) {
	n := 100_0000
	ids := testId(n, UUID)
	for id, r := range ids {
		t.Logf("uuid %s conflict %d times in %d times", id, r, n)
	}
}

func TestUUID(t *testing.T) {
	uuid := UUID()
	t.Log(uuid)
}

func TestULID(t *testing.T) {
	ulid := ULID()
	t.Log(ulid)
}
