const { execSync } = require("child_process");
const os = require("os");
const path = require("path");

const platform = os.platform();
const arch = os.arch();

let binaryUrl;
if (platform === "darwin" && arch === "x64") {
  binaryUrl = "https://github.com/chidinma21/env-lint/releases/latest/download/env-lint-darwin-amd64";
} else if (platform === "darwin" && arch === "arm64") {
  binaryUrl = "https://github.com/chidinma21/env-lint/releases/latest/download/env-lint-darwin-arm64";
} else if (platform === "linux" && arch === "x64") {
  binaryUrl = "https://github.com/chidinma21/env-lint/releases/latest/download/env-lint-linux-amd64";
} else if (platform === "win32" && arch === "x64") {
  binaryUrl = "https://github.com/chidinma21/env-lint/releases/latest/download/env-lint-windows-amd64.exe";
} else {
  console.error(`Unsupported platform: ${platform} ${arch}`);
  process.exit(1);
}

const outputPath = path.join(__dirname, "..", "env-lint" + (platform === "win32" ? ".exe" : ""));
execSync(`curl -L ${binaryUrl} -o ${outputPath}`);
execSync(`chmod +x ${outputPath}`);
