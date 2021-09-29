# luguz

Luguz is a captcha implementation in go. Its a rest api server that serves
base64 encoded captcha challenges along with its solutions. The rendered captcha
image is basically a blank background with a bunch of random circles one of them
is open. The solution is the coordinates of a rectangle that borders the open
circle.

Example request:

```
curl http://host:port/captcha 

{"solution":{"X":350,"Y":50,"W":142,"H":142},"data":"BASE_ENCODED_PNG_IMAGE"}
````

## Make luguz

```shell
# to run tests:
make test
# to build luguz binary:
make build
```

## Cli flags

```
  -cache int
        The number of the pre-rendered captcha images to cache in memory. (default no cache)
  -circles int
        The number of the circles in the rendered captcha image. (default 10)
  -height int
        The height of the rendered captcha image. (default 200)
  -width int
        The width of the rendered captcha image. (default 500)
```

## Cached pre-rendering mode

In this mode, luguz will build in-memory cached list of pre-rendered captchas and serve the requests from it.
At the same time, a process is running in background and preparing new cashes, replacing the main cache when it depletes. 
If the incoming requests are faster than the cache regenerating process, lugus will serve a random dirty(used before) captchas until the cache is renewed.

Cached mode is enabled by passing option ``-cache=1000`` for a size 1000 cache.

## Testing hacks

I use this following bash line to test the ouput real quick:

```
curl http://localhost:8080/captcha |  jq -r ".data" | base64 -d | display
```

`dispaly` is part of imagemagick tool and `jq` to parse the json output.
