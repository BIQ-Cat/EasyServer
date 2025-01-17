# EasyServer

## A light multy-lingual framework written in Go

### Features
- Multi-lingual structure, so you can use many languages in one project with ease
- Customisable system: EasyServer can change it's structure for (and _with_) your project
- Source-based configuration, which can make your server faster, is JSON-styled, easy-to-learn and comes as a __default__ variant (it means you can _always_ replace it with a JSON variant)
- FastCGI is supported. (coming soon)

## Installation
### From source
First of all, you should build EasyMaker or download pre-built variant from `Releases`. There are scripts `build_windows.bat` and `build_linux.sh` for it.
Then, you can run
```console
./easymaker install -target=all
```
for building all parts of EasyServer. Use _help_ target to see other options.

To build sources without installation, you can use `build` command with the same targets.

In the `build` folder, you can see your compiled EasyServer and (or) libraries.

## Documentation
Visit [pkg.go.dev](https://pkg.go.dev/github.com/BIQ-Cat/easyserver)




