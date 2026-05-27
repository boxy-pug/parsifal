# parsifal

A hand-written JSON parser in Go. This is a learning project inspired by [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball.

I'm starting with JSON, and then i'll add YAML and TOML maybe, and a CLI interface maybe.

Built from scratch to understand how parsers and lexers work under the hood. AI is used for concept explanation and test writing, but source code is written by hand.

## Things I've learned:

### About JSON

- json is a data serialization format thats easy to parse for computers and human readable. All keys are strings. All strings are in double quotes. objects{"key": val}, arrays [val, val], string, number, bool, null basically. doesnt care about whitespace.
- Duplicate keys, what happens? json spec doesnt say!
- Numbers are always floats, if representing a big int you should make it a string and parse as int on the other end.

### About Unicode

- Unicode has over a million possible chars, and it just assigns a number for each one basically. One 16bit int \uXXXX, (2 bytes) only gives you about 65000 of them. This is called the "Basic Multilingual Plane" (BMP). To access chars that live above this range you need to use some "tricks". Unicode actually uses 21 bits.
- utf8 and utf16 are two different "packaging formats". utf8 = 1-byte chunks and utf16 = 2-byte chunks. They use different tricks to express chars above the BMP.
- utf8: First byte tells you how many bytes a char is. Then every byte after that starts with a continuation pattern, so that if you jump straight into a stream of bytes you'll be able to see that you're in the middle of a multi byte char.
- utf16: It splits the 21 bit real unicode code point into two 10 bit halves. Then surrogate markers are added, high and low. If the first 2-byte chunk starts with binary 110110 then you know you should read in the next two bytes and combine them in one code point to get a char that lives above BMP. Or if you start in the middle of a stream and encounter binary 110111 you know you need the previous 2-byte chunk to correctly get the char.

```
High surrogate always:  1101 10xx xxxx xxxx
Low surrogate always:   1101 11xx xxxx xxxx
                                 ^^^^^^^^^^^
                                 these 10 bits from each = the actual character
```

- JSON is utf8 but inside json strings you can write '\uXXXX' which is a utf 16 code unit.
- A Go rune is just and int32. It holds the raw 21 bit Unicode point number in a 32 bit int. Wasteful/overkill for ascii but convenient and doesn't really matter much. When you range over a string in Go it "secretly" parses as runes.

### About lexing and parsing

- A lexer turns a stream of bytes into chunks/tokens/understandable bites baseed on syntax. It's the "dumb part", parsing is the "smart" part that checks if things are in valid order, the bracket is closed etc.
- Rule of thumb: If you eat a token advance it! It's easy to lose track of where you are with next() and peek() etc. the func that reads a token shoudl advance to next token before returning.

### About Go

- io.reader is the raw interface (implements Read([]byte) (int, error)) for reading from anything, file, network, some data stream. you can ask for a certain amount of bytes but you dont know how many you will get with io.reader. bufio.Reader wraps it and adds a buffer, default is 4096 bytes. then you read from that buffer, for example rune by rune, no underlying io.Reader.Read() happens until needed. Makes it easy to read rune by rune with ReadRune()
- interface{} is Go for "a box that can hold anything". Every type automatically satisfies interface{}. It imposes zero requirements. You use type assertion to unwrap it and use the underlying value again. like

```go
obj := node.(Object) // i know this holds an Object, open it, panics if wrong type

// so you should use ok pattern instead
obj, ok := node.(Object)
if !ok {
// node wasn't an object, handle
}
```
