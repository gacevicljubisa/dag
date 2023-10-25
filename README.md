# dag

This is example of DAG (Directed Acyclic Graph) in Go.

## Usage

```
$ go run .
```

## Info

Error handling is omitted for simplicity.  
In production, instead of `panic`, error should be returned to caller.
Id should be unique, but difference in the case of `string` is ignored for simplicity.