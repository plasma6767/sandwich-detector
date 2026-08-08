// Tests for hostOnly(). These don't make any network calls, so they run
// fast and safely as part of `make test`.
package ethnode

import "testing"

// TestHostOnly checks hostOnly against several example URLs, using a table
// test: a list of input/expected-output cases looped over, instead of one
// function per case.
func TestHostOnly(t *testing.T) {
	cases := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{
			name:   "strips path and secret key, keeps host",
			rawURL: "https://eth-mainnet.g.alchemy.com/v2/super-secret-key",
			want:   "eth-mainnet.g.alchemy.com",
		},
		{
			name:   "keeps port if present",
			rawURL: "http://localhost:8545",
			want:   "localhost:8545",
		},
		{
			name:    "errors on malformed URL",
			rawURL:  "://not-a-valid-url",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		// t.Run creates a named sub-test per case, so a failure points
		// directly at which case broke.
		t.Run(tc.name, func(t *testing.T) {
			got, err := hostOnly(tc.rawURL)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got host %q, want %q", got, tc.want)
			}
		})
	}
}
