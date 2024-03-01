# bimiprobe

Reads in a list of domains from stdin and outputs any associated BIMI record

## Install

If you have Go installed and configured (i.e. with `$GOPATH/bin` in your `$PATH`):

```
go install github.com/cybercdh/bimiprobe@latest
```

## Usage

`$ echo sub.example.com | bimiprobe`

or 

`$ cat <file> | bimiprobe`

### Options

```
Usage of bimiprobe:
  -c int
      set the concurrency level (default 20)
  -dns string
      Custom DNS resolver address (ip:port)
  -port string
      DNS server port (default "53")
```