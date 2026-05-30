Go-Core
================================

> [Aeon Digital](http://www.aeondigital.com.br)  
> rianna@aeondigital.com.br

&nbsp;

A modular repository of general-purpose packages and utilities for Go applications. 
Managed as a monorepo, each directory represents an independent module with its own versioning, allowing you to import only what your application needs.


&nbsp;
&nbsp;


________________________________________________________________________________

## Repository Structure

```text
Go-Core/
├── .github/workflows/    # CI/CD automation pipelines
├── tools/                # Core debugging and formatting utilities
└── README.md             # Global documentation
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Modules

### Tools

A collection of standalone utilities for debugging, dynamic type allocation, and cross-language format translations.

* **Import path:** `://github.com`
* **Features:** JSON debugging dump, advanced log and error formatting, dynamic generic instance allocation, and universal datetime layout translation.


&nbsp;
&nbsp;


________________________________________________________________________________

## Install

Open your shell and use `go get` pointing directly to the desired module path:

```shell
  go get ://github.com@latest
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Test Suite

Each module contains its own isolated test suite. To scan and run tests across all modules locally, use:

```shell
  for dir in \((find . -name "go.mod" -exec dirname {} \;); do cd "\)dir" && go test -v -cover ./... && cd - > /dev/null; done
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Licence

This project is offered under the [MIT license](LICENCE.md).
