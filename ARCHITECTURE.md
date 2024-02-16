# Project Design (TODO)

The overall architecture is based on domain-driven development and ports-adapters.

Project structure:

```sh
internal/
  common/
  infra/
  stock/
  usecase/
```

## Domains

There are X domains in this project:

- common
- stock
- user
- trade
- screener

### Stock

The app is mainly dealing with stock time series data and behaviours around them.
