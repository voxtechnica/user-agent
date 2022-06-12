# UserAgent

UserAgent parses an HTTP [User-Agent](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent) string to
determine basic device, operating system, and client application characteristics. It's designed to be fast and
reasonably accurate at identifying browsers, bots, and mobile applications with a large market share. Smaller clients
will either be classified as "Other", or as the browser they're derived from (e.g. Waterfox will be classified as
Firefox).

UserAgent is intended to help you answer questions like the following:

* What percentage of our site traffic is coming from mobile phones, tablets, and desktop/laptop computers?
* Which browsers should we be using to test our web application? Can we stop supporting Internet Explorer?
* Which browser versions are people using today? Can we use this particular CSS feature?
* If we were going to develop native application(s), which operating system(s) should we support?

If you care more about the specific make/model of the mobile device someone is using, you would be better served
using something like the Matomo Analytics [Device Detector](https://github.com/matomo-org/device-detector). It's
an actively-maintained collection of YAML files with regular expressions focused on accurate device detection.
This accuracy comes at a performance cost, though. A cascade of regular expression matches is going to be slow.

If you'd like a decent sample of real-world User-Agent strings to test with, the
file [user_agents.json](cmd/sample_data/user_agents.json) contains a collection of 17,548 unique User-Agent strings
from 458,143 webpage views in Q2 of 2022. They're grouped, having been characterized using a Java port of
the [Device Detector](https://github.com/mngsk/device-detector). Note that they may have been mis-characterized in
those groupings.

## Installation

```bash
go get github.com/voxtechnica/user-agent
```

## Usage

Use `user_agent.Parse(header)` to create a `UserAgent`, as illustrated below.

```go
header := "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"
ua := user_agent.Parse(header)
fmt.Println(ua.String()) // Browser Chrome 101.0 Desktop Windows 10.0
```

## Performance

The User-Agent parser is pretty fast. It's based on `strings.Contains` instead of using regular expressions.
It takes 3-4 microseconds to parse a User-Agent header, as indicated in the example benchmark results below.

```text
go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/voxtechnica/user-agent
cpu: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz
BenchmarkParse/Parse-Googlebot-12   420231   2684 ns/op    835 B/op   12 allocs/op
BenchmarkParse/Parse-Chrome-12      332326   3270 ns/op   1589 B/op   15 allocs/op
BenchmarkParse/Parse-Firefox-12     373551   2997 ns/op   1200 B/op   16 allocs/op
BenchmarkParse/Parse-Safari-12      275498   4242 ns/op   2032 B/op   19 allocs/op
PASS
ok    github.com/voxtechnica/user-agent    4.654s
```
