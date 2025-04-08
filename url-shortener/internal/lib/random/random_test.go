package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			assert.Len(t, str1, tt.size)

			str2 := NewRandomString(tt.size)
			assert.Len(t, str2, tt.size)

			// При быстром вызове тестов время может совпасть, из-за этого значения могут быть одинаковы.
			//Поэтому после NewRandomString() поставил другую функцию, чтобы искусственно
			//создать задержку времени.
			assert.Equal(t, str1, str2)
		})
	}
}
