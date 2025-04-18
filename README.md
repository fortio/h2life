# h2life

Streams game of life over http as ANSI codes

###

Demo/deployed under:
```bash
curl -N https://demo.fortio.org/life
```


### Try locally:

```bash
go install github.com/fortio/h2life@latest
h2life &
curl -N http://localhost:31337/life
```
