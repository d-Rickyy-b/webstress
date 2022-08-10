package util

import "testing"

func Test_numOfChars(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "Test Zero",
			input: 0,
			want:  1,
		},
		{
			name:  "Test 1 digit",
			input: 4,
			want:  1,
		},
		{
			name:  "Test 2 digits",
			input: 25,
			want:  2,
		},
		{
			name:  "Test 3 digits",
			input: 481,
			want:  3,
		},
		{
			name:  "Test 4 digits",
			input: 8481,
			want:  4,
		},
		{
			name:  "Test 5 digits",
			input: 18481,
			want:  5,
		},
		{
			name:  "Test 6 digits",
			input: 180481,
			want:  6,
		},
		{
			name:  "Test negative 1 digit",
			input: -3,
			want:  2,
		},
		{
			name:  "Test negative 2 digits",
			input: -12,
			want:  3,
		},
		{
			name:  "Test negative 3 digits",
			input: -123,
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumOfChars(tt.input); got != tt.want {
				t.Errorf("NumOfChars() = %v, want %v", got, tt.want)
			}
		})
	}
}
