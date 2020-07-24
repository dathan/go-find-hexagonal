package cache

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewService(t *testing.T) {

	store, err := NewService()
	if err != nil {
		t.Errorf("Error with NewService(): %s", err)
		t.Fail()
	}

	tests := []struct {
		name    string
		want    Store
		wantErr bool
	}{
		{"newstore", store, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = got.Set("defaultInfo", tests)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tInfo := []struct {
				name    string
				want    Store
				wantErr bool
			}{}

			found, err := got.Get("defaultInfo", tInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Did TInfoWork: %v %s\n", found, spew.Sdump(tInfo))
		})
	}
}
