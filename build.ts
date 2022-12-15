import Client, { connect } from "@dagger.io/dagger"

 // initialize Dagger client
connect(async (client: Client) => {

	let src = client
		.host()
		.workdir({exclude: [".gitignore", "node_modules", "build.ts"]})

	let gocache = client
		.cacheVolume("gomodcache")

	let goarches = [
		"amd64",
		//"arm64",
	]
	let geese = ["linux"]

	await Promise.all(goarches.map(async (goarch) => {
		await Promise.all(geese.map(async (goos) => {

			let go = client
				.container({ platform: `${goos}/${goarch}` })
				.from("golang:1.19")


			let deps = go
				.withExec(["apt-get", "update"])
				.withExec(["apt-get", "install", "-y", "xorg-dev", "libgl1-mesa-dev", "libopenal1", "libopenal-dev", "libvorbis0a", "libvorbis-dev", "libvorbisfile3"])

			let builder = deps
				.withMountedCache("/cache", gocache)
				.withEnvVariable("GOMODCACHE", "/cache")
				.withMountedDirectory("/app", src)
				.withWorkdir("/app/cmd/greet3d")
				.withEnvVariable("GOOS", goos)
				.withEnvVariable("GOARCH", goarch)
				.withExec(["go", "build", "-o", `greet3d`])

			await builder
				.directory("/app/cmd/greet3d/images")
				.export(`build/${goos}/${goarch}/images`)

			await Promise.all([ "earth.frag", "earth.vert" ].map(async (f) => {
				await builder
					.file(`/app/cmd/greet3d/${f}`)
					.export(`build/${goos}/${goarch}/${f}`)
			}))

			await builder
				.file("greet3d")
				.export(`build/${goos}/${goarch}/greet3d`)
		}))
	}))
}, {LogOutput: process.stdout});
