# SIS Project — Checkpoint 1: Statistics Calculator

Individual project for *Golang Application Development* (KBTU, Fall 2026).

A small console application written in Go. A **task** is a list of numbers.
The application validates the input, processes the task, and stores the
resulting statistics (count, min, max, sum, mean, median) in the same task.
All data is stored in memory.

Example: input `[10, 20, 30]` → `count=3 min=10 max=30 sum=60 mean=20 median=20`.

## Requirements

Go 1.27

## Run

```bash
go run .
```

## Checks and tests

```bash
gofmt -l .
go vet ./...
go test ./...
```

## Project structure

```
main.go             demo of the application
task/task.go        Task, State, Result, errors, validation, state changes
task/processor.go   Processor interface and StatsProcessor (statistics calculation)
task/service.go     Service: creating, finding and processing tasks (in-memory map)
task/task_test.go   unit tests
```

## Task lifecycle

```
Pending -> Processing -> Completed
                     \-> Failed
```

| From         | Allowed next states   |
|--------------|-----------------------|
| `Pending`    | `Processing`          |
| `Processing` | `Completed`, `Failed` |
| `Completed`  | none (final state)    |
| `Failed`     | none (final state)    |

- `Pending` — the task is created but not processed yet.
- `Processing` — the task is being processed.
- `Completed` — processing finished successfully and the result is saved.
- `Failed` — an error occurred during processing (for example, the sum of the numbers is too large).

Any other transition returns `ErrInvalidTransition`.
Processing a `Completed` or `Failed` task again is forbidden, because the
result has already been produced.

## Input validation

Input is invalid (`ErrInvalidInput`) if the list is empty or contains `NaN`
or `±Inf`. Statistics of such values are meaningless.

## Errors

| Error                  | When it occurs                              |
|------------------------|---------------------------------------------|
| `ErrInvalidInput`      | empty list or invalid number                |
| `ErrTaskNotFound`      | no task with the given ID                   |
| `ErrInvalidTransition` | forbidden state transition                  |
| `ErrProcessingFailed`  | an error occurred during processing         |

Expected errors are returned as `error` values. `panic` is not used.

## Design decisions

- **`Processor` interface:** `Service` does not depend on a concrete
  processing algorithm.
- **`map[int]*Task`:** fast lookup by ID; a pointer allows changing the task
  directly in the storage.
- **Pointer receivers** (`*Task`, `*Service`) are used where the state
  changes; **value receivers** (`StatsProcessor`) where there is no data to change.
- **Input slice is copied** when a task is created, so external code cannot
  change the task's data.
- **IDs** are sequential integers starting from 1.
