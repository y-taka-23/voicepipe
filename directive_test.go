package main

import (
	"testing"
)

func TestNewDirective(t *testing.T) {
	in := "stub/voicepipe_test.yml"
	want := &Directive{
		Repository: "user/app",
		ImageDirectives: []*ImageDirective{
			&ImageDirective{
				Tag:         "tag_1",
				Description: "description 1",
				Parameters: map[string]string{
					"NAME_1": "value_1",
				},
			},
			&ImageDirective{
				Tag:         "tag_2",
				Description: "description 2",
				Parameters: map[string]string{
					"NAME_2": "value_2",
				},
			},
		},
	}
	got, err := NewDirective(in)
	if err != nil {
		t.Errorf("NewDirective(%s) returns %s", in, err)
	}
	if got.Repository != want.Repository {
		t.Errorf(
			"NewDirective(%s).Repository == %s, want %s",
			in, got.Repository, want.Repository,
		)
	}
	for i, id := range got.ImageDirectives {
		wantedID := want.ImageDirectives[i]
		if id.Tag != wantedID.Tag {
			t.Errorf(
				"NewDirective(%s).ImageDirectives[%d].Tag == %s, want %s",
				in, i, id.Tag, wantedID.Tag,
			)
		}
		if id.Description != wantedID.Description {
			t.Errorf(
				"NewDirective(%s).ImageDirectives[%d].Description == %s, want %s",
				in, i, id.Description, wantedID.Description,
			)
		}
		if len(id.Parameters) != len(wantedID.Parameters) {
			t.Errorf(
				"the number of NewDirective(%s).ImageDirectives[%d].Parameters is %d, want %d",
				in, i, len(id.Parameters), len(wantedID.Parameters),
			)
		}
		for k, v := range id.Parameters {
			if v != wantedID.Parameters[k] {
				t.Errorf(
					"NewDirective(%s).ImageDirectives[%d].Parameter[%s] == %s, want %s",
					in, i, k, v, wantedID.Parameters[k],
				)
			}
		}
	}
}
