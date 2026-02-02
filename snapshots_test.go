package restic

import (
	"encoding/json"
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid 64-char hex ID",
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: false,
		},
		{
			name:    "invalid length - too short",
			input:   "0123456789abcdef",
			wantErr: true,
		},
		{
			name:    "invalid length - too long",
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef00",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "0123456789abcdefghij0123456789abcdef0123456789abcdef0123456789abcd",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.input {
					t.Errorf("ID.String() = %v, want %v", got.String(), tt.input)
				}
			}
		})
	}
}

func TestID_String(t *testing.T) {
	input := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	id, err := ParseID(input)
	if err != nil {
		t.Fatalf("ParseID() failed: %v", err)
	}

	got := id.String()
	if got != input {
		t.Errorf("ID.String() = %v, want %v", got, input)
	}
}

func TestID_MarshalJSON(t *testing.T) {
	input := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	id, err := ParseID(input)
	if err != nil {
		t.Fatalf("ParseID() failed: %v", err)
	}

	got, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() failed: %v", err)
	}

	want := `"` + input + `"`
	if string(got) != want {
		t.Errorf("MarshalJSON() = %v, want %v", string(got), want)
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid JSON ID",
			input:   `"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"`,
			wantErr: false,
		},
		{
			name:    "invalid length",
			input:   `"0123456789abcdef"`,
			wantErr: true,
		},
		{
			name:    "missing quotes",
			input:   `0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef`,
			wantErr: true,
		},
		{
			name:    "invalid hex",
			input:   `"0123456789abcdefghij0123456789abcdef0123456789abcdef01234567890"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := id.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestID_JSONRoundTrip(t *testing.T) {
	original := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	id1, err := ParseID(original)
	if err != nil {
		t.Fatalf("ParseID() failed: %v", err)
	}

	data, err := json.Marshal(id1)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	var id2 ID
	if err := json.Unmarshal(data, &id2); err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if id1 != id2 {
		t.Errorf("Round-trip failed: got %v, want %v", id2, id1)
	}
}

func TestSnapshot_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"short_id": "01234567",
		"time": "2024-01-01T12:00:00Z",
		"tree": "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210",
		"paths": ["/home/user"],
		"hostname": "testhost",
		"username": "testuser",
		"uid": 1000,
		"gid": 1000,
		"tags": ["daily", "important"],
		"program_version": "restic 0.16.0"
	}`

	var snapshot Snapshot
	err := json.Unmarshal([]byte(jsonData), &snapshot)
	if err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if snapshot.ShortID != "01234567" {
		t.Errorf("ShortID = %v, want %v", snapshot.ShortID, "01234567")
	}
	if snapshot.Hostname != "testhost" {
		t.Errorf("Hostname = %v, want %v", snapshot.Hostname, "testhost")
	}
	if len(snapshot.Tags) != 2 {
		t.Errorf("len(Tags) = %v, want %v", len(snapshot.Tags), 2)
	}
	if len(snapshot.Paths) != 1 || snapshot.Paths[0] != "/home/user" {
		t.Errorf("Paths = %v, want [/home/user]", snapshot.Paths)
	}
}
