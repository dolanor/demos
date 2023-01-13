# [<img height="50px" src="https://dagger.io/logo.svg">](https://dagger.io) Demos

This holds the demos I present during our [Community Calls](https://dagger.io/events).

## Content

Currently, the project is a simple Go library for greeting. It is located in the [greeting](./greeting) package.

It is used by 2 apps:
- a [CLI](./cmd/greet/main.go)
- a [3D frontend](./cmd/greet3d/main.go)

### Dagger Pipeline using the [Dagger Node SDK](https://docs.dagger.io/sdk/nodejs)

### CLI

The CLI is straightforward, and thanks to Go, we can cross-compile in our pipeline just changing the 
classic Go environment variables `GOOS` and `GOARCH`. A simple [`for each`](https://github.com/dolanor/demos/blob/42331888d99c0e890d1b75f4f333d1629fbccfa3/ci/greet.ts#L16-L17) allows to cross compile.

N.B.: the `for each` is a little more involved than the JS `forEach` as we have `async` functions to call.
N.B.2: I'm not a JS/TypeScript expert :p

### Greet 3D

The app shows a rotating Earth in space. The api uses the `greeting` package to generate a greeting 2D Texture that is then applied on the 3D model of the Earth.

The simple cross-compiling of Go cannot operate here. The [g3n 3D game engine](http://g3n.rocks) access some native libraries (mostly the 3D hardware drivers). Therefore, we need `CGO`. And we need the native libraries to be present on the system when compiling.
We use Dagger capacity to use multi arch images, and pull the right image for the `GOARCH` we're building for.
The `apt-get` will get the right libraries (`amd64` or `arm64`), and then it will build and link correctly for the aimed platform.

The other tricky part is that a 3D app needs some assets (images/textures) or external scripts (shader files). We use Dagger to copy them alongside the binary so we can execute the binary and it can find its assets.

N.B.: in Go, I could actually use `go:embed` to embed those assets in my binary, but for other 3D projects in C++, you don't have this possibility natively and you need some strategy for this. So I decided to show how to do it if you were in this case.
