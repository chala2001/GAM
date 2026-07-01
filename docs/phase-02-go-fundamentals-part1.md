# Phase 2 — Go Crash Course, Part 1

Runnable code for this phase lives in [`docs/phase-02/`](phase-02/), one tiny `main` program per concept under `cmd/`. Nothing here is part of the real services yet — it's throwaway code whose only job is to make the next section's concepts click.

## New Go concepts in this phase

Toolchain (`go mod`, `go run`, `go build`), packages, variables & types, zero values, `struct`s vs Java classes, methods with receivers, pointers, error handling (`error` as a value, not exceptions), slices & maps.

---

## 1. Toolchain

Go ships one binary (`go`) that replaces Maven/Gradle + `javac` + `java`:

| Command | Java equivalent | What it does |
|---|---|---|
| `go mod init <module>` | `pom.xml` / `build.gradle` creation | Declares a module (like a Maven `groupId:artifactId`) and creates `go.mod`. |
| `go run ./path` | `mvn compile exec:java` | Compiles and runs in one step — no separate `.class`/jar artifact left behind. |
| `go build ./...` | `mvn package` | Compiles everything into binaries, without running them. |
| `go vet ./...` | a linter pass | Static analysis for common mistakes (unused results, format-string errors, etc). |
| `gofmt` | a code formatter (like `spotless`) | Canonical formatting — the Go community doesn't argue about style; `gofmt` decides. |

There's no build file listing dependency versions and plugins the way `pom.xml` does — `go.mod` is much smaller because most of what a build tool does in Java (compilation, packaging) is just built into the `go` command.

Run it:
```
cd docs/phase-02
go run ./cmd/01-hello-toolchain
```

## 2. Packages

Every Go file starts with `package <name>`. A directory is a package; a package named `main` with a `func main()` is an executable entry point — the equivalent of a Java class with `public static void main(String[] args)`, except the class-vs-package distinction doesn't exist in Go. There's no `public class HelloWorld { }` wrapper — top-level functions are normal, not a special case.

`import "fmt"` pulls in the standard library's formatted-I/O package (`fmt.Println`, `fmt.Sprintf`, etc.) — the rough equivalent of `System.out.println` plus `String.format`.

## 3. Variables, types, and zero values

```go
var name string
var requestCount int
var latencyMs float64
var isHealthy bool
fmt.Println("zero values:", name, requestCount, latencyMs, isHealthy)
```
Output: `zero values:  0 0 false`

This is the biggest early gotcha coming from Java: **Go has no `null` for these types.** A `var` declaration without an initializer doesn't leave the variable pointing at nothing — it gets the type's *zero value*: `""` for strings, `0` for numbers, `false` for bools, `nil` only for reference-like types (pointers, slices, maps, interfaces, channels, functions). This eliminates an entire category of `NullPointerException`-shaped bugs by construction, at the cost of needing an explicit sentinel (like `-1` or a second `ok bool`) when you actually need to represent "not set."

Two ways to declare with a value:
```go
var apiName string = "orders-api"  // explicit type
version := 1                        // := infers the type, only valid inside a function
```
`:=` is Go's most-used declaration form and only works for new local variables — there's no top-level equivalent, mirroring how Java requires explicit types for fields but lets `var` infer locals since Java 10.

Code: [`cmd/02-variables-types/main.go`](phase-02/cmd/02-variables-types/main.go)

## 4. `struct`s vs Java classes

```go
type RegisteredAPI struct {
	Name     string
	Upstream string
	Path     string
}
```
A `struct` is a plain data record — think a Java class with only fields and no methods, and (critically) **no inheritance**. Go has no `extends`. Composition (embedding one struct inside another) is the tool Go gives you instead, and we'll hit that in a later phase once it's motivated by real code (e.g. embedding a common `BaseHandler`-style config across services).

Field capitalization is significant: `Name` (capital) is exported (visible outside the package, like Java's `public`); a lowercase `name` field would only be visible inside the same package (closer to Java's package-private than `private`, since there's no equivalent of a strictly single-class-only `private` at the struct-field level).

## 5. Methods with receivers

```go
func (a RegisteredAPI) Describe() string {
	return fmt.Sprintf("%s -> %s%s", a.Name, a.Upstream, a.Path)
}
```
There's no `class` keyword to hang a method inside. Instead, a function gets a *receiver* — `(a RegisteredAPI)` — which is Go's way of saying "this function is a method on `RegisteredAPI`, and inside the body, `a` is like Java's implicit `this`." You call it exactly like a Java method: `orders.Describe()`.

Zero-value structs print all fields at their zero value — `{Name: Upstream: Path:}` — reinforcing point 3: an unconstructed `RegisteredAPI{}` is a valid, non-null value, not a crash waiting to happen.

Code: [`cmd/03-structs-methods/main.go`](phase-02/cmd/03-structs-methods/main.go)

## 6. Pointers

```go
func (q Quota) TryConsumeByValue() { q.Remaining-- }
func (q *Quota) TryConsumeByPointer() { q.Remaining-- }
```
Output:
```
after value-receiver call: 100
after pointer-receiver call: 99
```
Every Go value passed to a function — including a method's receiver — is **copied** by default, the same as Java passing a primitive `int` by value. `TryConsumeByValue` mutates its own private copy of `q`, so the caller's `Quota` is untouched. `TryConsumeByPointer` takes `*Quota` — a pointer, Go's explicit version of "pass by reference" — so the mutation is visible to the caller.

This is different from Java, where every object variable is *already* a reference, so this choice never comes up explicitly. In Go you decide it per method: read-only/small structs usually take value receivers; anything that mutates state, or is large enough that copying is wasteful, takes a pointer receiver. This exact decision matters a lot once we're writing service structs later in the roadmap (e.g. a `RateLimiter` whose internal counters must be shared, not copied, across calls).

`&q` takes the address of `q` (produces a `*Quota`); `*qPtr` dereferences it back to the `Quota` value. There's no pointer arithmetic like C — Go pointers only ever point at one value of their declared type.

Code: [`cmd/04-pointers/main.go`](phase-02/cmd/04-pointers/main.go)

## 7. Error handling as a value

```go
func validateAPIName(name string) error {
	if name == "" {
		return errors.New("api name cannot be empty")
	}
	return nil
}
```
Go has no `throw`/`try`/`catch`. `error` is just an interface type, and by convention it's the last return value of any function that can fail. The caller is *forced* to look at it because Go doesn't silently propagate exceptions up a call stack:
```go
if err := validateAPIName(name); err != nil {
	fmt.Println("rejected:", err)
	continue
}
```
The upside: control flow is always visible at the call site — you can see exactly which calls can fail and what happens when they do, without needing to trace which checked/unchecked exceptions a Java method might throw. The downside, which you'll feel almost immediately: `if err != nil { return err }` shows up after nearly every call, which reads as boilerplate-heavy until it becomes muscle memory. This pattern is non-negotiable across the whole project — every gRPC call, every SQL query, every Redis call in later phases returns an `error` this same way.

Code: [`cmd/05-errors/main.go`](phase-02/cmd/05-errors/main.go)

## 8. Slices & maps

```go
var routes []string
routes = append(routes, "/orders/*")
```
A slice is Go's dynamic array — closest to `ArrayList<String>`, except `append` doesn't mutate in place the way `ArrayList.add` does; it *returns* a (possibly new, possibly resized) slice, which is why you always reassign: `routes = append(routes, ...)`.

```go
rateLimits := map[string]int{"orders-api": 100, "payments-api": 50}
limit, ok := rateLimits["orders-api"]
```
A map is Go's `HashMap<String, Integer>` equivalent. The two-value form `limit, ok := rateLimits[key]` is idiomatic Go: `ok` tells you whether the key existed, so a missing key gives you `(zero value, false)` instead of Java's `null` (or an `Optional` if the map API returns one) — same "no null" philosophy as point 3, applied to lookups.

Code: [`cmd/06-slices-maps/main.go`](phase-02/cmd/06-slices-maps/main.go)

---

## Why this way, not another way

- **Why cover zero values so early, twice (vars and structs)?** It's the single biggest mental model shift from Java and it silently affects every struct we'll write for the next 16 phases (Postgres row structs, gRPC request/response structs, config structs). Better to internalize it on throwaway code than debug it inside the real Auth service.
- **Why introduce pointers via a mutation bug rather than the syntax first?** Seeing `TryConsumeByValue` silently fail to mutate is the kind of thing you hit once, understand permanently, and never forget — closer to how it'll actually bite you later (e.g. forgetting a pointer receiver on a service struct's method and wondering why a counter never increments).
- **Why not touch interfaces, goroutines, or `context.Context` yet?** They all depend on a comfortable grasp of the above first (interfaces need methods+structs to be second nature; goroutines/`context` need to not be fighting the type system while also learning concurrency). Those are the whole subject of Phase 3.

## How to run/verify this phase

```
cd docs/phase-02
go run ./cmd/01-hello-toolchain
go run ./cmd/02-variables-types
go run ./cmd/03-structs-methods
go run ./cmd/04-pointers
go run ./cmd/05-errors
go run ./cmd/06-slices-maps

gofmt -l .   # should print nothing — everything already formatted
go vet ./... # should print nothing — no static-analysis warnings
```

All six programs and both checks were run clean before this doc was written.
