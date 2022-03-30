# Vatz SDK for plugin

## How to bootstrap your new plugin

**For now, `vatz` repository is private so many part of bootstrapping should be done by manually.**

Clone `vatz` repository.

```
~$ git clone github.com/dsrvlabs/vatz
```

Copy template file into user new project.

```
~$ ./bootstrip/main.template <YOUR PROJECT HOME>/main.go
```

Creat new go project by creating mod file.

```
~$ go mod init <YOUR PROJECT NAME>
```

Add this line on `go.mod` file.

```go.mod
replace github.com/dsrvlabs/vatz/sdk => /home/rootwarp/git/vatz/sdk
```

At last, check packages.

```
~$ go mod tidy
```
