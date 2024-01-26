package digest

import (
	"testing"
)

func TestSHA256(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{"Empty string", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"Example string", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := SHA256(tc.input)
			if output != tc.expectedOutput {
				t.Errorf("Expected %s, but got %s", tc.expectedOutput, output)
			}
		})
	}
}
