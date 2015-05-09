// Copyright 2014 The fav Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

type Program struct {
	//0:service, 1:command, 2:application(GUI), run at new window
	Type int

	//program path
	Path string

	//arguments
	Args []string

	//PID
	PID int64
}
