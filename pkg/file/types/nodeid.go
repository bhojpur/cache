package types

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/bhojpur/cache/pkg/file/crypto"
)

var (
	// ErrInvalidIP is an error thrown when the ip to construct a NodeID is invalid.
	ErrInvalidIP = errors.New("IP address to construct a NodeID is invalid")

	// ErrInvalidPort is an error thrown when the port to construct a NodeID is invalid.
	ErrInvalidPort = errors.New("port to construct a NodeID is invalid")
)

// NodeID contains the necessary data for referencing and connecting to a node.
type NodeID struct {
	IP   string       `json:"ip"`   // The node's IP address
	Port int          `json:"port"` // The port on which the node is hosted
	Hash *crypto.Hash `json:"hash"` // The hash of the node
}

// NewNodeID constructs a new NodeID.
func NewNodeID(ip string, port int) (*NodeID, error) {
	// Check that the IP address is valid
	addr := net.ParseIP(ip)
	if addr == nil {
		return nil, ErrInvalidIP
	}

	// Check that the port is valid
	if port == 0 {
		return nil, ErrInvalidPort
	}

	// Create the NodeID
	newNodeID := &NodeID{
		IP:   ip,
		Port: port,
		Hash: nil,
	}

	hash := crypto.Sha3(newNodeID.Bytes())
	newNodeID.Hash = &hash

	return newNodeID, nil
}

/* ----- BEGIN HELPER FUNCTIONS ----- */

// Bytes returns the bytes of a NodeID.
func (nodeID *NodeID) Bytes() []byte {
	json, _ := json.MarshalIndent(*nodeID, "", "  ")
	return json
}

// String converts the NodeID to a string.
func (nodeID *NodeID) String() string {
	json, _ := json.MarshalIndent(*nodeID, "", "  ")
	return string(json)
}

/* ----- END HELPER FUNCTIONS ----- */
