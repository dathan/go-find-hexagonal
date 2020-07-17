package filesystem

import (
	"testing"

	"github.com/dathan/go-find-hexagonal/pkg/filter"
	"github.com/dathan/go-find-hexagonal/pkg/find"
)

func Test_filesytemRepository_Find(t *testing.T) {
	type args struct {
		fo find.FilterOptions
	}

	fileSystemFilter := filter.NewGenericOptions("../../")

	fFunc := find.FilterFunc(func(fr *find.FindResult) bool {
		if fr != nil && fr.Name == "filesystem.go" {
			return true
		}
		return false
	})

	fileSystemFilter.SetFilterFunc(fFunc)

	testWalk := args{
		fo: fileSystemFilter,
	}

	tests := []struct {
		name    string
		f       FilesytemRepository
		args    args
		want    int
		wantErr bool
	}{
		{"currentfiles", FilesytemRepository{}, testWalk, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Find(tt.args.fo)
			if (err != nil) != tt.wantErr {
				t.Errorf("filesytemRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !(len(got) == tt.want) {
				t.Errorf("filesytemRepository.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
