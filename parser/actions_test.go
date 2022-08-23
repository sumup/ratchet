package parser

import (
	"reflect"
	"testing"
)

func TestActions_Parse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		exp  []string
	}{
		{
			name: "mostly_empty_file",
			in: `
runs:
`,
			exp: []string{},
		},
		{
			name: "uses",
			in: `
runs:
  steps:
    - uses: 'actions/checkout@v3'
    - name: test step with ubuntu 
      uses: 'docker://ubuntu:20.04'
    - uses: 'docker://ubuntu@sha256:47f14534bda344d9fe6ffd6effb95eefe579f4be0d508b7445cf77f61a0e5724'
      with:
        uses: 'foo/bar@v0'
    - uses: './github/workflows/other.yml'
    - uses: 'org/repo/.github/workflows/other@v0'
    - name: test step with shell
      id: shell_test
      shell: bash
      run: echo "uses"
`,
			exp: []string{
				"actions://actions/checkout@v3",
				"actions://org/repo/.github/workflows/other@v0",
				"container://ubuntu:20.04",
				"container://ubuntu@sha256:47f14534bda344d9fe6ffd6effb95eefe579f4be0d508b7445cf77f61a0e5724",
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := helperStringToYAML(t, tc.in)

			refs, err := new(Actions).Parse(m)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := refs.Refs(), tc.exp; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %q to be %q", got, want)
			}
		})
	}
}
