import Client, { connect } from "@dagger.io/dagger"
import fs from "node:fs"

 // initialize Dagger client
connect(async (client: Client) => {
	await client
		.container()
		.from("node")
		.withNewFile("/tmp/filedir/myTextFile", { contents: "my text content" })
		.withExec([ "cat", "/tmp/filedir/myTextFile" ])
		.stdout()
}, {LogOutput: process.stdout});
