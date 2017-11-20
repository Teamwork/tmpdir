package tmpdir

import (
	"os"
	"strings"
	"testing"

	"github.com/teamwork/test"
)

func TestMkTemp(t *testing.T) {
	err := Setup("tmpdir_test")
	if err != nil {
		t.Fatalf("Setup error: %v", err)
	}
	if _, err := os.Stat(Dir); err != nil {
		t.Fatalf("Dir was not created: %v", err)
	}

	defer func() {
		err := Cleanup()
		if err != nil {
			t.Errorf("Cleanup error: %v", err)
		}
		if _, err := os.Stat(Dir); err == nil {
			t.Fatalf("Dir was not cleaned up: %v", err)
		}
	}()

	cases := []struct {
		in, want, wantErr string
	}{
		{"hello.pdf", "hello.pdf", ""},
		{"hello.pdf", "hello-1.pdf", ""},
		{"hello/world.pdf", "hello-world.pdf", ""},
		{"../../../xxx", "..-..-..-xxx", ""},
		{"//xxx", "xxx", ""},
		{strings.Repeat("x", 300), strings.Repeat("x", 250-len(Dir)-1), ""},
		{strings.Repeat("x", 250) + ".ext", strings.Repeat("x", 250-len(Dir)-4) + ".ext", ""},
		{strings.Repeat("x", 250) + ".ext", strings.Repeat("x", 250-len(Dir)-4) + "-1.ext", ""},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			out, err := MkTemp(tc.in)

			tc.want = Dir + tc.want

			if !test.ErrorContains(err, tc.wantErr) {
				t.Fatalf("wrong error\nout:  %v\nwant: %#v\n", err, tc.wantErr)
			}
			if out.Name() != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out.Name(), tc.want)
			}
		})
	}
}

func TestSplitExt(t *testing.T) {
	cases := []struct {
		in, wantBase, wantExt string
	}{
		{"hello.pdf", "hello", "pdf"},
		{"/path/hello.pdf", "/path/hello", "pdf"},
		{"/path.hello.pdf", "/path.hello", "pdf"},
		{".pdf", "", "pdf"},
		{"hello.", "hello", ""},
		{"hello", "hello", ""},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			base, ext := splitExt(tc.in)
			if base != tc.wantBase {
				t.Errorf("\nout:  %#v\nwant: %#v\n", base, tc.wantBase)
			}
			if ext != tc.wantExt {
				t.Errorf("\nout:  %#v\nwant: %#v\n", ext, tc.wantExt)
			}
		})
	}
}