# Damplify

<img src="./img/shout.svg" align="right" width="250">

Damplify is a DNS amplification attack tool implemented in Go. 
It was created for educational purposes. 
Please don't misuse it for illegal activities.

Damplify works on Linux only.

## Background

DNS amplification attacks are a form of reflective DDOS. 
The attacker sends DNS queries with spoofed source IP to recursive resolvers. 
The resolvers send their answers to the spoofed source IP.

As the responses from the recursive resolvers can be many times bigger than the query, the traffic 
is amplified. 
Furthermore, blocking the traffic is hindered as it originates from a large number of legitimate systems.

## Installation

Build damplify with the following command (requires go):
```
go install github.com/nick-triller/damplify
```

## Usage

```
A DNS ampflification attack tool

Usage:
  damplify <targetIP> <targetPort> [flags]

Flags:
  -h, --help                   help for damplify
      --resolversPath string   Path to file containing resolver IPs (default "resolvers.txt")
      --workers int            Number of worker routines (default 10)
```

Example invocation:
```
damplify --resolversPath resolvers.txt 192.168.178.40 9998
```

## Limitations

As this is a proof of concept, there are some limitations:
- The DNS query is hardcoded (TXT query for cloudflare.com, amplification ~x10).
- The tool runs on linux only. Golang doesn't implement raw sockets on Windows.
- Resolvers and target IPs must be IPv4.

## Attributions

The logo was created by [Freepik](https://www.flaticon.com/authors/freepik).

The included DNS resolver list [resolvers.txt](./resolvers.txt) was taken
from [public-dns.info](https://public-dns.info/).
