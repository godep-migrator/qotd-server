This is a QOTD server written as a small project for exploring Go.

> You can get *anywhere* in ten minutes if you drive fast enough.

## Try it out

- Install the server: `go get -u github.com/nixterrimus/go-qotd`
- Start the server: `go-qotd wisdom.txt`
- In another terminal make a request: `nc localhost 3333`

## Quote File Format

This QOTD server has a mandatory files argument of a quote file, for
example to run with the `wisdom.txt` file included in the git repository
you would use: `go-qotd wisdom.txt`.

The quote file format is exactly the same as the [fortune file
format](http://en.wikipedia.org/wiki/Fortune_(Unix)#Fortune_files), that
is to say:

- It is a text file
- Each quote is separated by a `%` character of its own line

## The QOTD Protocol

RFC 865 defines the [Quote of the Day
Protocol](http://tools.ietf.org/html/rfc865). The specification is
really short.  Here's the gist of it:

- A TCP port is opened on port 17
- On connection a quote is served
- The connection is closed
- The service may also listen on UDP
- Quotes should be less than 512 characters long

This server is not yet RFC 865 compliant.  But I'm working on it.

## Next Steps / Project Goals

- [X] Accept file as CLI argument
- [ ] Accept Quotes via Standard In
- [ ] Accept a URL for a quote file
- [X] Accept port as CLI argument
- [X] Server should run in 865 compliant or non-compliant mode (512
character truncation)
- [ ] Server should listen on UDP
- [ ] Server should be well tested
- [ ] Installable through homebrew
- [ ] Advertise QOTD over MDNS

## Client Goals

- [ ] Find Server over MDNS
- [ ] Allow `host port` as args

## Install through homebrew

Notes:
  - https://github.com/Homebrew/homebrew-binary
  - https://github.com/Homebrew/homebrew/pull/20895
  - https://github.com/Homebrew/homebrew/pull/21810/files
  - https://github.com/Homebrew/homebrew/issues/23703

## Author

This was written by [Nick Rowe](http://dcxn.com) in the summer of
2014 with guidance by Rafe Colton.
