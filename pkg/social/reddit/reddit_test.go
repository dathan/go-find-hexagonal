package reddit

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestUpVotes_List(t *testing.T) {
	type args struct {
		handle string
	}
	tests := []struct {
		name string
		args args
	}{
		{"testlist", args{"dathanvp"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewService()
			if err != nil {
				t.Errorf("ERROR: %s\n", err)
			}

			upv, err := client.Favorites.CacheList(tt.args.handle)
			if err != nil {
				t.Errorf("ERROR: %s\n", err)
			}

			spew.Dump(upv)

		})
	}
}
