<!-- Generated from internal/guidelines/guidelines.json; DO NOT EDIT. -->

# Modern Go Guidelines Explained

This file provides a more detailed description of the features supported in the Go Modern Guidelines.

**Modernizer legend:**

- [x] — supported by a modernize analyzer
- [ ] — no modernizer available

**Impact legend:**

- Critical — found in almost every project, dozens of occurrences
- High — found often, 5–20 occurrences per project
- Medium — found regularly, 1–5 occurrences per project
- Low — found rarely or in specific code

## Guidelines

| Category | Guideline | Modernizer | Go | Impact |
|----------|-----------|------------|----|--------|
| Types | [`generic_methods`](#generic_methods) | [ ] | 1.27 | Medium |
| JSON | [`json_v2`](#json_v2) | [ ] | 1.27 | High |
| Types | [`promoted_field_literals`](#promoted_field_literals) | [x] | 1.27 | Medium |
| Strings | [`strings_bytes_cut_last`](#strings_bytes_cut_last) | [x] | 1.27 | Medium |
| Utilities | [`stdlib_uuid`](#stdlib_uuid) | [ ] | 1.27 | Medium |
| URL | [`url_clone`](#url_clone) | [ ] | 1.27 | Low |
| Utilities | [`new_expression`](#new_expression) | [x] | 1.26 | High |
| Errors | [`errors_as_type`](#errors_as_type) | [x] | 1.26 | Medium |
| Sync | [`sync_waitgroup_go`](#sync_waitgroup_go) | [x] | 1.25 | High |
| Context | [`testing_t_context`](#testing_t_context) | [x] | 1.24 | High |
| JSON | [`json_omitzero`](#json_omitzero) | [x] | 1.24 | Medium |
| Testing | [`testing_b_loop`](#testing_b_loop) | [x] | 1.24 | Medium |
| Strings | [`strings_split_seq`](#strings_split_seq) | [x] | 1.24 | High |
| Collections | [`maps_keys_values_iter`](#maps_keys_values_iter) | [ ] | 1.23 | High |
| Collections | [`slices_collect`](#slices_collect) | [ ] | 1.23 | Medium |
| Collections | [`slices_sorted`](#slices_sorted) | [ ] | 1.23 | Medium |
| Time | [`time_tick_gc`](#time_tick_gc) | [ ] | 1.23 | Low |
| Loops | [`range_over_int`](#range_over_int) | [x] | 1.22 | Critical |
| Loops | [`loopvar_capture`](#loopvar_capture) | [x] | 1.22 | High |
| Utilities | [`cmp_or`](#cmp_or) | [ ] | 1.22 | High |
| Types | [`reflect_type_for`](#reflect_type_for) | [x] | 1.22 | Low |
| HTTP | [`http_servemux_patterns`](#http_servemux_patterns) | [ ] | 1.22 | Medium |
| Collections | [`min_max`](#min_max) | [x] | 1.21 | High |
| Collections | [`clear`](#clear) | [ ] | 1.21 | Medium |
| Collections | [`slices_contains`](#slices_contains) | [x] | 1.21 | Critical |
| Collections | [`slices_index`](#slices_index) | [ ] | 1.21 | Medium |
| Collections | [`slices_index_func`](#slices_index_func) | [ ] | 1.21 | Medium |
| Collections | [`slices_sort_func`](#slices_sort_func) | [ ] | 1.21 | High |
| Collections | [`slices_sort`](#slices_sort) | [x] | 1.21 | High |
| Collections | [`slices_max_min`](#slices_max_min) | [ ] | 1.21 | Medium |
| Collections | [`slices_reverse`](#slices_reverse) | [ ] | 1.21 | Medium |
| Collections | [`slices_compact`](#slices_compact) | [ ] | 1.21 | Low |
| Collections | [`slices_clip`](#slices_clip) | [x] | 1.21 | Low |
| Collections | [`slices_clone`](#slices_clone) | [x] | 1.21 | Medium |
| Collections | [`maps_clone`](#maps_clone) | [ ] | 1.21 | Medium |
| Collections | [`maps_copy`](#maps_copy) | [x] | 1.21 | Medium |
| Collections | [`maps_delete_func`](#maps_delete_func) | [ ] | 1.21 | Low |
| Sync | [`sync_once_func`](#sync_once_func) | [ ] | 1.21 | Medium |
| Sync | [`sync_once_value`](#sync_once_value) | [ ] | 1.21 | Medium |
| Context | [`context_after_func`](#context_after_func) | [ ] | 1.21 | Medium |
| Context | [`context_timeout_deadline_cause`](#context_timeout_deadline_cause) | [ ] | 1.21 | Low |
| Strings | [`bytes_clone`](#bytes_clone) | [x] | 1.20 | Medium |
| Strings | [`strings_cut_prefix_suffix`](#strings_cut_prefix_suffix) | [x] | 1.20 | High |
| Errors | [`errors_join`](#errors_join) | [ ] | 1.20 | High |
| Context | [`context_cancel_cause`](#context_cancel_cause) | [ ] | 1.20 | Medium |
| Fmt | [`fmt_appendf`](#fmt_appendf) | [x] | 1.19 | Medium |
| Sync | [`atomic_types`](#atomic_types) | [x] | 1.19 | Medium |
| Types | [`any`](#any) | [x] | 1.18 | Critical |
| Strings | [`bytes_cut`](#bytes_cut) | [x] | 1.18 | Medium |
| Strings | [`strings_clone`](#strings_clone) | [ ] | 1.18 | Medium |
| Strings | [`strings_cut`](#strings_cut) | [x] | 1.18 | Medium |
| Errors | [`errors_is`](#errors_is) | [ ] | 1.13 | Critical |
| Time | [`time_until`](#time_until) | [ ] | 1.8 | Medium |
| Time | [`time_since`](#time_since) | [ ] | 1.0 | High |

---

<a id="generic_methods"></a>

## `generic_methods`

**Category: Types · Go 1.27+ · Impact: Medium · Modernizer: no**

Use generic methods instead of package-level generic helper functions when the operation naturally belongs to the type itself.

Generic methods keep operations in the namespace of the type that owns them. Keep package-level helpers for operations that do not naturally belong to one receiver type.

### Example

**Before:**

```go
type Set[T comparable] map[T]struct{}

func Map[T comparable, U any](s Set[T], f func(T) U) []U {
	out := make([]U, 0, len(s))
	for value := range s {
		out = append(out, f(value))
	}
	return out
}

names := Map(users, func(user User) string {
	return user.Name
})
```

**After:**

```go
type Set[T comparable] map[T]struct{}

func (s Set[T]) Map[U any](f func(T) U) []U {
	out := make([]U, 0, len(s))
	for value := range s {
		out = append(out, f(value))
	}
	return out
}

names := users.Map(func(user User) string {
	return user.Name
})
```

---

<a id="json_v2"></a>

## `json_v2`

**Category: JSON · Go 1.27+ · Impact: High · Modernizer: no**

Use `encoding/json/v2` for new JSON code in Go 1.27+; leave existing `encoding/json` code unchanged unless migration is explicitly requested.

`encoding/json/v2` has stricter, more interoperable defaults: it rejects invalid UTF-8 and duplicate object names, encodes `nil` slices and maps as empty arrays and objects, and reports invalid JSON-related Go types. Prefer these defaults in new code. Do not replace `encoding/json` in existing code unless migration is explicitly requested, because even a compiling import change can alter wire behavior. For an explicit migration, begin with `jsonv1.DefaultOptionsV1()`, test serialized output and accepted input, then remove or override compatibility options deliberately; later options take precedence.

### Example 1

**Before:**

```go
import "encoding/json"

type Pet struct {
	Name      string
	Nicknames []string
}

body, err := json.Marshal(Pet{Name: "Remi"})
// body is {"Name":"Remi","Nicknames":null}
```

**After:**

```go
import "encoding/json/v2"

type Pet struct {
	Name      string
	Nicknames []string
}

// New code uses v2's stricter defaults.
body, err := json.Marshal(Pet{Name: "Remi"})
// body is {"Name":"Remi","Nicknames":[]}
```

### Example 2

**Before:**

```go
import "encoding/json"

type Pet struct {
	Name      string
	Nicknames []string
}

// Leave existing v1 code unchanged unless migration is requested.
body, err := json.Marshal(Pet{Name: "Remi"})
// Existing consumers receive {"Name":"Remi","Nicknames":null}.
```

**After:**

```go
import (
	jsonv1 "encoding/json"
	json "encoding/json/v2"
)

type Pet struct {
	Name      string
	Nicknames []string
}

// For an explicit migration, retain existing wire behavior first.
body, err := json.Marshal(
	Pet{Name: "Remi"},
	jsonv1.DefaultOptionsV1(),
)
// Existing consumers still receive {"Name":"Remi","Nicknames":null}.
```

### Example 3

**Before:**

```go
import (
	jsonv1 "encoding/json"
	json "encoding/json/v2"
)

type Pet struct {
	Name      string
	Nicknames []string
}

body, err := json.Marshal(
	Pet{Name: "Remi"},
	jsonv1.DefaultOptionsV1(),
)
// body is {"Name":"Remi","Nicknames":null}
```

**After:**

```go
import (
	jsonv1 "encoding/json"
	json "encoding/json/v2"
)

type Pet struct {
	Name      string
	Nicknames []string
}

// Then adopt v2 behavior only after verifying consumers.
body, err := json.Marshal(
	Pet{Name: "Remi"},
	jsonv1.DefaultOptionsV1(),
	json.FormatNilSliceAsNull(false),
)
// body is {"Name":"Remi","Nicknames":[]}
```

---

<a id="promoted_field_literals"></a>

## `promoted_field_literals`

**Category: Types · Go 1.27+ · Impact: Medium · Modernizer: yes**

Set embedded struct fields directly with promoted field names in Go 1.27+ struct literals instead of constructing the embedded struct explicitly.

Go 1.27 allows keyed struct literals to set fields promoted from embedded struct types directly. Use `Field: value` in the outer literal instead of `Embedded: Embedded{Field: value}` when you are just setting promoted fields. Do not mix a promoted field with the embedded field that promotes it; pointer-embedded paths are not supported.

### Example

**Before:**

```go
type AuditInfo struct {
	CreatedBy string
	UpdatedBy string
}

type Document struct {
	AuditInfo
	Name string
	Path string
}

doc := Document{
	AuditInfo: AuditInfo{
		CreatedBy: "alice",
		UpdatedBy: "alice",
	},
	Name: "report.pdf",
	Path: "/documents/report.pdf",
}
```

**After:**

```go
type AuditInfo struct {
	CreatedBy string
	UpdatedBy string
}

type Document struct {
	AuditInfo
	Name string
	Path string
}

doc := Document{
	CreatedBy: "alice",
	UpdatedBy: "alice",
	Name:      "report.pdf",
	Path:      "/documents/report.pdf",
}
```

---

<a id="strings_bytes_cut_last"></a>

## `strings_bytes_cut_last`

**Category: Strings · Go 1.27+ · Impact: Medium · Modernizer: yes**

Use `strings.CutLast` and `bytes.CutLast` instead of `LastIndex` plus manual slicing around the last separator.

`CutLast` splits around the final separator and returns before, after, and found. It removes `LastIndex` checks and manual separator-length slicing.

### Example 1

**Before:**

```go
i := strings.LastIndex(path, "/")
if i < 0 {
	return "", path, false
}
dir, file := path[:i], path[i+1:]
```

**After:**

```go
dir, file, found := strings.CutLast(path, "/")
```

### Example 2

**Before:**

```go
i := bytes.LastIndex(line, []byte(":"))
if i < 0 {
	return nil, line, false
}
name, value := line[:i], line[i+1:]
```

**After:**

```go
name, value, found := bytes.CutLast(line, []byte(":"))
```

---

<a id="stdlib_uuid"></a>

## `stdlib_uuid`

**Category: Utilities · Go 1.27+ · Impact: Medium · Modernizer: no**

Use the standard library `uuid` package instead of third-party libraries or custom UUID implementations when targeting Go 1.27+.

The standard `uuid` package covers common UUID generation and parsing without an external dependency. Keep a third-party package only when targeting older Go versions or needing behavior outside the standard API.

### Example 1

**Before:**

```go
import googleuuid "github.com/google/uuid"

id := googleuuid.New()
text := id.String()
```

**After:**

```go
import "uuid"

id := uuid.New()
text := id.String()
```

### Example 2

**Before:**

```go
id, err := googleuuid.Parse(raw)
if err != nil {
	return err
}
```

**After:**

```go
id, err := uuid.Parse(raw)
if err != nil {
	return err
}
```

---

<a id="url_clone"></a>

## `url_clone`

**Category: URL · Go 1.27+ · Impact: Low · Modernizer: no**

Use the `URL.Clone` and `Values.Clone` methods from `net/url` to copy URLs and `URL` values instead of manual copying.

`URL.Clone` and `Values.Clone` create deep copies using `net/url`'s own representation. They avoid partial manual copies that can accidentally share mutable slices or miss `URL` fields.

### Example 1

**Before:**

```go
copy := new(url.URL)
*copy = *base
if base.User != nil {
	username := base.User.Username()
	if password, ok := base.User.Password(); ok {
		copy.User = url.UserPassword(username, password)
	} else {
		copy.User = url.User(username)
	}
}
```

**After:**

```go
copy := base.Clone()
```

### Example 2

**Before:**

```go
copy := make(url.Values, len(values))
for key, items := range values {
	copy[key] = append([]string(nil), items...)
}
```

**After:**

```go
copy := values.Clone()
```

---

<a id="new_expression"></a>

## `new_expression`

**Category: Utilities · Go 1.26+ · Impact: High · Modernizer: yes**

Use `new(value)` for pointer fields or arguments instead of generic/type-specific pointer helper functions or temporary variables used only for `&value`.

`new(value)` creates a `*T` from a value expression. In struct literals, prefer `Field: new(value)` for pointer fields over helper calls whose only purpose is returning `&value`; keep helpers only when they add behavior.

### Example 1

**Before:**

```go
func Pointer[T any](value T) *T {
	return &value
}

cfg := Config{
	Timeout: Pointer(30),
	Debug:   Pointer(true),
}
```

**After:**

```go
cfg := Config{
	Timeout: new(30),
	Debug:   new(true),
}
```

### Example 2

**Before:**

```go
timeout := 30
debug := true
cfg := Config{
	Timeout: &timeout,
	Debug:   &debug,
}
```

**After:**

```go
cfg := Config{
	Timeout: new(30),
	Debug:   new(true),
}
```

---

<a id="errors_as_type"></a>

## `errors_as_type`

**Category: Errors · Go 1.26+ · Impact: Medium · Modernizer: yes**

Use `errors.AsType[T](err)` when checking whether an error matches a specific type.

`errors.AsType` returns the matched error value and a boolean directly. It avoids a separate temporary variable and the pointer-to-target pattern required by `errors.As`.

### Example

**Before:**

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
	handle(pathErr)
}
```

**After:**

```go
if pathErr, ok := errors.AsType[*os.PathError](err); ok {
	handle(pathErr)
}
```

---

<a id="sync_waitgroup_go"></a>

## `sync_waitgroup_go`

**Category: Sync · Go 1.25+ · Impact: High · Modernizer: yes**

Use `wg.Go` when spawning goroutines tracked by a `sync.WaitGroup`.

`WaitGroup.Go` starts a goroutine and handles the matching `Add` and `Done` calls. Use it when the goroutine's lifetime is exactly what the `WaitGroup` should track.

### Example

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

<a id="testing_t_context"></a>

## `testing_t_context`

**Category: Context · Go 1.24+ · Impact: High · Modernizer: yes**

Use `t.Context()` when a test function needs a context tied to the test lifetime.

`t.Context()` returns a context tied to the test lifetime. It removes manual background context setup when helper work should stop as the test is ending.

### Example

**Before:**

```go
func TestFoo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := doSomething(ctx)
}
```

**After:**

```go
func TestFoo(t *testing.T) {
	ctx := t.Context()
	result := doSomething(ctx)
}
```

---

<a id="json_omitzero"></a>

## `json_omitzero`

**Category: JSON · Go 1.24+ · Impact: Medium · Modernizer: yes**

Use `omitzero` on JSON-tagged bool, numeric, struct, and time fields whose zero value should be omitted; keep `omitempty` for empty strings, slices, and maps.

When adding or editing JSON struct tags, use `omitzero` when the Go zero value means the field is absent, such as `false`, `0`, or a zero time. Use `omitempty` when JSON empty-value semantics are intended, especially empty strings, slices, and maps.

### Example

**Before:**

```go
type CacheEntry struct {
	Name      string    `json:"name,omitempty"`
	Warm      bool      `json:"warm,omitempty"`
	Hits      int64     `json:"hits,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}
```

**After:**

```go
type CacheEntry struct {
	Name      string    `json:"name,omitempty"`
	Warm      bool      `json:"warm,omitzero"`
	Hits      int64     `json:"hits,omitzero"`
	ExpiresAt time.Time `json:"expiresAt,omitzero"`
}
```

---

<a id="testing_b_loop"></a>

## `testing_b_loop`

**Category: Testing · Go 1.24+ · Impact: Medium · Modernizer: yes**

Use `b.Loop()` for the main loop in benchmark functions.

`b.Loop()` is the modern benchmark loop. It manages benchmark iteration mechanics for you and can remove the need for manual timer control in the benchmark body.

### Example

**Before:**

```go
func BenchmarkFoo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doWork()
	}
}
```

**After:**

```go
func BenchmarkFoo(b *testing.B) {
	for b.Loop() {
		doWork()
	}
}
```

---

<a id="strings_split_seq"></a>

## `strings_split_seq`

**Category: Strings · Go 1.24+ · Impact: High · Modernizer: yes**

Use `strings.SplitSeq`, `strings.FieldsSeq`, `bytes.SplitSeq`, or `bytes.FieldsSeq` when iterating over split results.

The `SplitSeq` and `FieldsSeq` helpers stream substrings instead of allocating a slice of all parts.

### Example 1

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

### Example 2

**Before:**

```go
for _, field := range bytes.Fields(b) {
	process(field)
}
```

**After:**

```go
for field := range bytes.FieldsSeq(b) {
	process(field)
}
```

---

<a id="maps_keys_values_iter"></a>

## `maps_keys_values_iter`

**Category: Collections · Go 1.23+ · Impact: High · Modernizer: no**

Use `maps.Keys` or `maps.Values` directly as iterators instead of manually looping over a map.

`maps.Keys` and `maps.Values` can be ranged over directly when a slice is not needed. This avoids allocating a temporary collection just to iterate.

### Example

**Before:**

```go
for k := range m {
	process(k)
}
```

**After:**

```go
for k := range maps.Keys(m) {
	process(k)
}
```

---

<a id="slices_collect"></a>

## `slices_collect`

**Category: Collections · Go 1.23+ · Impact: Medium · Modernizer: no**

Use `slices.Collect` to build a slice from an iterator.

`slices.Collect` materializes an iterator into a slice when a slice is actually required. Prefer ranging over the iterator directly when streaming is enough.

### Example

**Before:**

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

---

<a id="slices_sorted"></a>

## `slices_sorted`

**Category: Collections · Go 1.23+ · Impact: Medium · Modernizer: no**

Use `slices.Sorted` to collect and sort iterator values in one step.

`slices.Sorted` collects iterator values and sorts them in one call. It is useful for deterministic output from unordered sources such as map keys.

### Example

**Before:**

```go
keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
slices.Sort(keys)
```

**After:**

```go
keys := slices.Sorted(maps.Keys(m))
```

---

<a id="time_tick_gc"></a>

## `time_tick_gc`

**Category: Time · Go 1.23+ · Impact: Low · Modernizer: no**

Use `time.Tick` when it fits; Go 1.23 can recover unreferenced tickers without requiring `Stop` for GC.

Go 1.23 made unreferenced tickers recoverable by the garbage collector. Use `time.Tick` for simple forever loops, and use `time.NewTicker` when you need `Stop` or `Reset`.

### Example

**Before:**

```go
ticker := time.NewTicker(time.Second)
defer ticker.Stop()
for range ticker.C {
	poll()
}
```

**After:**

```go
for range time.Tick(time.Second) {
	poll()
}
```

---

<a id="range_over_int"></a>

## `range_over_int`

**Category: Loops · Go 1.22+ · Impact: Critical · Modernizer: yes**

Use `for i := range n` when iterating from `0` to `n-1`.

Ranging over an integer is the concise Go 1.22 form for zero-based count loops. Keep a traditional for loop when you need a non-zero start, custom step, or changing bound.

### Example

**Before:**

```go
for i := 0; i < len(items); i++ {
	process(items[i])
}
```

**After:**

```go
for i := range len(items) {
	process(items[i])
}
```

---

<a id="loopvar_capture"></a>

## `loopvar_capture`

**Category: Loops · Go 1.22+ · Impact: High · Modernizer: yes**

Do not add redundant loop-variable copies before closures or taking addresses; Go 1.22 gives each iteration its own variables.

Go 1.22 gives each loop iteration its own variables, so defensive copies like `v := v` are usually redundant before closures, goroutines, deferred functions, or appending `&v`. Keep an explicit copy only when it serves another purpose. If you need a pointer to the original slice element rather than the per-iteration copy, use `&slice[i]`.

### Example 1

**Before:**

```go
for _, item := range items {
	item := item
	go func() {
		process(item)
	}()
}
```

**After:**

```go
for _, item := range items {
	go func() {
		process(item)
	}()
}
```

### Example 2

**Before:**

```go
var selected []*Item
for _, item := range items {
	item := item
	if item.Enabled {
		selected = append(selected, &item)
	}
}
```

**After:**

```go
var selected []*Item
for _, item := range items {
	if item.Enabled {
		selected = append(selected, &item)
	}
}
```

---

<a id="cmp_or"></a>

## `cmp_or`

**Category: Utilities · Go 1.22+ · Impact: High · Modernizer: no**

Use `cmp.Or` to pick the first non-zero value from a fallback chain.

`cmp.Or` returns the first non-zero value from its arguments. It is concise for simple fallback chains, but remember that all arguments are evaluated before the call.

### Example

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

---

<a id="reflect_type_for"></a>

## `reflect_type_for`

**Category: Types · Go 1.22+ · Impact: Low · Modernizer: yes**

Use `reflect.TypeFor[T]()` instead of `reflect.TypeOf((*T)(nil)).Elem()`.

`reflect.TypeFor` returns the `reflect.Type` for a compile-time type parameter or type argument. It avoids `nil`-pointer tricks and is clearer for interface types.

### Example

**Before:**

```go
typ := reflect.TypeOf((*T)(nil)).Elem()
```

**After:**

```go
typ := reflect.TypeFor[T]()
```

---

<a id="http_servemux_patterns"></a>

## `http_servemux_patterns`

**Category: HTTP · Go 1.22+ · Impact: Medium · Modernizer: no**

Use method-aware `ServeMux` patterns and `r.PathValue` for path parameters.

The modern `ServeMux` pattern syntax can include an HTTP method and named path wildcards. `r.PathValue` retrieves wildcard values without manual path trimming.

### Example

**Before:**

```go
mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/")
	handleID(w, r, id)
})
```

**After:**

```go
mux.HandleFunc("GET /api/{id}", func(w http.ResponseWriter, r *http.Request) {
	handleID(w, r, r.PathValue("id"))
})
```

---

<a id="min_max"></a>

## `min_max`

**Category: Collections · Go 1.21+ · Impact: High · Modernizer: yes**

Use built-in `min` and `max` instead of handwritten comparisons.

The `min` and `max` built-ins express ordered comparisons directly. They remove branch boilerplate when the logic is simply choosing the smaller or larger value.

### Example

**Before:**

```go
if b > a {
	a = b
}
```

**After:**

```go
a = max(a, b)
```

---

<a id="clear"></a>

## `clear`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `clear(m)` to delete all map entries or `clear(s)` to zero slice elements.

`clear` is the idiomatic built-in for removing all map entries or zeroing the elements of a slice. For slices it preserves the slice length and capacity.

### Example

**Before:**

```go
for k := range m {
	delete(m, k)
}
```

**After:**

```go
clear(m)
```

---

<a id="slices_contains"></a>

## `slices_contains`

**Category: Collections · Go 1.21+ · Impact: Critical · Modernizer: yes**

Use `slices.Contains` instead of a manual search loop.

`slices.Contains` is the direct membership test for slices of comparable values. It replaces a manual search loop whose only result is whether the element is present.

### Example

**Before:**

```go
found := false
for _, item := range items {
	if item == x {
		found = true
		break
	}
}
```

**After:**

```go
found := slices.Contains(items, x)
```

---

<a id="slices_index"></a>

## `slices_index`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `slices.Index` to find the index of an element, returning `-1` when absent.

`slices.Index` is the direct index lookup for slices of comparable values. It returns `-1` when the value is absent, matching the common handwritten convention.

### Example

**Before:**

```go
index := -1
for i, item := range items {
	if item == x {
		index = i
		break
	}
}
```

**After:**

```go
index := slices.Index(items, x)
```

---

<a id="slices_index_func"></a>

## `slices_index_func`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `slices.IndexFunc` to find an element by predicate.

`slices.IndexFunc` is the standard predicate-based lookup. Use it when the match depends on fields, derived values, or any condition other than direct equality.

### Example

**Before:**

```go
index := -1
for i, item := range items {
	if item.ID == id {
		index = i
		break
	}
}
```

**After:**

```go
index := slices.IndexFunc(items, func(item Item) bool {
	return item.ID == id
})
```

---

<a id="slices_sort_func"></a>

## `slices_sort_func`

**Category: Collections · Go 1.21+ · Impact: High · Modernizer: no**

Use `slices.SortFunc` with `cmp.Compare` instead of `sort.Slice` for typed comparisons.

`slices.SortFunc` compares typed elements directly, avoiding `sort.Slice` closures that index back into the slice. `cmp.Compare` is the usual comparator for ordered fields.

### Example

**Before:**

```go
sort.Slice(items, func(i, j int) bool {
	return items[i].X < items[j].X
})
```

**After:**

```go
slices.SortFunc(items, func(a, b Item) int {
	return cmp.Compare(a.X, b.X)
})
```

---

<a id="slices_sort"></a>

## `slices_sort`

**Category: Collections · Go 1.21+ · Impact: High · Modernizer: yes**

Use `slices.Sort` for slices of ordered values.

`slices.Sort` is the generic sort for slices of ordered values. It replaces older type-specific helpers and simple `sort.Slice` calls.

### Example

**Before:**

```go
sort.Ints(values)
```

**After:**

```go
slices.Sort(values)
```

---

<a id="slices_max_min"></a>

## `slices_max_min`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `slices.Max` and `slices.Min` instead of manual loops over ordered values.

`slices.Max` and `slices.Min` capture the common non-empty-slice scan for an ordered maximum or minimum. Keep explicit logic when empty slices or custom comparison rules need special handling.

### Example

**Before:**

```go
maxValue := values[0]
for _, value := range values[1:] {
	if value > maxValue {
		maxValue = value
	}
}
```

**After:**

```go
maxValue := slices.Max(values)
```

---

<a id="slices_reverse"></a>

## `slices_reverse`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `slices.Reverse` instead of a manual swap loop.

`slices.Reverse` performs an in-place reversal without hand-written index math. It is clearer and less error-prone than open-coded swap loops.

### Example

**Before:**

```go
for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
	items[i], items[j] = items[j], items[i]
}
```

**After:**

```go
slices.Reverse(items)
```

---

<a id="slices_compact"></a>

## `slices_compact`

**Category: Collections · Go 1.21+ · Impact: Low · Modernizer: no**

Use `slices.Compact` to remove consecutive duplicates in place.

`slices.Compact` removes consecutive duplicate values in place. If the goal is to remove all duplicates, sort or otherwise group equal values first.

### Example

**Before:**

```go
out := values[:0]
for i, value := range values {
	if i == 0 || value != values[i-1] {
		out = append(out, value)
	}
}
values = out
```

**After:**

```go
values = slices.Compact(values)
```

---

<a id="slices_clip"></a>

## `slices_clip`

**Category: Collections · Go 1.21+ · Impact: Low · Modernizer: yes**

Use `slices.Clip` to remove unused capacity.

`slices.Clip` limits a slice's capacity to its length. Use it to stop retaining excess backing-array capacity or to prevent later appends from reusing hidden capacity.

### Example

**Before:**

```go
s = s[:len(s):len(s)]
```

**After:**

```go
s = slices.Clip(s)
```

---

<a id="slices_clone"></a>

## `slices_clone`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: yes**

Use `slices.Clone` to copy a slice.

`slices.Clone` returns a shallow copy of a slice using the standard helper. It is clearer than append-copy boilerplate and preserves `nil` slices.

### Example

**Before:**

```go
copied := append([]T(nil), values...)
```

**After:**

```go
copied := slices.Clone(values)
```

---

<a id="maps_clone"></a>

## `maps_clone`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `maps.Clone` instead of manual map iteration.

`maps.Clone` returns a shallow copy of a map using the standard helper. It preserves `nil` maps and makes the copy operation explicit.

### Example

**Before:**

```go
copied := make(map[string]int, len(src))
for k, v := range src {
	copied[k] = v
}
```

**After:**

```go
copied := maps.Clone(src)
```

---

<a id="maps_copy"></a>

## `maps_copy`

**Category: Collections · Go 1.21+ · Impact: Medium · Modernizer: yes**

Use `maps.Copy` to copy entries from one map into another.

`maps.Copy` copies all entries from one map into another, overwriting existing keys in the destination. It replaces loops whose body only assigns `dst[k] = v`.

### Example

**Before:**

```go
for k, v := range src {
	dst[k] = v
}
```

**After:**

```go
maps.Copy(dst, src)
```

---

<a id="maps_delete_func"></a>

## `maps_delete_func`

**Category: Collections · Go 1.21+ · Impact: Low · Modernizer: no**

Use `maps.DeleteFunc` to delete map entries that match a predicate.

`maps.DeleteFunc` deletes entries selected by a predicate. It keeps the filtering condition in one callback and avoids hand-written delete loops.

### Example

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
maps.DeleteFunc(m, func(k string, v int) bool {
	return shouldDelete(k, v)
})
```

---

<a id="sync_once_func"></a>

## `sync_once_func`

**Category: Sync · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `sync.OnceFunc` instead of `sync.Once` plus a wrapper closure.

`sync.OnceFunc` wraps a function so that the returned function runs it at most once. It is a compact fit for idempotent cleanup or one-time initialization hooks.

### Example

**Before:**

```go
var once sync.Once
cleanup := func() {
	once.Do(func() {
		close(ch)
	})
}
```

**After:**

```go
cleanup := sync.OnceFunc(func() {
	close(ch)
})
```

---

<a id="sync_once_value"></a>

## `sync_once_value`

**Category: Sync · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `sync.OnceValue` to memoize a computed value.

`sync.OnceValue` memoizes a computed value safely across goroutines. It replaces a `sync.Once` plus a separate result variable and getter closure.

### Example

**Before:**

```go
var once sync.Once
var value T
getter := func() T {
	once.Do(func() {
		value = computeValue()
	})
	return value
}
```

**After:**

```go
getter := sync.OnceValue(func() T {
	return computeValue()
})
```

---

<a id="context_after_func"></a>

## `context_after_func`

**Category: Context · Go 1.21+ · Impact: Medium · Modernizer: no**

Use `context.AfterFunc` to run cleanup when a context is canceled.

`context.AfterFunc` registers work to run after cancellation and returns a stop function. It avoids starting a goroutine whose only job is to wait on `ctx.Done()`.

### Example

**Before:**

```go
go func() {
	<-ctx.Done()
	cleanup()
}()
```

**After:**

```go
stop := context.AfterFunc(ctx, cleanup)
defer stop()
```

---

<a id="context_timeout_deadline_cause"></a>

## `context_timeout_deadline_cause`

**Category: Context · Go 1.21+ · Impact: Low · Modernizer: no**

Use timeout and deadline contexts with causes when callers need to inspect the cancellation reason.

Timeout and deadline causes let callers distinguish the reason for cancellation with `context.Cause`. Use them when `ctx.Err()` is too coarse for diagnostics or control flow.

### Example

**Before:**

```go
ctx, cancel := context.WithTimeout(parent, d)
defer cancel()
```

**After:**

```go
ctx, cancel := context.WithTimeoutCause(parent, d, errTimeout)
defer cancel()
```

---

<a id="bytes_clone"></a>

## `bytes_clone`

**Category: Strings · Go 1.20+ · Impact: Medium · Modernizer: yes**

Use `bytes.Clone` to copy a byte slice.

`bytes.Clone` communicates that a new byte slice is required. It avoids append-copy boilerplate and preserves the standard library's `nil`-slice behavior.

### Example

**Before:**

```go
copied := append([]byte(nil), b...)
```

**After:**

```go
copied := bytes.Clone(b)
```

---

<a id="strings_cut_prefix_suffix"></a>

## `strings_cut_prefix_suffix`

**Category: Strings · Go 1.20+ · Impact: High · Modernizer: yes**

Use `strings.CutPrefix` or `strings.CutSuffix` when you need both the trimmed result and whether it matched.

`CutPrefix` and `CutSuffix` combine the match check and trimming operation. They avoid writing `HasPrefix` or `HasSuffix` followed by a second operation that repeats the same condition.

### Example

**Before:**

```go
if strings.HasPrefix(s, "pre:") {
	rest := strings.TrimPrefix(s, "pre:")
	use(rest)
}
```

**After:**

```go
if rest, ok := strings.CutPrefix(s, "pre:"); ok {
	use(rest)
}
```

---

<a id="errors_join"></a>

## `errors_join`

**Category: Errors · Go 1.20+ · Impact: High · Modernizer: no**

Use `errors.Join` to combine multiple errors while preserving error matching.

`errors.Join` combines multiple non-`nil` errors while preserving error matching through `errors.Is` and `errors.As`. It returns `nil` when there are no non-`nil` errors.

### Example

**Before:**

```go
if err1 != nil && err2 != nil {
	return fmt.Errorf("%v; %w", err1, err2)
}
```

**After:**

```go
return errors.Join(err1, err2)
```

---

<a id="context_cancel_cause"></a>

## `context_cancel_cause`

**Category: Context · Go 1.20+ · Impact: Medium · Modernizer: no**

Use `context.WithCancelCause` and `context.Cause` when cancellation needs to carry an error cause.

Cancellation causes let code record why a context was canceled. Callers can inspect `context.Cause(ctx)` instead of only seeing the broad `ctx.Err` result.

### Example

**Before:**

```go
ctx, cancel := context.WithCancel(parent)
cancel()
```

**After:**

```go
ctx, cancel := context.WithCancelCause(parent)
cancel(err)
cause := context.Cause(ctx)
```

---

<a id="fmt_appendf"></a>

## `fmt_appendf`

**Category: Fmt · Go 1.19+ · Impact: Medium · Modernizer: yes**

Use `fmt.Appendf` when appending formatted text to a byte slice and an intermediate `fmt.Sprintf` string is unnecessary.

`fmt.Appendf` writes formatted output directly into a byte slice. It is useful when code is already accumulating `[]byte` and the extra string allocation from `fmt.Sprintf` is not needed.

### Example

**Before:**

```go
buf = append(buf, []byte(fmt.Sprintf("x=%d", x))...)
```

**After:**

```go
buf = fmt.Appendf(buf, "x=%d", x)
```

---

<a id="atomic_types"></a>

## `atomic_types`

**Category: Sync · Go 1.19+ · Impact: Medium · Modernizer: yes**

Use typed atomics such as `atomic.Bool`, `atomic.Int64`, and `atomic.Pointer[T]` instead of untyped atomic functions.

Typed atomic wrapper values keep the storage and the atomic operations together. They make the value type visible, reduce accidental non-atomic access, and avoid old pointer-alignment pitfalls.

### Example 1

**Before:**

```go
var enabled int32
atomic.StoreInt32(&enabled, 1)
if atomic.LoadInt32(&enabled) != 0 {
	run()
}
```

**After:**

```go
var enabled atomic.Bool
enabled.Store(true)
if enabled.Load() {
	run()
}
```

### Example 2

**Before:**

```go
var ptr unsafe.Pointer
atomic.StorePointer(&ptr, unsafe.Pointer(cfg))
```

**After:**

```go
var ptr atomic.Pointer[Config]
ptr.Store(cfg)
```

---

<a id="any"></a>

## `any`

**Category: Types · Go 1.18+ · Impact: Critical · Modernizer: yes**

Use `any` instead of `interface{}`.

`any` is the built-in alias for `interface{}` introduced with generics. Use it for unconstrained values and type parameters so the code reads in current Go terminology.

### Example

**Before:**

```go
func Decode(v interface{}) error {
	return nil
}
```

**After:**

```go
func Decode(v any) error {
	return nil
}
```

---

<a id="bytes_cut"></a>

## `bytes_cut`

**Category: Strings · Go 1.18+ · Impact: Medium · Modernizer: yes**

Use `bytes.Cut` instead of `bytes.Index` plus manual slicing.

`bytes.Cut` returns before, after, and found in one call. It avoids duplicated Index checks and separator-length arithmetic while returning slices of the original byte slice.

### Example

**Before:**

```go
i := bytes.Index(b, sep)
if i < 0 {
	return nil, nil, false
}
before, after := b[:i], b[i+len(sep):]
```

**After:**

```go
before, after, found := bytes.Cut(b, sep)
```

---

<a id="strings_clone"></a>

## `strings_clone`

**Category: Strings · Go 1.18+ · Impact: Medium · Modernizer: no**

Use `strings.Clone` to copy a string without retaining shared backing memory.

`strings.Clone` forces a string copy. Use it when retaining a small string derived from a much larger string would otherwise keep the larger backing memory alive.

### Example

**Before:**

```go
copied := string([]byte(s))
```

**After:**

```go
copied := strings.Clone(s)
```

---

<a id="strings_cut"></a>

## `strings_cut`

**Category: Strings · Go 1.18+ · Impact: Medium · Modernizer: yes**

Use `strings.Cut` instead of `strings.Index` plus manual slicing.

`strings.Cut` returns before, after, and found in one call. It makes split-at-first-separator logic explicit and removes manual Index checks and slicing.

### Example

**Before:**

```go
i := strings.Index(s, ":")
if i < 0 {
	return "", "", false
}
key, value := s[:i], s[i+1:]
```

**After:**

```go
key, value, found := strings.Cut(s, ":")
```

---

<a id="errors_is"></a>

## `errors_is`

**Category: Errors · Go 1.13+ · Impact: Critical · Modernizer: no**

Use `errors.Is(err, target)` instead of `err == target` so wrapped errors are handled correctly.

`errors.Is` walks wrapped errors and honors custom `Is` methods. Direct equality only matches the exact sentinel value and misses errors wrapped with `fmt.Errorf` or custom wrappers.

### Example

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

---

<a id="time_until"></a>

## `time_until`

**Category: Time · Go 1.8+ · Impact: Medium · Modernizer: no**

Use `time.Until(deadline)` instead of `deadline.Sub(time.Now())`.

`time.Until` is the standard spelling for remaining time before a deadline. It is equivalent to `deadline.Sub(time.Now())` but reads in the direction the caller cares about.

### Example

**Before:**

```go
remaining := deadline.Sub(time.Now())
```

**After:**

```go
remaining := time.Until(deadline)
```

---

<a id="time_since"></a>

## `time_since`

**Category: Time · Go 1.0+ · Impact: High · Modernizer: no**

Use `time.Since(start)` instead of `time.Now().Sub(start)`.

`time.Since` is the standard spelling for elapsed time from a recorded start. It keeps the `time.Now` call inside the time package helper and makes the intent easier to scan.

### Example

**Before:**

```go
elapsed := time.Now().Sub(start)
```

**After:**

```go
elapsed := time.Since(start)
```
