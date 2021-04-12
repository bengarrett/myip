# MyIP Tetrad

[![goreleaser](https://github.com/bengarrett/myip/actions/workflows/release.yml/badge.svg)](https://github.com/bengarrett/myip/actions/workflows/release.yml)

MyIP Tetrad is a simple to use terminal tool to determine your Internet-facing IP address and location from four remote sources. Developed on Go, it's a portable self-contained binary with no dependencies.

It is an excellent tool for quickly determining if your machine or network is connected to the Internet or to see if a VPN is activated.

## Downloads

### Packages

##### macOS [Homebrew](https://brew.sh/)
```sh
brew install bengarrett/homebrew-myip/myip
```

##### Windows [Scoop](https://scoop.sh/)
```sh
scoop bucket add bengarrett https://github.com/bengarrett/myip.git
scoop install bengarrett/myip
```

### Intel
- [Windows](https://github.com/bengarrett/myip/releases/latest/download/myip_Windows_Intel.zip)
- [macOS](https://github.com/bengarrett/myip/releases/latest/download/myip_macOS_Intel.tar.gz
)
- [FreeBSD](https://github.com/bengarrett/myip/releases/latest/download/myip_FreeBSD_Intel.tar.gz
)
- [Linux](https://github.com/bengarrett/myip/releases/latest/download/myip_Linux_Intel.tar.gz
)
- - [APK](https://github.com/bengarrett/myip/releases/latest/download/myip.apk
) (Alpine package)<br>`apk add myip.apk`
- - [DEB](https://github.com/bengarrett/myip/releases/latest/download/myip.deb) (Debian package)<br>`dpkg -i myip.deb`
- - [RPM](https://github.com/bengarrett/myip/releases/latest/download/myip.rpm) (Redhat package)<br>`rpm -i myip.rpm`

### arm
- [macOS on Apple M chips](https://github.com/bengarrett/myip/releases/latest/download/myip_macOS_M-series.tar.gz
)
- [Linux arm32](https://github.com/bengarrett/myip/releases/latest/download/myip_Linux_arm32_.tar.gz
) (Raspberry Pi and other single-board computers)
- [Linux arm64](https://github.com/bengarrett/myip/releases/latest/download/myip_Linux_arm64.tar.gz
)

## Usage

```sh
myip -help
# MyIP Usage:
#     myip [options]:
#
#     -h, --help       show this list of options
#     -f, --first      returns the first reported IP address and its location
#     -s, --simple     simple mode only displays the IP address
#     -t, --timeout    https request timeout in milliseconds (default: 5000 [5 seconds])
#     -v, --version    version and information for this program
```

```sh
myip
# (1/4) 93.184.216.34, Norwell, United States
# (2/4) 93.184.216.34, Norwell, United States
# (3/4) 93.184.216.34, Norwell, United States
# (4/4) 93.184.216.34, Norwell, United States
```

```sh
myip -first
# (1/1) 93.184.216.34, Norwell, United States
```

```sh
myip -simple
# 93.184.216.34
# 93.184.216.34
# 93.184.216.34
# 93.184.216.34
```

```sh
myip -simple -first
# 93.184.216.34
```

```sh
myip -timeout=900
# (1/4) 93.184.216.34, Norwell, United States
# ip4.seeip.org: timeout
# (3/4) 93.184.216.34, Norwell, United States
# api.ipify.org: timeout
```

## Build

[Go](https://golang.org/doc/install) supports dozens of architectures and operating systems letting MyIP to [be built for most platforms](https://golang.org/doc/install/source#environment).

```sh
# clone this repo
git clone git@github.com:bengarrett/myip.git

# access the repo
cd myip

# target and build the app for the host system
go build

# target and build for Windows 7+ 32-bit
env GOOS=windows GOARCH=386 go build

# target and build for OpenBSD
env GOOS=openbsd GOARCH=amd64 go build

# target and build for Linux on MIPS CPUs
env GOOS=linux GOARCH=mips64 go build
```

---

#### MyIP uses the following online APIs.

- [ipify API](https://www.ipify.org)
- [MYIP.com](https://www.myip.com)
- [Workshell MyIP](https://www.my-ip.io)
- [SeeIP](https://seeip.org)

I found [Steve Azzopardi's excellent _import "context"_](https://steveazz.xyz/blog/import-context/) post useful for understanding context library in Go.