# Introduction

* `github.com/zhangyoufu/lumina`

  A Go library speaking Hex-Rays IDA lumina protocol.

* `github.com/zhangyoufu/lumina/cmd/lumina-proxy`

  A proxy server that may help if you have any privacy concerns using the Hex-Rays' official lumina server.

# lumina-proxy

## Why do I need a proxy server?

* When you connect to a lumina server, a hello packet sends your IDA license key to the server.
  This discloses your IP (and potentially geolocation) to Hex-Rays, along with your name / company on your IDA license.
* When you push metadata to lumina server, your hostname, idb path, input binary path, input binary MD5 are disclosed to the server.
  This doesn't sound good.

By using a proxy server that understands lumina protocol, you can benefit from and make contributions to Hex-Rays' lumina server without privacy concerns.

## Quick Start Guide

1. get a copy of `ida.key` (a legitimate license is required)
2. generate TLS private key and certificate (or PKI)
3. `docker run -p 443:443 -v $(pwd)/ida.key:/ida.key:ro -v $(pwd)/cert.pem:/cert.pem:ro -v $(pwd)/key.pem:/key.pem:ro youfu/lumina-proxy -listen :443 -tls`
4. configure your clients to use and trust this proxy server (the following instruction is only applicable to IDA < 8.0)
   * modify `LUMINA_HOST` and `LUMINA_PORT` in `cfg/ida.cfg`
   * copy your TLS certificate into IDA installation directory and rename it to `hexrays.crt`
   * restart IDA

# Further Readings

* https://www.synacktiv.com/publications/investigating-ida-lumina-feature.html
* https://abda.nl/posts/introducing-lumen/
* https://lumen.abda.nl/

# Other Implementations

* https://hex-rays.com/products/ida/lumina/ (official)
* https://github.com/synacktiv/lumina_server
* https://github.com/naim94a/lumen

# Warning

It's not wise to use any untrusted 3rd-party lumina server.
Anyone who has access to your IDA license key can talk to Hex-Rays' lumina server with your identity.
