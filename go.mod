module github.com/creachadair/jrpc2

require (
	github.com/fortytw2/leaktest v1.3.0
	github.com/google/go-cmp v0.5.9
	golang.org/x/sync v0.2.0
)

require github.com/creachadair/mds v0.0.1

go 1.19

// A bug in handler.New could panic a wrapped handler on pointer arguments.
retract [v0.21.2, v0.22.0]

// Checksum mismatch due to accidental double tag push. Safe to use, but warns.
retract v0.23.0
