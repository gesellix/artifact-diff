package artifact_diff

import (
	"reflect"
	"testing"
)

func TestArtifactInfo_WithFlattenedAndSortedFileInfos(t *testing.T) {
	type fields struct {
		Path      string
		Count     int
		FileInfos FileInfos
	}
	fileInfos := FileInfos{}
	fileInfos["2"] = &FileInfo{
		"zip!/dir/com.example.foo.stuff.jar",
		22,
		"222",
	}
	fileInfos["4"] = &FileInfo{
		"zip!/dir/com.example.foo.bar_1.2.3.4711.jar",
		44,
		"444",
	}
	fileInfos["1"] = &FileInfo{
		"zip!/dir/com.example.foo.conf.jar",
		11,
		"111",
	}
	fileInfos["3"] = &FileInfo{
		"zip!/dir/com.example.foo.jar",
		33,
		"333",
	}
	fileInfos["6"] = &FileInfo{
		"zip!/dir/com.example.foo-bar.jar",
		66,
		"666",
	}
	fileInfos["5"] = &FileInfo{
		"zip!/dir/com.example.foo.bar.jar",
		55,
		"555",
	}
	expectedFileInfos := []*FileInfo{
		{
			"zip!/dir/com.example.foo-bar.jar",
			66,
			"666",
		},
		{
			"zip!/dir/com.example.foo.bar.jar",
			55,
			"555",
		},
		{
			"zip!/dir/com.example.foo.bar_1.2.3.4711.jar",
			44,
			"444",
		},
		{
			"zip!/dir/com.example.foo.conf.jar",
			11,
			"111",
		},
		{
			"zip!/dir/com.example.foo.jar",
			33,
			"333",
		},
		{
			"zip!/dir/com.example.foo.stuff.jar",
			22,
			"222",
		},
	}
	tests := []struct {
		name   string
		fields fields
		want   *FlatArtifactInfo
	}{
		{
			"ensure sorted fileInfos",
			fields{
				"a/path",
				0,
				fileInfos,
			},
			&FlatArtifactInfo{
				"a/path",
				0,
				expectedFileInfos,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ArtifactInfo{
				Path:      tt.fields.Path,
				Count:     tt.fields.Count,
				FileInfos: tt.fields.FileInfos,
			}
			if got := d.WithFlattenedAndSortedFileInfos(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithFlattenedAndSortedFileInfos() = %v, want %v", got, tt.want)
			}
		})
	}
}
