package user

import "testing"

func TestPasswordHash(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "Simple password",
			password: "password123",
			expected: "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
		},
		{
			name:     "Another password",
			password: "helloWorld",
			expected: "11d4ddc357e0822968dbfd226b6e1c2aac018d076a54da4f65e1dc8180684ac3",
		},
		{
			name:     "Empty password",
			password: "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PasswordHash(tc.password)
			if got != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func BenchmarkPasswordHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PasswordHash("11d4ddc357e0822968dbfd226b6e1c2aac018d076a54da4f65e1dc8180684ac3")
	}
}
