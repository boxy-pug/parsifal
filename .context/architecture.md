# Architecture Decisions

## Shared AST Package (`ast/`)

- `ast/` is a leaf package at the project root. No parser imports it recursively.
- All parsers (`jsonparser`, `yamlparser`, `tomlparser`) produce the same `Node` tree types.
- Each node struct has a `Tag string` field for format-specific type annotations (e.g., YAML's `!!str`, `!Custom`). Empty when unused.
- `jsonparser/node.go` contains type aliases (`type Node = ast.Node`) so existing JSON parser code compiles unchanged.
- Errors stay format-specific. `jsonparser/errors.go` is not shared — YAML and TOML will have their own error sets.

## YAML Lexer Design

### Line-oriented, not rune-oriented

JSON can be lexed rune-by-rune because structure is explicit (`{`, `}`, `[`, `]`). YAML block style is fundamentally line-oriented — indentation is structural. The lexer reasons about entire lines, measures leading spaces, and compares to the indent stack.

### Indent stack

The lexer maintains `indentStack []int` storing actual column positions (not level numbers). Root is always column 0.

- New line more indented than stack top → emit `INDENT`, push column
- New line same indent → same level, no token
- New line less indented → pop stack, emit `DEDENT` for each level closed. If the new column doesn't match any value in the stack, that's an error.

### Coarse tokens, parser resolves types

Unlike the JSON lexer which classifies `true` → `BOOL` and `42` → `NUMBER`, the YAML lexer emits `SCALAR` for all values (quoted or unquoted). The parser resolves scalar types (string, number, bool, null) during parsing because YAML scalar resolution is contextual.

### Token set

| Token | Purpose |
|-------|---------|
| `SCALAR` | Any scalar value (plain, single-quoted, double-quoted, block). Raw string in `Literal`. |
| `DASH` | Block sequence entry (`- ` at start of line). Standalone token — the parser handles nesting via indent context. |
| `COLON` | Key-value separator. |
| `COMMA` | Flow style only. |
| `INDENT` / `DEDENT` | Structural nesting signals (like `{` and `}` in JSON). |
| `NEWLINE` | Line terminator. |
| `LBRACE` / `RBRACE` / `LBRACKET` / `RBRACKET` | Flow style. |
| `EOF` / `ILLEGAL` | Standard. |

### No `KEY` token

The lexer does not emit `KEY`. Keys are just `SCALAR` tokens that happen to be followed by `COLON`. The parser infers "this scalar is a key" from the `SCALAR COLON` pattern. This avoids peek-ahead logic in the lexer and supports non-string keys (e.g., `123: value`).

### Block vs. flow mode

YAML supports mixing block and flow style freely. The parser switches between modes. Flow style (`{}`, `[]`) is tokenized like JSON. Block style uses the indent stack.

### Type resolution

Single-pass: the parser calls a `resolveScalar(raw string, quoted bool) ast.Node` function when it needs a typed value. No second pass over the input.

## Future additions

- YAML: `!!timestamp`, `!!binary`, `!!set`, `!!omap` → new AST node types (`Timestamp`, `Binary`, `Set`, `OrderedMap`).
- YAML: anchors (`&`) and aliases (`*`) → symbol table during parsing.
- TOML: datetime variants (offset datetime, local datetime, local date, local time).
