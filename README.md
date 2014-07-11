This is a QOTD server written as a small project for exploring Go.

> You can get *anywhere* in ten minutes if you drive fast enough.

## Try it out

- Install the server: `go get -u github.com/nixterrimus/go-qotd`
- Start the server: `go-qotd`
- In another terminal make a request: `nc localhost 3333`

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

- [ ] Accept file as CLI argument
- [X] Accept port as CLI argument
- [ ] Server should run in 865 compliant or non-compliant mode (512
character truncation)
- [ ] Server should listen on UDP
- [ ] Server should be well tested
- [ ] Installable through homebrew

## Install through homebrew

Notes:
  - https://github.com/Homebrew/homebrew-binary
  - https://github.com/Homebrew/homebrew/pull/20895
  - https://github.com/Homebrew/homebrew/pull/21810/files
  - https://github.com/Homebrew/homebrew/issues/23703

## Author

This was written by [Nick Rowe](http://dcxn.com) in the summer of
2014 with guidance by Rafe Colton.
