package core_test

import (
	"bytes"
	"testing"

	"github.com/Petroviiic/GoGit/internal/core"
)

func TestParseTree(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		content []byte
		want    *core.Tree
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := core.ParseTree(tt.content)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseTree() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseTree() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ParseTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserialize(t *testing.T) {
	tests := []struct {
		name    string
		input   *core.BaseObject
		wantErr bool
	}{
		{
			name: "Blob test",
			input: &core.BaseObject{
				Type:    "blob",
				Content: []byte("zdravo ovo sam ja"),
			},
			wantErr: false,
		},
		{
			name: "Tree test",
			input: &core.BaseObject{
				Type:    "tree",
				Content: []byte(`[{"Mode":"100644","Name":"file.txt","Hash":"abc"},{"Mode":"100644","Name":"file1.txt","Hash":"abcd"},{"Mode":"040000","Name":"folder","Hash":"abcde"}]`),
			},
			wantErr: false,
		},
		{
			name:    "Error case - invalid data",
			input:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			if tt.input != nil {
				data, err = tt.input.Serialize()
				if err != nil {
					t.Fatalf("Failed to serialize for test: %v", err)
				}
			} else {
				data = []byte("random data, should fail")
			}

			got, err := core.Deserialize(data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.GetType() != tt.input.Type {
				t.Errorf("Type mismatch: got %v, want %v", got.GetType(), tt.input.Type)
			}

			if !bytes.Equal(got.GetContent(), tt.input.Content) {
				t.Errorf("Content mismatch: got %s, want %s", string(got.GetContent()), string(tt.input.Content))
			}
		})
	}
}
