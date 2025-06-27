# Bedrock Snippets

## Running Locally
### Requirements
- Go
- Node & NPM

### Installation
```
npm install
```

### Running
Firstly, move into the website directory:
```
cd website
```
Then there are three ways you can run it locally.

Build:
```
go run .
```

Build and then start a server at [localhost:8080](localhost:8080):
```
go run . -dev
```

Using [air](https://github.com/air-verse/air?tab=readme-ov-file#installation) to automatically rebuild and restart the server on a file change:
```
air
```