This is a QOTD server written as a small project for exploring Go.

> You can get *anywhere* in ten minutes if you drive fast enough.

## Try it out

- Install the server: `go get -u github.com/nixterrimus/qotd-server`
- Start the server: `qotd-server wisdom.txt`
- In another terminal make a request: `nc localhost 3333`

## Quote File Format

This QOTD server has a mandatory files argument of a quote file, for
example to run with the `wisdom.txt` file included in the git repository
you would use: `qotd-server wisdom.txt`.

The quote file format is exactly the same as the [fortune file
format](http://en.wikipedia.org/wiki/Fortune_(Unix)#Fortune_files), that
is to say:

- It is a text file
- Each quote is separated by a `%` character of its own line

## The QOTD Protocol / RFC 865

RFC 865 defines the [Quote of the Day
Protocol](http://tools.ietf.org/html/rfc865). The specification is
really short.  Here's the gist of it:

- The Server listens on TCP (port 17, by convention)
- The Server may also listen on UDP (also port 17, by convention)
- On connection a quote is served
- The connection is closed immediately after a quote is serve
- Quotes should be less than 512 characters long

## TCP and UDP

By default the server listens on both TCP and UDP on the same port.  If
that doesn't fit your needs, you have a couple of flags:

- `no-tcp` - Starts the server and only listens on UDP
- `no-udp` - Starts the server and only listens on TCP

You can test the TCP interface with netcat: `nc localhost 3333` And you can test on UDP
as well: `echo -n " " | nc -4u -w1 localhost 3333`.

## Next Steps / Project Goals

- [X] Accept file as CLI argument
- [ ] Accept Quotes via Standard In
- [ ] Accept a URL for a quote file
- [X] Accept port as CLI argument
- [ ] Server should run in 865 compliant or non-compliant mode (512
character truncation)
- [X] Server should listen on UDP
- [ ] Server should be well tested
- [ ] Installable through homebrew
- [ ] Advertise QOTD over MDNS
- [ ] Never serve clients the same quote twice

## Client Goals

- [ ] Find Server over MDNS
- [ ] Allow `host port` as args

## Install through homebrew

Notes:
  - https://github.com/Homebrew/homebrew-binary
  - https://github.com/Homebrew/homebrew/pull/20895
  - https://github.com/Homebrew/homebrew/pull/21810/files
  - https://github.com/Homebrew/homebrew/issues/23703

## Author & Thanks

This was written by [Nick Rowe](http://dcxn.com) in the summer of
2014.

Thanks to Rafe Colton for encouraging me to get started with go and
providing guidance on getting started.

Thanks to Parker Moore for hacking on UDP support.
