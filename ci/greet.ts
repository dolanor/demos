import Client, { connect } from "@dagger.io/dagger"

 // initialize Dagger client
connect(async (client: Client) => {
	const go = await client
		.container()
		.from("golang:1.20.2-alpine3.17")

	const src = await client
		.host()
		.directory(".", { exclude: [".gitignore", "node_modules", "build.ts"] })

	const goarches = ["amd64", "arm64"]
	const geese = ["windows", "darwin", "linux"]

	await Promise.all(goarches.map(async (goarch) => {
		await Promise.all(geese.map(async (goos) => {
			await go
				.withMountedDirectory("/app", src)
				.withWorkdir("/app/build")
				.withEnvVariable("GOOS", goos)
				.withEnvVariable("GOARCH", goarch)
				.withExec(["go", "build", "../cmd/greet"])
				.directory(".")
				.export(`build/${goos}/${goarch}/`)
		}))
	}))
}, {LogOutput: process.stdout});
