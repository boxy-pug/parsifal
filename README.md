# parsifal

A hand-written JSON parser in Go. This is a learning project inspired by [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball.

I'm starting with JSON, and then i'll add YAML and TOML maybe, and a CLI interface maybe.

Built from scratch to understand how parsers and lexers work under the hood. AI is used for concept explanation and test writing, but source code is written by hand.

## Things I've learned:

- A lexer turns a stream of bytes into chunks/tokens/understandable bites baseed on syntax. It's the "dumb part", parsing is the "smart" part that checks if things are in valid order, the bracket is closed etc.
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
- json is a data serialization format thats easy to parse for computers and human readable. All keys are strings. All strings are in double quotes. objects{"key": val}, arrays [val, val], string, number, bool, null basically. doesnt care about whitespace.
- io.reader is the raw interface for reading, you dont know how many bytes it will give you. bufio.Reader wraps it and adds a buffer. Makes it easy to read rune by rune with ReadRune()
