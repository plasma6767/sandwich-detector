package ethnode

import "testing"

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
