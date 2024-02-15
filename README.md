# Stocker Back

Back-end server for stocker app built upon `pocketbase`.

## Project Architecture

Based on domain-driven development.

```
domains/
  stocks/
    models/
      entities/
      valueobjects/
      aggregates/
    repositories/
    usecases/
    infra/
      pocketbase/
  users/
    models/
    repositories/
    usecases/
    infra/
```
