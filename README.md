# go-gellato-membership

Saintropic membership RESTful service. Built using Go language with go-siris framework (The fastest web framework for Golang | fasthttp) and dynamodb as database. Please support [go-siris framework](https://github.com/go-siris/siris) (The fastest web framework for Golang!) and read more about it to properly use and develop this project, as [go-siris](https://github.com/go-siris/siris) was the core of this project.

List of dependency used in this project: 
* [go-siris](https://github.com/go-siris/siris)
* [go-qrcode](https://github.com/skip2/go-qrcode)
* [jwt-go](https://github.com/dgrijalva/jwt-go)
* [sendgrid-go](https://github.com/sendgrid/sendgrid-go)
* [aws-sdk-go dynamodb as database](https://github.com/aws/aws-sdk-go)
* [go.uuid](https://github.com/satori/go.uuid)
* [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)

# Table of contents

* [Installation](#installation)
* [License](#license)

# Installation

[Go Programming Language](https://golang.org/dl/), at least version 1.8

```sh
$ go get 
```

> Run the command to get all package dependencies.

Create new file env.json and copy all of example-env.json file to newly created env.json file. 
Change the value of each element in env.json according to your specification.


# License

This project was distributed under the MIT License found in [LICENSE file](LICENSE).