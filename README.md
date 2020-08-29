# Typing practice on command line

![2020-08-29-214923_1918x1059_scrot](https://user-images.githubusercontent.com/24639564/91644991-a29a0c80-ea41-11ea-9051-315036ddb1cb.png)

I love [cslarsen/wpm](https://github.com/cslarsen/wpm), i use it almost everyday! so i wanted to make my own version on golang.

### Install
Build the project
```
$ go build
```

Move it in a directory in your $PATH
```
# mv gopm /usr/local/bin/
```

### Usage
```
$ gopm FILES...
```

Example: start practicing on files `test.txt` and `test2.txt`
```
$ gopm test.txt test2.txt
```

Or run it without building building the project
```
go run main.go -- [...FILES]
```

### Test
```sh
$ go test ./...
```

### Keybindings
| quit | \<Ctrl-c\> |
|------|-------|
