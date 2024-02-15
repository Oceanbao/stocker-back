# Project Design

The overall architecture is based on domain-driven development and ports and adapters.

Project structure:

```sh
domains/
  stocks/
    models.go
    repository.go
    usecase.go
  users/
    models.go
    repository.go
    usecase.go
```

## Domains

There are 2 domains in this project:

- stocks
- users

### Stocks

The app is mainly dealing with stock time series data and behaviours around them.

#### Models

```go
// Ticker (entity) is a unique stock object.
type Ticker struct {
  Name string
  Cap float64
}
```
