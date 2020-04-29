package jwt

import (
	"fmt"
	"testing"
)

func TestAudiences(t *testing.T) {
	testTable := []struct {
		name string
		aud  string
		err  error
	}{
		{
			name: "misc: not enough slashes",
			aud:  "/projects/1234",
			err:  fmt.Errorf("audience \"/projects/1234\" must follow the format \"/projects/PROJECT_NUMBER/SERVICE_DETAILS\""),
		},
		{
			name: "app engine: valid",
			aud:  "/projects/1234/apps/fake-project-id",
			err:  nil,
		},
		{
			name: "app engine: missing service details",
			aud:  "/projects/1234/",
			err:  fmt.Errorf("audience \"/projects/1234/\" is missing service details"),
		},
		{
			name: "global: valid",
			aud:  "/projects/1234/global/backendServices/1234",
			err:  nil,
		},
		{
			name: "global: missing service details",
			aud:  "/projects/1234/",
			err:  fmt.Errorf("audience \"/projects/1234/\" is missing service details"),
		},
	}

	for _, tc := range testTable {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseAudience(tc.aud)

			switch {
			case err == nil && tc.err == nil:
				// noop
			case err != nil && tc.err == nil:
				t.Error("expected no error, got error:", err)
			case err == nil && tc.err != nil:
				t.Error("expected error, got no error:", tc.err)
			case err != nil && tc.err != nil:
				if err.Error() != tc.err.Error() {
					t.Error("unexpected error got: ", err)
					t.Error("expected error      : ", tc.err)
				}
			}

		})
	}
}
