# Go packages

Common packages to be used in projects using Golang.

## Installation


### Configure `go get` to download over SSH

```bash
git config --global url."git@gitlab.vpc-xpto-01:".insteadOf "https://gitlab.vpc-xpto-01/"
```

### Export the environment variable GOPRIVATE

```bash
export GOPRIVATE='gitlab.vpc-xpto-01'
```

### Check golang version

This project requires [Go Modules](https://blog.golang.org/using-go-modules) which is available since version 1.11. Check your version using the following command:

```bash
go version
```

## Usage

```go
package main

import (
	"fmt"

	"gitlab.vpc-xpto-01/oss/gopkg.git/math"
)

func main() {
	v := math.Sum(1,2)
	fmt.Println(v)
}
```