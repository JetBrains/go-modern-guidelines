# Modern Go Guidelines Explained
**Work in progress** — inconsistencies may be present.

This file provides a more detailed description of the features supported in the Go Modern Guidelines.

**Modernizer legend:**
- [x] — available in go fix as modernizer
- [ ] — not implemented

**Impact legend:**
- Critical — found in almost every project, dozens of occurrences
- High — found often, 5-20 occurrences per project
- Medium — found regularly, 1-5 occurrences per project
- Low — found rarely or in specific code

| Category | Diagnostic | Modernizer | Go | Impact |
|----------|------------|------------|-----|--------|
| Collections | [`slicescontains`](#slicescontains) | [x]        | 1.21 | Critical |
| Collections | [`sortslice`](#sortslice) | [x]        | 1.21 | High |
| Collections | [`minmax`](#minmax) | [x]        | 1.21 | High |
| Collections | [`mapkeysvalues`](#mapkeysvalues) | [ ]        | 1.23 | High |
| Collections | [`mapsloop`](#mapsloop) | [x]        | 1.21 | Medium |
| Collections | [`mapsclone`](#mapsclone) | [ ]        | 1.21 | Medium |
| Collections | [`rangeoverindex`](#rangeoverindex) | [ ]        | 1.21 | Medium |
| Collections | [`slicesmaxmin`](#slicesmaxmin) | [ ]        | 1.21 | Medium |
| Collections | [`slicesreverse`](#slicesreverse) | [ ]        | 1.21 | Medium |
| Collections | [`clearcollection`](#clearcollection) | [ ]        | 1.21 | Medium |
| Collections | [`mapdeletefunc`](#mapdeletefunc) | [ ]        | 1.21 | Low |
| Collections | [`slicescompact`](#slicescompact) | [ ]        | 1.21 | Low |
| Strings | [`stringscutprefix`](#stringscutprefix) | [x]        | 1.20 | High |
| Strings | [`stringsseq`](#stringsseq) | [x]        | 1.24 | High |
| Strings | [`stringsclone`](#stringsclone) | [ ]        | 1.20 | Medium |
| Strings | [`bytescut`](#bytescut) | [ ]        | 1.18 | Medium |
| Loops | [`rangeint`](#rangeint) | [x]        | 1.22 | Critical |
| Loops | [`forvar`](#forvar) | [x]        | 1.22 | High |
| Types | [`efaceany`](#efaceany) | [x]        | 1.18 | Critical |
| Types | [`reflecttypefor`](#reflecttypefor) | [ ]        | 1.22 | Low |
| Errors | [`erris`](#erris) | [ ]        | 1.13 | Critical |
| Errors | [`errorsjoin`](#errorsjoin) | [ ]        | 1.20 | High |
| Time | [`timesince`](#timesince) | [ ]        | 1.0 | High |
| Time | [`timeuntil`](#timeuntil) | [ ]        | 1.8 | Medium |
| Context | [`testingcontext`](#testingcontext) | [x]        | 1.24 | High |
| Context | [`contextcause`](#contextcause) | [ ]        | 1.20 | Medium |
| Context | [`contextafterfunc`](#contextafterfunc) | [ ]        | 1.21 | Medium |
| Context | [`contexttimeoutcause`](#contexttimeoutcause) | [ ]        | 1.21 | Low |
| Sync | [`waitgroup`](#waitgroup) | [x]        | 1.25 | High |
| Sync | [`synconcefunc`](#synconcefunc) | [ ]        | 1.21 | Medium |
| Sync | [`atomicvalue`](#atomicvalue) | [ ]        | 1.19 | Medium |
| Fmt | [`fmtappendf`](#fmtappendf) | [x]        | 1.19 | Medium |
| Testing | [`bloop`](#bloop) | [x]        | 1.25 | Medium |
| JSON | [`omitzero`](#omitzero) | [x]        | 1.24 | Medium |
| Utilities | [`cmpor`](#cmpor) | [ ]        | 1.22 | High |
| Utilities | [`newliteral`](#newliteral) | [x]        | 1.26 | High |
| HTTP | [`httpmux`](#httpmux) | [ ]        | 1.22 | Medium |


---

## 1. Collections (slices, maps)

<a id="slicescontains"></a>

### [x] 1.1 `slicescontains` — replace search loop with `slices.Contains`

**Go 1.21+ | Impact: Critical**

One of the most common patterns — checking if element exists in slice.

**Before:**
```go
found := false
for _, v := range s {
    if v == needle {
        found = true
        break
    }
}
```

**After:**
```go
found := slices.Contains(s, needle)
```

---

<a id="sortslice"></a>

### [x] 1.2 `sortslice` — replace `sort.Slice` with `slices.Sort`

**Go 1.21+ | Impact: High**

Sorting slices is a common operation, new API is simpler and faster.

**Before:**
```go
sort.Slice(users, func(i, j int) bool {
    return users[i].Name < users[j].Name
})
```

**After:**
```go
slices.SortFunc(users, func(a, b User) int {
    return cmp.Compare(a.Name, b.Name)
})
```

---

<a id="minmax"></a>

### [x] 1.3 `minmax` — replace conditional assignment with `min`/`max`

**Go 1.21+ | Impact: High**

Very common pattern in any code with numbers.

**Before:**
```go
if a < b {
    result = a
} else {
    result = b
}
```

**After:**
```go
result := min(a, b)
```

---

<a id="mapsloop"></a>

### [ ] 1.4 `mapsclone` — explicit use of `maps.Clone`

**Go 1.21+ | Impact: Medium**

Extension of `mapsloop` for more precise recognition of copy pattern.

**Before:**
```go
copy := make(map[string]int, len(m))
for k, v := range m {
    copy[k] = v
}
```

**After:**
```go
copy := maps.Clone(m)
```

---

<a id="mapkeysvalues"></a>

### [ ] 1.5 `mapkeysvalues` — use `maps.Keys` / `maps.Values`

**Go 1.23+ | Impact: High**

Very common pattern — get list of keys or values from map.

**Before (keys):**
```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
```

**After:**
```go
keys := slices.Collect(maps.Keys(m))
```

**Before (values):**
```go
vals := make([]User, 0, len(m))
for _, v := range m {
    vals = append(vals, v)
}
```

**After:**
```go
vals := slices.Collect(maps.Values(m))
```

---

<a id="rangeoverindex"></a>

### [ ] 1.6 `rangeoverindex` — replace manual index search with `slices.Index`

**Go 1.21+ | Impact: Medium**

**Before:**
```go
idx := -1
for i, v := range s {
    if v == needle {
        idx = i
        break
    }
}
```

**After:**
```go
idx := slices.Index(s, needle)
```

---

<a id="slicesmaxmin"></a>

### [ ] 1.7 `slicesmaxmin` — use `slices.Max` / `slices.Min`

**Go 1.21+ | Impact: Medium**

**Before:**
```go
max := nums[0]
for _, n := range nums[1:] {
    if n > max {
        max = n
    }
}
```

**After:**
```go
max := slices.Max(nums)
```

---

<a id="slicesreverse"></a>

### [ ] 1.8 `slicesreverse` — use `slices.Reverse`

**Go 1.21+ | Impact: Medium**

**Before:**
```go
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
}
```

**After:**
```go
slices.Reverse(s)
```

---

<a id="mapdeletefunc"></a>

### [ ] 1.9 `mapdeletefunc` — use `maps.DeleteFunc`

**Go 1.21+ | Impact: Low**

**Before:**
```go
for k, v := range m {
    if shouldDelete(k, v) {
        delete(m, k)
    }
}
```

**After:**
```go
maps.DeleteFunc(m, shouldDelete)
```

---

<a id="slicescompact"></a>

### [ ] 1.10 `slicescompact` — remove consecutive duplicates

**Go 1.21+ | Impact: Low**

Specific pattern, but when found — obvious replacement.

**Before:**
```go
result := s[:1]
for i := 1; i < len(s); i++ {
    if s[i] != s[i-1] {
        result = append(result, s[i])
    }
}
```

**After:**
```go
result := slices.Compact(s)
```

---

<a id="clearcollection"></a>

### [ ] 1.11 `clearcollection` — use `clear()` for clearing

**Go 1.21+ | Impact: Medium**

**Before (map):**
```go
for k := range m {
    delete(m, k)
}
```

**After:**
```go
clear(m)
```

**Before (slice — zeroing elements):**
```go
for i := range s {
    s[i] = zero
}
```

**After:**
```go
clear(s)
```

---

## 2. Strings and bytes

<a id="stringscutprefix"></a>

### [x] 2.1 `stringscutprefix` — replace `HasPrefix` + `TrimPrefix` with `CutPrefix`

**Go 1.20+ | Impact: High**

Very common pattern when parsing strings.

**Before:**
```go
if strings.HasPrefix(s, "prefix:") {
    rest := strings.TrimPrefix(s, "prefix:")
    // use rest
}
```

**After:**
```go
if rest, ok := strings.CutPrefix(s, "prefix:"); ok {
    // use rest
}
```

---

<a id="stringsseq"></a>

### [x] 2.2 `stringsseq` — replace `Split`/`Fields` with `SplitSeq`/`FieldSeq`

**Go 1.24+ | Impact: High**

Avoids intermediate slice allocation when iterating.

**Before:**
```go
for _, part := range strings.Split(s, ",") {
    process(part)
}
```

**After:**
```go
for part := range strings.SplitSeq(s, ",") {
    process(part)
}
```

---

<a id="stringsclone"></a>

### [ ] 2.3 `stringsclone` — use `strings.Clone` / `bytes.Clone`

**Go 1.20+ | Impact: Medium**

**Before (strings):**
```go
s2 := string([]byte(s))
// or
s2 := "" + s
```

**After:**
```go
s2 := strings.Clone(s)
```

**Before (bytes):**
```go
b2 := append([]byte(nil), b...)
```

**After:**
```go
b2 := bytes.Clone(b)
```

---

<a id="bytescut"></a>

### [ ] 2.4 `bytescut` — use `bytes.Cut`

**Go 1.18+ | Impact: Medium**

**Before:**
```go
idx := bytes.IndexByte(b, ':')
if idx >= 0 {
    key, value := b[:idx], b[idx+1:]
    // ...
}
```

**After:**
```go
if key, value, ok := bytes.Cut(b, []byte{':'}); ok {
    // ...
}
```

---

## 3. Loops and iteration

<a id="rangeint"></a>

### [x] 3.1 `rangeint` — replace 3-clause loops with `for i := range n`

**Go 1.22+ | Impact: Critical**

One of the most common patterns in Go code.

**Before:**
```go
for i := 0; i < n; i++ {
    // ...
}
```

**After:**
```go
for i := range n {
    // ...
}
```

---

<a id="forvar"></a>

### [x] 3.2 `forvar` — remove `x := x` variables in loops

**Go 1.22+ | Impact: High**

Previously required for capturing variable in closure, no longer needed.

**Before:**
```go
for _, v := range items {
    v := v // no longer needed
    go func() {
        process(v)
    }()
}
```

**After:**
```go
for _, v := range items {
    go func() {
        process(v)
    }()
}
```

---

## 4. Types and interfaces

<a id="efaceany"></a>

### [x] 4.1 `efaceany` — replace `interface{}` with `any`

**Go 1.18+ | Impact: Critical**

Found in huge amounts of legacy code.

**Before:**
```go
func Process(data interface{}) interface{} {
    // ...
}
```

**After:**
```go
func Process(data any) any {
    // ...
}
```

---

<a id="reflecttypefor"></a>

### [ ] 4.2 `reflecttypefor` — use `reflect.TypeFor`

**Go 1.22+ | Impact: Low**

Specific to code with reflection.

**Before:**
```go
t := reflect.TypeOf((*MyType)(nil)).Elem()
```

**After:**
```go
t := reflect.TypeFor[MyType]()
```

---

## 5. Error handling

<a id="erris"></a>

### [ ] 5.1 `erris` — use `errors.Is` / `errors.As`

**Go 1.13+ | Impact: Critical**

One of the most common patterns. Direct error comparison is an anti-pattern.

**Before:**
```go
if err == os.ErrNotExist {
    return nil
}
```

**After:**
```go
if errors.Is(err, os.ErrNotExist) {
    return nil
}
```

**Before (type assert):**
```go
if e, ok := err.(*os.PathError); ok {
    // use e
}
```

**After:**
```go
var e *os.PathError
if errors.As(err, &e) {
    // use e
}
```

---

<a id="errorsjoin"></a>

### [ ] 5.2 `errorsjoin` — use `errors.Join` for multiple errors

**Go 1.20+ | Impact: High**

Very common pattern when multiple operations can fail.

**Before:**
```go
var errs []error
if err := op1(); err != nil {
    errs = append(errs, err)
}
if err := op2(); err != nil {
    errs = append(errs, err)
}
if len(errs) > 0 {
    return fmt.Errorf("multiple errors: %v", errs)
}
```

**After:**
```go
var errs []error
if err := op1(); err != nil {
    errs = append(errs, err)
}
if err := op2(); err != nil {
    errs = append(errs, err)
}
if len(errs) > 0 {
    return errors.Join(errs...)
}
```

Works with `errors.Is`/`errors.As` — they check all wrapped errors.

---

## 6. Time

<a id="timesince"></a>

### [ ] 6.1 `timesince` — use `time.Since`

**Go 1.0+ | Impact: High**

Very common pattern for measuring execution time.

**Before:**
```go
elapsed := time.Now().Sub(start)
```

**After:**
```go
elapsed := time.Since(start)
```

---

<a id="timeuntil"></a>

### [ ] 6.2 `timeuntil` — use `time.Until`

**Go 1.8+ | Impact: Medium**

**Before:**
```go
timeout := deadline.Sub(time.Now())
```

**After:**
```go
timeout := time.Until(deadline)
```

---

## 7. Context

<a id="testingcontext"></a>

### [x] 7.1 `testingcontext` — replace `context.WithCancel` with `t.Context` in tests

**Go 1.24+ | Impact: High**

Common pattern in tests.

**Before:**
```go
func TestFoo(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    // ...
}
```

**After:**
```go
func TestFoo(t *testing.T) {
    ctx := t.Context()
    // ...
}
```

---

<a id="contextcause"></a>

### [ ] 7.2 `contextcause` — use `context.WithCancelCause`

**Go 1.20+ | Impact: Medium**

Improves debugging by providing cancellation reason.

**Before:**
```go
ctx, cancel := context.WithCancel(parent)
// ...
cancel()
// ctx.Err() just returns "context canceled"
```

**After:**
```go
ctx, cancel := context.WithCancelCause(parent)
// ...
cancel(fmt.Errorf("shutting down: %w", reason))
// context.Cause(ctx) returns the actual reason
```

---

<a id="contextafterfunc"></a>

### [ ] 7.3 `contextafterfunc` — use `context.AfterFunc`

**Go 1.21+ | Impact: Medium**

**Before:**
```go
go func() {
    <-ctx.Done()
    cleanup()
}()
```

**After:**
```go
context.AfterFunc(ctx, cleanup)
```

---

<a id="contexttimeoutcause"></a>

### [ ] 7.4 `contexttimeoutcause` — use `context.WithTimeoutCause`

**Go 1.21+ | Impact: Low**

Improves diagnostics but not critical.

**Before:**
```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
```

**After:**
```go
ctx, cancel := context.WithTimeoutCause(parent, 5*time.Second, errors.New("operation timeout"))
```

---

## 8. Concurrency (sync)

<a id="waitgroup"></a>

### [x] 8.1 `waitgroup` — use `WaitGroup.Go`

**Go 1.25+ | Impact: High**

Simplifies typical WaitGroup pattern.

**Before:**
```go
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func() {
        defer wg.Done()
        process(item)
    }()
}
wg.Wait()
```

**After:**
```go
var wg sync.WaitGroup
for _, item := range items {
    wg.Go(func() {
        process(item)
    })
}
wg.Wait()
```

---

<a id="synconcefunc"></a>

### [ ] 8.2 `synconcefunc` — migrate to `sync.OnceFunc` / `sync.OnceValue`

**Go 1.21+ | Impact: Medium**

**Before:**
```go
var (
    once sync.Once
    cfg  *Config
)

func GetConfig() *Config {
    once.Do(func() {
        cfg = loadConfig()
    })
    return cfg
}
```

**After:**
```go
var GetConfig = sync.OnceValue(loadConfig)
```

---

<a id="atomicvalue"></a>

### [ ] 8.3 `atomicvalue` — use typed `atomic.Pointer[T]` / `atomic.Int64`

**Go 1.19+ | Impact: Medium**

Type-safe atomic operations without type assertions.

**Before:**
```go
var val atomic.Value
val.Store(myData)
x := val.Load().(*MyData)  // type assertion, can panic
```

**After:**
```go
var ptr atomic.Pointer[MyData]
ptr.Store(myData)
x := ptr.Load()  // type-safe, no assertion needed
```

Also: `atomic.Int32`, `atomic.Int64`, `atomic.Uint32`, `atomic.Uint64`, `atomic.Bool`.

---

## 9. Formatting (fmt)

<a id="fmtappendf"></a>

### [x] 9.1 `fmtappendf` — replace `[]byte(fmt.Sprintf(...))` with `fmt.Appendf`

**Go 1.19+ | Impact: Medium**

**Before:**
```go
buf := []byte(fmt.Sprintf("Hello, %s!", name))
```

**After:**
```go
buf := fmt.Appendf(nil, "Hello, %s!", name)
```

---

## 10. Testing

<a id="bloop"></a>

### [x] 10.1 `bloop` — replace loops in benchmarks with `b.Loop()`

**Go 1.25+ | Impact: Medium**

**Before:**
```go
func BenchmarkFoo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        foo()
    }
}
```

**After:**
```go
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        foo()
    }
}
```

---

## 11. JSON (encoding/json)

<a id="omitzero"></a>

### [x] 11.1 `omitzero` — replace `omitempty` with `omitzero`

**Go 1.24+ | Impact: Medium**

More correct behavior for zero values.

**Before:**
```go
type Config struct {
    Timeout time.Duration `json:"timeout,omitempty"`
}
```

**After:**
```go
type Config struct {
    Timeout time.Duration `json:"timeout,omitzero"`
}
```

---

## 12. Utilities (cmp)

<a id="cmpor"></a>

### [ ] 12.1 `cmpor` — use `cmp.Or` for default values

**Go 1.22+ | Impact: High**

Very common pattern — fallback to default value.

**Before:**
```go
name := os.Getenv("NAME")
if name == "" {
    name = "default"
}
```

**After:**
```go
name := cmp.Or(os.Getenv("NAME"), "default")
```

**Before (chain):**
```go
var result string
if a != "" {
    result = a
} else if b != "" {
    result = b
} else {
    result = "default"
}
```

**After:**
```go
result := cmp.Or(a, b, "default")
```

---

<a id="newliteral"></a>

### [ ] 12.2 `newliteral` — use `new(value)` for pointer to literal

**Go 1.26+ (upcoming) | Impact: High**

Very common pattern — getting pointer to a literal value.

**Before:**
```go
x := 42
p := &x
// or one-liner workaround:
p := func() *int { v := 42; return &v }()
```

**After:**
```go
p := new(42)  // *int with value 42
```

Works for any type: `new("hello")`, `new(true)`, `new(MyStruct{...})`.

---

## 13. HTTP

<a id="httpmux"></a>

### [ ] 13.1 `httpmux` — new routing patterns in `http.ServeMux`

**Go 1.22+ | Impact: Medium**

**Before:**
```go
mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", 405)
        return
    }
    id := strings.TrimPrefix(r.URL.Path, "/users/")
    // ...
})
```

**After:**
```go
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    // ...
})
```
