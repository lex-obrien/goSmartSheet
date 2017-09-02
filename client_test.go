package goSmartSheet

import "testing"

func Test_validateURL(t *testing.T) {
	tests := []struct {
		name        string
		u           string
		wantIsValid bool
		wantErr     bool
	}{
		{name: "1", u: "http://www.test.com/foo", wantIsValid: true, wantErr: false},
		{name: "2", u: "https://www.test.com/foo", wantIsValid: true, wantErr: false},
		{name: "3", u: "https://wasdasdww.test.coddm/foo", wantIsValid: true, wantErr: false},
		{name: "4", u: "4432", wantIsValid: false, wantErr: true},
		{name: "5", u: "www.test.com", wantIsValid: false, wantErr: true},
		{name: "6", u: "www.test.com/foo", wantIsValid: false, wantErr: true},
		{name: "7", u: "http://www.test.co", wantIsValid: false, wantErr: true},
		{name: "8", u: "", wantIsValid: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsValid, err := validateURL(tt.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("validateURL() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}
