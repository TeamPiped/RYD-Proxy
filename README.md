# RYD-Proxy

An open-source, non-logging proxy for the "Return YouTube Dislike" API server.

Only fetching the Dislike count is supported, and will ever be supported.

Requests are spoofed to look like how they would appear when using the Tor Browser with the official extension.

IPs are rotated every 2 minutes to avoid rate limiting.

## Example Request

https://ryd-proxy.kavin.rocks/votes/dQw4w9WgXcQ

```js
{"id":"dQw4w9WgXcQ","dateCreated":"2022-04-09T22:01:38.222268Z","likes":14589269,"dislikes":390375,"rating":4.8957585373858015,"viewCount":1232906190,"deleted":false}
```

## Why

https://github.com/Anarios/return-youtube-dislike/issues/344
