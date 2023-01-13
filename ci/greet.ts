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

	let goarches = ["amd64", "arm64"]
	let geese = ["windows", "darwin", "linux"]

	await Promise.all(goarches.map(async (goarch) => {
		await Promise.all(geese.map(async (goos) => {
			await go
				.withMountedDirectory("/app", srcId.id)
				.withWorkdir("/app")
				.withEnvVariable("GOOS", goos)
				.withEnvVariable("GOARCH", goarch)
				.exec(["go", "build", "-o", `greetings`])
				.file("greetings")
				.export(`build/${goos}/${goarch}/dagger-greetings`)
		}))
	}))
}, {LogOutput: process.stdout});
