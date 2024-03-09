#!/bin/sh

fmtg() {
	find \
		. \
		-type f \
		-name '*.go' |
		xargs \
			gofmt \
			-s \
			-w
}

fmts() {
	which shfmt &&
	find \
		. \
		-type f \
		-name '*.sh' |
		xargs \
			-r \
			shfmt --write
}

fmtg
fmts
