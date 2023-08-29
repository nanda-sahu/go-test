# Chapter 1: Basic REST Server

The aim of this chapter is to get you set up and ensure you can build and run a Go application, and therefore is more a set of instructions rather than a coding exercise as most of the subsequent chapters.

## Module Setup

The first part is to create a new Go module (and an accompanying git repo). 

Let's create the directory for the module:
```
mkdir ecommerce-workshop
cd ecommerce-workshop
```

Now let's create a new repository on HPE's github under your user with the name 'ecommerce-workshop'. Once the repo exists, we can init git and link it to the repository. This will allow us to use branches for each chapter to isolate the work, and have code reviews for pull requests.
```
git init
git remote add origin git@github.hpe.com:{your-gh-user}/ecommerce-workshop.git
```

Let's init our Go module using the `go mod` command and then make our first git commit. This will result in us having a first commit with a `go.mod` file. This file defines basic info about the module such as the minimum version of Go required and the moudle name. Later this file will also contain a list of all dependencies required for the project (along with their versions).
```
go mod init github.hpe.com/{your-gh-user}/ecommerce-workshop
git add .
git commit -m "Initial commit"
```

## HTTP Server

Like other C type languages, the entry point for a Go program is through a 'main' function. By convention, this file goes in a `main.go` file that lives in the `cmd/{app-name}` directory (for a summary of the repo layout we follow, check out this [example](https://github.com/golang-standards/project-layout). 

Let's go ahead and make our cmd directory and copy across our sample application. We can also copy across the sample Makefile and .gitignore to make our life easier.
```
git checkout -b initial-rest-server
mkdir -p cmd/rest-server
cp /path/to/go-chapter/new-hire-workshop/section-1_rest/chapter-1_basic-rest-server/materials/main.go ./cmd/rest-server/
cp /path/to/go-chapter/new-hire-workshop/section-1_rest/chapter-1_basic-rest-server/materials/Makefile .
cp /path/to/go-chapter/new-hire-workshop/section-1_rest/chapter-1_basic-rest-server/materials/sample-gitignore ./.gitignore
git add .
git commit -m "Copied across sample REST server"
```

Now let's run our server and test that it works. Taking a look at the code, you'll see a very simple web server serving on port 8080. This is done through the [net/http](https://pkg.go.dev/net/http) library that is included in Golang’s standard library. The standard library web server is production ready and commonly used in our production applications.

For the future exercises we can use the Makefile commands, but for the first time we'll do it using the Go tool directly to build up familiarity.
```
mkdir ./dist                  # make a dir that is ignored by git to keep our binaries
go build -o dist ./cmd/...    # tell Go to recursively build all commands, and put then in dist
```

## Run and Test

Finally, let's run our server
```
./dist/rest-server
```
And test it from another terminal:
```
curl http://localhost:8080/test
Hello on path "/hello"
```

If you are able to run the server and get a response from your curl requests then this section is complete. The only task that remains is to merge you branch into the main branch in preparation for the workshop.
```
git checkout main
git merge initial-rest-server
```

## TODO: Creating and Syncing a GitHub Repo
