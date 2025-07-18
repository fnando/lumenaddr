# lumenaddr

Generate vanity addresses for stellar (addresses with specific suffixes).

## Downloads

I have already built
[executables available for download](https://github.com/fnando/lumenaddr/releases).

## Usage

Just run the command providing the words you're looking for. Just remember that
longer the word, the more time it'll take. The following command will lookup for
keys that match the words `STELLAR` and `LUMENS`, and output them to the
console.

```
lumenaddr STELLAR LUMENS
```

To generate mnemonic (recovery phrase), you can use `-mnemonic` but this is
considerably slower.

```
lumenaddr -mnemonic STELLAR LUMENS
```

## Screenshots

![](https://github.com/fnando/lumenaddr/raw/main/screenshots/save.png)

![](https://github.com/fnando/lumenaddr/raw/main/screenshots/print-keys.png)

![](https://github.com/fnando/lumenaddr/raw/main/screenshots/output.png)

![](https://github.com/fnando/lumenaddr/raw/main/screenshots/dataclips.png)

## License

(The MIT License)

Copyright (c) 2018 Nando Vieira

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the 'Software'), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
