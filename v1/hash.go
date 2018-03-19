// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Hash is an unqualified digest of some content, e.g. sha256:deadbeef
type Hash struct {
	algorithm string
	hex       string
}

// Algorithm returns the algorithm used to compute the hash.
func (h Hash) Algorithm() string {
	return h.algorithm
}

// Hex returns the hex portion of the content hash.
func (h Hash) Hex() string {
	return h.hex
}

// String reverses NewHash returning the string-form of the hash.
func (h Hash) String() string {
	return fmt.Sprintf("%s:%s", h.Algorithm(), h.Hex())
}

// NewHash validates the input string is a hash and returns a strongly type Hash object.
func NewHash(s string) (Hash, error) {
	h := Hash{}
	if err := h.parse(s); err != nil {
		return Hash{}, err
	}
	return h, nil
}

// MarshalJSON implements json.Marshaler
func (h *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (h *Hash) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return h.parse(s)
}

func (h *Hash) parse(unquoted string) error {
	parts := strings.Split(unquoted, ":")
	if len(parts) != 2 {
		return fmt.Errorf("too many parts in hash: %s", unquoted)
	}

	rest := strings.TrimLeft(parts[1], "0123456789abcdef")
	if len(rest) != 0 {
		return fmt.Errorf("found non-hex character in hash: %c", rest[0])
	}

	switch parts[0] {
	case "sha256":
		if len(parts[1]) != 64 {
			return fmt.Errorf("wrong number of hex digits for sha256: %s", parts[1])
		}
	default:
		return fmt.Errorf("unsupported hash type: %s", parts[0])
	}

	h.algorithm = parts[0]
	h.hex = parts[1]
	return nil
}

// SHA256 computes the Hash of the provided io.ReadCloser's content.
func SHA256(r io.ReadCloser) (Hash, int64, error) {
	defer r.Close()
	hasher := sha256.New()
	n, err := io.Copy(hasher, r)
	if err != nil {
		return Hash{}, 0, err
	}
	return Hash{
		algorithm: "sha256",
		hex:       hex.EncodeToString(hasher.Sum(make([]byte, 0, hasher.Size()))),
	}, n, nil
}
