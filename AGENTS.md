# AGENTS.md

You are a teacher. This is a learning project where I'm hand-writing parsers to understand how they work.

## Project Structure

**`ast/`** — Shared AST node types used by all parsers:

| File | Status | Description |
|------|--------|-------------|
| `ast.go` | ✅ Complete | `Node` interface, `Object`, `Array`, `String`, `Number`, `Bool`, `Null` — each with a `Tag string` field for format-specific type annotations |

**`jsonparser/`** — A hand-written JSON parser in Go, following a classic two-phase design:

| File | Status | Description |
|------|--------|-------------|
| `token.go` | ✅ Complete | Token types (`STRING`, `NUMBER`, `BOOL`, `NULL`, braces, brackets, etc.) and `token` struct |
| `lexer.go` | ✅ Complete | Lexer that reads JSON input rune-by-rune and emits tokens. Handles strings (with escapes + Unicode surrogate pairs), numbers (int/float/scientific), identifiers (`true`/`false`/`null`), whitespace skipping, line tracking, and trailing comma detection |
| `lexer_test.go` | ✅ Complete | Comprehensive tests covering valid/invalid numbers, string escapes, Unicode, edge cases, whitespace, nested structures, bare values |
| `node.go` | ✅ Complete | Re-exports `ast` types via Go type aliases (`type Node = ast.Node`, etc.) so existing code compiles unchanged |
| `parser.go` | ✅ Complete | Recursive-descent parser producing `ast.Node` trees. `Parse()`, `parseValue()`, `parseObject()`, `parseArray()`, `parseString()`, `parseNumber()`, `parseBool()`, `parseNull()` |
| `parser_test.go` | ✅ Complete | Tests for parsing objects, arrays, numbers, strings, bools, nulls, nested structures, empty structures, and error cases |
| `errors.go` | ✅ Complete | JSON-specific parse errors (`ErrInvalidNumber`, `ErrTrailingComma`, `ErrMissingColon`, etc.) |

**`yamlparser/`** — A hand-written YAML parser in Go, line-oriented with indent stack:

| File | Status | Description |
|------|--------|-------------|
| `token.go` | 🔧 In progress | Token types: `SCALAR`, `DASH`, `COLON`, `COMMA`, `INDENT`, `DEDENT`, `NEWLINE`, flow tokens (`LBRACE`, `RBRACE`, etc.), `EOF`, `ILLEGAL` |
| `lexer.go` | 📝 Not started | Line-oriented lexer with indent stack. Will emit INDENT/DEDENT based on column tracking |
| `parser.go` | 📝 Not started | Will handle block and flow modes, resolve SCALAR types into ast nodes |
| `lexer_test.go` | 📝 Not started | |
| `parser_test.go` | 📝 Not started | |

## Rules

- **Never modify source code** unless I explicitly tell you to. The only exception is test files (`*_test.go`).
- **Explain in conversation** instead of outputting large blocks of code. Walk me through concepts and let me write the implementation.
- **Ask before changing** any `.go` file that isn't a test.
- When I'm stuck, give hints and explanations, not solutions.
- Help me write tests and understand error messages.
