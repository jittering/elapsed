# elapsed

A simple utility which adds elapsed time markers on every input line.

```shell
$ while true; do echo "hello $(date)"; sleep $(ruby -e 'puts rand'); done | elapsed
[0ms     +0ms    ] hello Mon Jan  3 16:09:11 EST 2022
[274ms   +274ms  ] hello Mon Jan  3 16:09:12 EST 2022
[1343ms  +1068ms ] hello Mon Jan  3 16:09:13 EST 2022
[1643ms  +299ms  ] hello Mon Jan  3 16:09:13 EST 2022
[2469ms  +826ms  ] hello Mon Jan  3 16:09:14 EST 2022
[2600ms  +130ms  ] hello Mon Jan  3 16:09:14 EST 2022
[3361ms  +760ms  ] hello Mon Jan  3 16:09:15 EST 2022
[4316ms  +955ms  ] hello Mon Jan  3 16:09:16 EST 2022
```

The first number is the absolute elapsed time in milliseconds and the second is
the delta between lines of input.

## Installation

Simply download and install it with one of the following methods:

via homebrew (mac or linux):

```sh
brew install jittering/kegs/elapsed
```

or manually:

Download a [pre-built binary](https://github.com/jittering/elapsed/releases) or
build it from source:

```sh
# requires go
git clone https://github.com/jittering/elasped.git
cd elapsed
# binary will be written to build/elapsed
make build
```

## Usage

By default, with no command-line flags, `elapsed` displays the absolute elapsed
time as well as the delta between lines. The following flags are available which
alter this behavior:

```text
  -datetime
        Show date/time stamp when message was received
  -format string
        Date/time format (default: 2006-01-02T15:04:05Z07:00)
  -no-delta
        Do not print the delta elapsed time
  -no-elapsed
        Do not print the absolute elapsed time
  -slow
        Show only slow deltas (over a certain threshold)
  -slow-ms int
        Slow delta threshold in ms (default 500)
```

The `-format` flag modifies the date/time string format used in conjunction with
`-datetime`. The format strings should be as defined by the go [time/format](https://github.com/golang/go/blob/master/src/time/format.go#L92-L108) package. Some common formats:

```text
Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
ANSIC       = "Mon Jan _2 15:04:05 2006"
UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
RFC822      = "02 Jan 06 15:04 MST"
RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
RFC3339     = "2006-01-02T15:04:05Z07:00"
RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
Kitchen     = "3:04PM"

// Handy time stamps.
Stamp      = "Jan _2 15:04:05"
StampMilli = "Jan _2 15:04:05.000"
StampMicro = "Jan _2 15:04:05.000000"
StampNano  = "Jan _2 15:04:05.000000000"
```

## Configuration

If you always use the same set of flags, add them to `~/.elapsedrc`. This file
will be read on startup *if no other CLI flags are given*.

## License

[MIT](./LICENSE), (c) 2022, Pixelcop Research, Inc.
