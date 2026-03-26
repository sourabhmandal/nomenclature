package translation

import "testing"

func TestDetectInputType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "single string",
			input: "hello",
			want:  map[string]string{"hello": "string"},
		},
		{
			name:  "string number link",
			input: "hello 123 https://example.com",
			want: map[string]string{
				"hello":               "string",
				"123":                 "number",
				"https://example.com": "link",
			},
		},
		{
			name:  "boolean and number",
			input: "true 90",
			want: map[string]string{
				"true": "boolean",
				"90":   "number",
			},
		},
		{
			name:  "repeated tokens collapse to single map key",
			input: "hi there 10 20",
			want: map[string]string{
				"hi":    "string",
				"there": "string",
				"10":    "number",
				"20":    "number",
			},
		},
		{
			name:  "email and date",
			input: "dev@example.com 2026-03-26",
			want: map[string]string{
				"dev@example.com": "email",
				"2026-03-26":      "date",
			},
		},
		{
			name:  "empty input",
			input: "   ",
			want:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectInputType(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("DetectInputType(%q) length = %d, want %d", tt.input, len(got), len(tt.want))
			}

			for token, wantType := range tt.want {
				gotType, ok := got[token]
				if !ok {
					t.Fatalf("DetectInputType(%q) missing token %q", tt.input, token)
				}
				if gotType != wantType {
					t.Fatalf("DetectInputType(%q) token %q type = %q, want %q", tt.input, token, gotType, wantType)
				}
			}
		})
	}
}
