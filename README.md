# setop

Perform set operations on files.

## Install

```
go install github.com/benjaminheng/setop@latest
```

## Usage

Available commands:

```
setop intersection <file1> <file2>
setop difference <file1> <file2>
```

Examples:

```bash
$ seq 1 10 > /tmp/a

$ seq 5 15 > /tmp/b

$ setop intersection /tmp/a /tmp/b  # lines in both A and B
5
6
7
8
9
10

$ setop difference /tmp/a /tmp/b  # lines in A but not B
1
2
3
4
```
