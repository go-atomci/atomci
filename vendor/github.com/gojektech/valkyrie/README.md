# Valkyrie

<p align="center">
  <img src="doc/valkyrie-logo.png" width="360" alt="Valkyrie Logo">
</p>

<p align="center">
  <a href="https://travis-ci.org/gojektech/valkyrie">
	<img src="https://travis-ci.org/gojektech/valkyrie.svg?branch=master" alt="Build Status" />
  </a>
</p>

## Description
Valkyrie is a utility that helps you aggregate multiple errors in Go, while maintaining thread safety.

## Installation
```
go get -u github.com/gojektech/valkyrie
```

## Usage

Consider the case of an error prone operation, where there is a possibility of encountering multiple, mutually independant errors while running an operation. We can use Valkyrie to collate all the errors into a single error which we can then return:

```go
func errorProneOperation(n int) error {
	// Create a new Multierror instance, which implements the error interface
	errs := new(valkyrie.MultiError)

	for i := 0; i < n; i++ {
		errs.Push(fmt.Sprintf("error in iteration %d", i))
	}
	// When you have to return an error, call the `HasError` method
	// which returns nil if the length of errors is 0, and returns the errs instance itself if its not
	return errs.HasError()
}
```

Now, the `errorProneOperation` can be used just like a function that returns a single error:

```go
func main() {
	err := errorProneOperation(3)
	if err != nil {
		fmt.Println("err :", err)
	}
	// Outputs:
	// err : error in iteration 0, error in iteration 1, error in iteration 2
	err = errorProneOperation(0)
	if err != nil {
		fmt.Println("err :", err)
	}
	// Does not output anything
}
```

For more documentation, you can visit [godoc.org](https://www.godoc.org/github.com/gojektech/valkyrie).

## License

```
Copyright 2018, GO-JEK Tech (http://gojek.tech)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
