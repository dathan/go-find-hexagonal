package cache

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type TestThing struct {
	Name string
}

func TestNewService(t *testing.T) {

	store, err := NewService()
	if err != nil {
		t.Errorf("Error with NewService(): %s", err)
		t.Fail()
	}

	testStruct := TestThing{"testStruct"}

	tests := []struct {
		name    string
		want    Store
		wantErr bool
	}{
		{"newstore", store, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := store
			err := got.Set("defaultInfoxxx", testStruct)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var tInfo TestThing
			found, err := got.Get("defaultInfoxxx", &tInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Did TInfoWork: %v %s\n", found, spew.Sdump(tInfo))
		})
	}
}
