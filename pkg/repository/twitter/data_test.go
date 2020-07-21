// log into twitter, gather all likes,
// store likes in a datastore,
// search the datastore for link content when implementing the grep
// note the difference between find and search
package twitter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dathan/go-find-hexagonal/pkg/filter"
	"github.com/dathan/go-find-hexagonal/pkg/find"
)

func TestRepository_Find(t *testing.T) {

	var fFunc find.FilterFunc = func(fr *find.FindResult) bool {
		if fr != nil && strings.Contains(fr.Name, "golang") {
			//fmt.Printf("fr.Name: %s contains golang\n", fr.Name)
			return true
		}
		fmt.Printf("fr.Name: %s contains golang\n", fr.Name)

		return false
	}

	// set the filter func to find the specific files
	var filterInfo find.FilterOptions
	filterInfo = filter.NewGenericOptions(".").SetFilterFunc(fFunc)

	type args struct {
		fo find.FilterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    find.FindResults
		wantErr bool
	}{
		{"default", args{filterInfo}, find.FindResults{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Repository{}
			got, err := f.Find(tt.args.fo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
