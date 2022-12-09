import Client, { connect } from "@dagger.io/dagger"

 // initialize Dagger client
connect(async (client: Client) => {

	let srcId = await client
		.host()
		.workdir([".gitignore", "node_modules", "build.ts"])
		.id()

	let gocache = await client
		.cacheVolume("gomodcache")
		.id()

	let goarches = [
		"amd64",
		//"arm64",
	]
	let geese = ["linux"]

	await Promise.all(goarches.map(async (goarch) => {
		await Promise.all(geese.map(async (goos) => {

			let go = client
				.container(undefined, `${goos}/${goarch}`)
				.from("golang:1.19")


			let deps = go
				.exec(["apt-get", "update"])
				.exec(["apt-get", "install", "-y", "xorg-dev", "libgl1-mesa-dev", "libopenal1", "libopenal-dev", "libvorbis0a", "libvorbis-dev", "libvorbisfile3"])

			let builder = deps
				.withMountedCache("/cache", gocache.id)
				.withEnvVariable("GOMODCACHE", "/cache")
				.withMountedDirectory("/app", srcId.id)
				.withWorkdir("/app/cmd/greet3d")
				.withEnvVariable("GOOS", goos)
				.withEnvVariable("GOARCH", goarch)
				.exec(["go", "build", "-o", `greet3d`])

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
