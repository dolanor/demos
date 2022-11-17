import Client, { connect } from "@dagger.io/dagger"

 // initialize Dagger client
connect(async (client: Client) => {
	let go = await client
		.container()
		.from("golang:1.19")

	let srcId = await client
		.host()
		.workdir([".gitignore", "node_modules", "build.ts"])
		.id()

	let gocache = await client
		.cacheVolume("gomodcache")
		.id()

	let builder = await go
		.exec(["apt-get", "update"])
		.exec(["apt-get", "install", "-y", "xorg-dev", "libgl1-mesa-dev", "libopenal1", "libopenal-dev", "libvorbis0a", "libvorbis-dev", "libvorbisfile3"])

	// let goarches = ["amd64", "arm64"]
	// let geese = ["windows", "darwin", "linux"]
	let goarches = ["amd64"]
	let geese = ["linux"]

	await Promise.all(goarches.map(async (goarch) => {
		await Promise.all(geese.map(async (goos) => {
			let lbuilder = await builder
				.withMountedCache("/cache", gocache.id)
				.withEnvVariable("GOMODCACHE", "/cache")
				.withMountedDirectory("/app", srcId.id)
				.withWorkdir("/app/cmd/greet3d")
				.withEnvVariable("GOOS", goos)
				.withEnvVariable("GOARCH", goarch)
				.exec(["go", "build", "-o", `greet3d`])

			await lbuilder
				.directory("/app/cmd/greet3d/images")
				.export(`build/${goos}/${goarch}/images`)

			await Promise.all([ "earth.frag", "earth.vert" ].map(async (f) => {
				await lbuilder
					.file(`/app/cmd/greet3d/${f}`)
					.export(`build/${goos}/${goarch}/${f}`)
			}))



			await lbuilder
				.file("greet3d")
				.export(`build/${goos}/${goarch}/dagger-greetings`)
		}))
	}))
});
