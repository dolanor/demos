import Client, { connect } from "@dagger.io/dagger"
import fs from "node:fs"

 // initialize Dagger client
connect(async (client: Client) => {
	let textFile = fs.openSync("myTextFile", "w")
	fs.writeFileSync(textFile, "my text content\n")

	let packageFile = await client
		.host()
		.workdir(undefined, ["myTextFile"])
		.id()
	
	await client
		.container()
		.from("node")
		.withMountedDirectory("/tmp/filedir", packageFile.id)
		.exec([ "cat", "/tmp/filedir/myTextFile" ])
		.stdout()
		.contents()
}, {LogOutput: process.stdout});
