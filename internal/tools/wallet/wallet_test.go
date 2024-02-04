package wallet

import "testing"

func TestCreateWallet(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CreateWallet()
			if got != tt.want {
				t.Errorf("CreateWallet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CreateWallet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
