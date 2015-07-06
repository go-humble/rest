// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package rest

// DefaultId is a struct with an Id property and a getter
// called ModelId. You can embed it to satisfy the ModelId
// method of rest.Model.
type DefaultId struct {
	Id string
}

// ModelId satisfies the ModelId method of rest.Model.
func (d DefaultId) ModelId() string {
	return d.Id
}
