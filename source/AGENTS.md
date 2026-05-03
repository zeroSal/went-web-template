## Project Structure

- `main.go` — Entry point, creates root `cobra.Command`, uses embed.FS for templates/static
- `app/cmd/` — Command implementations using custom `went-command` pattern
- `app/ctrl/` — Controller implementations using custom `went-web` pattern (Iris handlers)
- `app/kernel.go` — DI kernel that creates and runs the fx.App
- `app/container.go` — FX module container with DI providers (dynamically assembled from services + auto-discovered controllers)
- `app/specs.go` — Build metadata (version, channel, build date)
- `app/bootstrap/` — Bootstrap providers:
  - `init.go` — Initialization functions: validates env, creates working dirs, initializes loggers, mounts controllers
- `app/service/env/env.go` — Env struct with `Load()` and `Validate()` methods
- `app/service/logger/` — Logger providers:
  - `audit.go` — AuditLogger provider (wraps went-logger FileLogger)
  - `error.go` — ErrorLogger provider (wraps went-logger FileLogger)
- `registry/` — Singleton registries for self-registration pattern:
  - `command.go` — Command registry (`registry.Command`)
  - `controller.go` — Controller registry (`registry.Controller`)
- `res/` — Embedded application resources (banner.template)
- `.env.dist` — Environment template file (must be copied to .env and customized)

---

## Environment Configuration

Copy `.env.dist` to `.env` and configure:

| Variable  | Default     | Description                          |
|-----------|-------------|--------------------------------------|
| `ENV`     | `dev`       | Environment (`dev` or `prod`)        |
| `VAR_DIR` | `var`       | Base directory for runtime files     |
| `HOST`    | `127.0.0.1` | The address where to bind the server |
| `PORT`    | 3096        | The port where to bind the server    |


The `Validate()` method enforces `ENV` as `dev` or `prod` and validates the port validity.

---

## Tech Stack

- **Cobra** — Root command framework (integrated via `went-command`)
- **Uber FX** — Dependency injection container
- **Iris v12** — HTTP web framework
- **went-clio** — CLI output formatting (console input/output, banners)
- **went-logger** — File logging (AuditLogger, ErrorLogger)
- **went-command** — Custom command pattern with registry
- **went-web** — Iris integration with controller pattern and registry
- **godotenv** — `.env` file loading
- **golangci-lint v2** — Linter aggregator

---

## Patterns

### Registry Pattern

Two global singleton registries in `registry/` package. Commands and controllers self-register via `init()`:

```go
// registry/command.go
var Command = command.NewRegistry()

// registry/controller.go
var Controller = controller.NewRegistry()
```

Adding a new command or controller only requires creating a new file with an `init()` block — no central wiring needed.

### Command Pattern

Commands implement `command.Interface` and use `GetHeader()` to define cobra metadata. They self-register via `init()`:

```go
var _ command.Interface = (*Serve)(nil)
type Serve struct {
    command.Base
}

func init() {
    registry.Command.Register(&Serve{})
}

func NewServe() command.Interface { ... }

func (c *Serve) GetHeader() command.Header {
    return command.Header{
        Use:   "serve",
        Short: "Start the web server",
        Long:  "...",
    }
}

func (c *Serve) Invoke() any {
    return c.run
}

func (c *Serve) run(ctx context.Context, env *env.Env, clio *clio.Clio, app *iris.Application) error { ... }
```

### Controller Pattern

Controllers implement `controller.Interface` with `Register(*iris.Application)` to define routes. They self-register via `init()` and are collected by FX via `group:"controllers"`:

```go
var _ controller.Interface = (*Home)(nil)
type Home struct {
    controller.Base
    env *env.Env
}

func init() {
    registry.Controller.Register(NewHome)
}

func NewHome(env *env.Env) *Home {
    return &Home{env: env}
}

func (c *Home) Register(app *iris.Application) {
    app.Get("/", c.index())
}

func (c *Home) index() iris.Handler { ... }
```

### Kernel Pattern

The `Kernel` wraps fx.App creation and execution:

```go
kernel := app.NewKernel(ctx, EmbedFS, specs, clio)
kernel.Run(command.Invoke())
```

`Run()` creates an FX app with the `Container`, supplies `Clio`/`EmbedFS`/`Specs`, provides context, invokes the given function, then starts and stops the app.

### FX Container

The container in `container.go` is assembled dynamically in `init()`:

```go
var services = []fx.Option{
    fx.Provide(logger.NewAuditLogger),
    fx.Provide(logger.NewErrorLogger),
    fx.Provide(env.Load),
    wentweb.Bundle,
    bootstrap.Init,
}

func init() {
    opts := services
    for _, constructor := range registry.Controller.All() {
        opts = append(opts, fx.Provide(
            fx.Annotate(
                constructor,
                fx.As(new(controller.Interface)),
                fx.ResultTags(`group:"controllers"`),
            ),
        ))
    }
    Container = fx.Options(opts...)
}
```

Controllers are dynamically collected from `registry.Controller` and annotated with `fx.As(new(controller.Interface))` + tagged into `"controllers"` group.

### Bootstrap

`bootstrap.Init` is an `fx.Options` bundle that runs at startup:

```go
var Init = fx.Options(
    fx.Invoke(initWorkingDirs),
    fx.Invoke(validateEnv),
    fx.Invoke(initLoggers),
    fx.Invoke(controller.Mount),
)
```

1. `initWorkingDirs` — Creates `var/logs/` directory
2. `validateEnv` — Calls `env.Validate()` to enforce dev/prod
3. `initLoggers` — Calls `Init()` on both AuditLogger and ErrorLogger
4. `controller.Mount` — From `went-web`, mounts all grouped controllers to the Iris app

### Env Config

`app/service/env/env.go` defines a single typed struct with `Load()` (reads from environment with defaults) and `Validate()` methods. It uses `godotenv` to load `.env` file. No raw `os.Getenv` calls outside this file.

### Runtime Directories

The bootstrap initialization creates the following directories at startup:
- `var/logs/` — Log files (audit.log, error.log)

---

## Layer Conventions

Both `app/service/` and `app/model/` follow the same rule: **one package per domain**. The package name defines the domain context; each filename describes the type of object it contains, never repeating the domain name.

### `app/service/`

```
app/service/
├── domain_1/
│   └── client.go            # package domain_1
└── domain_2/
    ├── client.go            # package domain_2
    ├── client_interface.go  # package domain_2
    └── loader.go            # package domain_2
```

Typical filenames: `client.go`, `client_interface.go`, `loader.go`, `resolver.go`, `cache.go`, `parser.go`

### `app/model/`

```
app/model/
├── domain_1/
│   └── entity.go            # package domain_1
└── domain_2/
    ├── entity.go            # package domain_2
    ├── dto.go               # package domain_2
    ├── request.go           # package domain_2
    └── response.go          # package domain_2
```

Typical filenames: `entity.go`, `dto.go`, `request.go`, `response.go`, `enum.go`, `event.go`

> **Why:** the domain lives in the package, the role lives in the filename.
> `github/github_client.go` ❌ → `github/client.go` ✅
> Imports stay self-documenting: `domain2.Client`, `domain2.Loader`, `release.Entity`.

---

## Code Conventions

### One object per file

Each file defines exactly one primary object (struct or interface). The filename must match the role of that object, as described in the layer conventions above.

### Constructor

Every object must have a constructor named `New<ObjectName>`. Its parameters are exclusively the dependencies to be injected. **No logic, initialization, or side effects of any kind inside the constructor** — it only assigns fields.

```go
func NewKernel(context context.Context, embedFS embed.FS, specs *Specs, clio *clio.Clio) *Kernel {
    return &Kernel{
        context,
        embedFS,
        specs,
        clio,
    }
}
```

### Function naming

Function and method names must always be a **verb** or a **verbNoun** (`load`, `fetchReleases`, `parseResponse`). Names must be self-explanatory in isolation.

Context already provided by the receiver type must not be repeated in the method name:

```go
// receiver is ElementManager
func (m *ElementManager) getAll() { ... }      // ✅
func (m *ElementManager) getElements() { ... } // ❌ redundant

// receiver is ReleaseClient
func (c *ReleaseClient) fetch() { ... }        // ✅
func (c *ReleaseClient) fetchRelease() { ... } // ❌ redundant
```

---

## Lint

- Config file: `.golangci.yml` in v2 format
- `staticcheck` enabled with exclusions: `-ST1000` (package comments)

---

## Build

- All builds use vendor mode: `GOFLAGS="-mod=vendor"`
- Production builds require `VERSION` and `CHANNEL` env vars to be set
- Build artifacts are output to `build/` with naming: `clitemplate-{VERSION}-{CHANNEL}-{OS}-{ARCH}`
- Use `make build-dev`, `make build-staging`, or `make build` (production)

### Makefile Targets

| Target | Description |
|---|---|
| `test` | Run `go test ./...` with coverage to `coverage/` |
| `go-mod-tidy` | Run `go mod tidy` |
| `go-mod-vendor` | Runs `go-mod-tidy` then `go mod vendor` |
| `build-linux-amd64` | Cross-compile Linux amd64 (requires `VERSION`, `CHANNEL`) |
| `build-linux-arm64` | Cross-compile Linux arm64 (requires `VERSION`, `CHANNEL`) |
| `build-darwin-amd64` | Cross-compile macOS amd64 (requires `VERSION`, `CHANNEL`) |
| `build-darwin-arm64` | Cross-compile macOS arm64 (requires `VERSION`, `CHANNEL`) |
| `build` | Runs `go-mod-vendor`, requires `VERSION` + `CHANNEL`, builds all 4 platforms |
| `build-staging` | Builds all 4 platforms with `VERSION=0.0.1 CHANNEL=staging` |
| `build-dev` | Builds single `build/app-dev` binary with `VERSION=0.0.1 CHANNEL=dev` |

---

## Bootstrap Flow

1. `main()` creates `clio`, reads banner template, builds `Specs`, creates `Kernel`
2. `main()` creates root `cobra.Command`
3. `command.Mount()` iterates `registry.Command.All()` — each command's `GetHeader()` defines cobra subcommand, `Invoke()` returns the run function
4. When a subcommand executes, `kernel.Run(command.Invoke())` is called
5. `kernel.Run()` creates an FX app with `Container` + supplied values + the invoke function
6. FX resolves all dependencies: `env.Load()`, loggers, `wentweb.Bundle`, controllers
7. `bootstrap.Init` runs: validates env, creates working dirs, initializes loggers, mounts controllers to Iris
8. The command's `run()` function executes (e.g., `app.Listen(":8080")`)
