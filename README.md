# AliDNS Forwarder
A 70 line program to expose authenticated alidns json api as dns-over-https endpoint in `/dns-query`. 

## Usage
`alidns_forwarder --uid <User ID> --ak <AccessKey ID> --secret <AccessKey Secret>`

Please refer to [AliYun help document](https://help.aliyun.com/dns/json-api-for-doh) for meanings of each parameters.

You may also supply `--port` option to make `alidns_forwarder` listen on different port. By default, it will listen on `8080`.

## Limit
- `DNS Wireformat` is not supported. PRs welcome.
- This program does NOT provide a way to authenticate received DoH requests. You should put it behind a reverse proxy like `caddy`!
- By default, `alidns_forwarder` will return an extra parameter `uid` in DoH response. [This is a upstream issue](https://help.aliyun.com/dns/json-api-for-doh) and we are not planned to remove this param in `alidns_forwarder`. PRs welcome.