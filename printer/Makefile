# Copyright 2013 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: zapi.go

zapi.go: printer.go
	GOOS=windows ./mksyscall_windows.pl $^ \
		| gofmt \
		> $@
