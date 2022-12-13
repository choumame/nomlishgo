# NomlishGo
Go上で[ノムリッシュ翻訳](https://racing-lagoon.info/)/[ビジネッシュ翻訳(意識高い系ビジネス用語翻訳)](https://bizwd.net)を行えるパッケージです。  
NomlishGo is a package for [Nomlish Translation](https://racing-lagoon.info/) and [Businessh Translation](https://bizwd.net) on Golang.

## How to use
### Instlation
```go 
go get github.com/chouzame/nomlishgo
```

### Usage
Import
```go
import "github.com/chouzame/nomlishgo"
```
Nomlish Translation
```go
nom, err := nomlishgo.ToNomlish("text", 2)
```
Businessh Translation
```go
bus, err := nomlishgo.ToBusinessh("text", 2)
```