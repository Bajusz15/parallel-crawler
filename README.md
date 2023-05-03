# parallel-crawler

## Description
This tool makes a request to the given url(s) and prints the md5 hash of the
response body along with the formatted url address.

The tool uses fan-out and fan-in pattern to send out the requests then "collect" the responses.

## Usage of the tool
`./myhttp -parallel [number of parallel requests (number)] [address1] [address2] [...]`
***
Example:
```
./myhttp -parallel 3 google.com facebook.com yahoo.com yandex.com twitter.com
```
Expected response:
```
2023/05/03 19:55:39 https://facebook.com d7dd8b081fc735d1d19d8170f812f70a
2023/05/03 19:55:39 https://yandex.com 8659c2f05bfc7b5dbc9b55a2b28edcc3
2023/05/03 19:55:39 https://google.com 4a52893a0e5c9bfce350fa1d98baea3f
2023/05/03 19:55:39 https://twitter.com 13d3ee4892a3a7ed0fa861c83b5c7820
2023/05/03 19:55:39 https://yahoo.com fd6b12d1f5cab48ec2052381919d3c5d
```

The parallel flag is optional.
