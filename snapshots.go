/*
BSD 2-Clause License

Copyright (c) 2014, Alexander Neumann <alexander@bumpern.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package restic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Snapshot is the state of a resource at one point in time.
type Snapshot struct {
	ID       *ID       `json:"id"`
	ShortID  string    `json:"short_id"`
	Time     time.Time `json:"time"`
	Parent   *ID       `json:"parent,omitempty"`
	Tree     *ID       `json:"tree"`
	Paths    []string  `json:"paths"`
	Hostname string    `json:"hostname,omitempty"`
	Username string    `json:"username,omitempty"`
	UID      uint32    `json:"uid,omitempty"`
	GID      uint32    `json:"gid,omitempty"`
	Excludes []string  `json:"excludes,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Original *ID       `json:"original,omitempty"`

	ProgramVersion string `json:"program_version,omitempty"`
}

// idSize contains the size of an ID, in bytes.
const idSize = sha256.Size

// ID references content within a repository.
type ID [idSize]byte

// ParseID converts the given string to an ID.
func ParseID(s string) (ID, error) {
	if len(s) != hex.EncodedLen(idSize) {
		return ID{}, fmt.Errorf("invalid length for ID: %q", s)
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return ID{}, fmt.Errorf("invalid ID: %s", err)
	}

	id := ID{}
	copy(id[:], b)

	return id, nil
}

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

// MarshalJSON returns the JSON encoding of id.
func (id ID) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 2+hex.EncodedLen(len(id)))

	buf[0] = '"'
	hex.Encode(buf[1:], id[:])
	buf[len(buf)-1] = '"'

	return buf, nil
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in id.
func (id *ID) UnmarshalJSON(b []byte) error {
	// check string length
	if len(b) != len(`""`)+hex.EncodedLen(idSize) {
		return fmt.Errorf("invalid length for ID: %q", b)
	}

	if b[0] != '"' {
		return fmt.Errorf("invalid start of string: %q", b[0])
	}

	// Strip JSON string delimiters. The json.Unmarshaler contract says we get
	// a valid JSON value, so we don't need to check that b[len(b)-1] == '"'.
	b = b[1 : len(b)-1]

	_, err := hex.Decode(id[:], b)
	if err != nil {
		return fmt.Errorf("invalid ID: %s", err)
	}

	return nil
}
