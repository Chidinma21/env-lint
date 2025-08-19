#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");

const binaryPath = path.join(__dirname, "..", "env-lint");
const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });

child.on("exit", code => process.exit(code));
