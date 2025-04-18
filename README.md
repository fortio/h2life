# h2life

Streams game of life over http as ANSI codes

Demo:

https://github.com/user-attachments/assets/8ad28d6e-169b-462f-9d9b-7ad5654e3d29


## Deployment

Demo/deployed under/try it:
```bash
curl --compressed https://demo.fortio.org/life
```
Compressed save my bandwidth (gzip).

or
```bash
curl https://demo.fortio.org/life
```

## Try locally:

```bash
go install github.com/fortio/h2life@latest
h2life -iter 1000 -delay 10ms &
curl http://localhost:31337/life
```

## Options

Main flags are:
```
h2life 0.1.0 usage:
	h2life [flags]
or 1 of the special arguments
	h2life {help|envhelp|version|buildinfo}
flags:
  -delay duration
    	Delay between frames (default 100ms)
  -iter int
    	Number of iterations per request (in addition to the initial) (default 79)
```
(see `h2life help` for full list)
