# env-lint

> A fast and flexible CLI tool to validate your `.env` files against a JSON schema.

Tame your environment configs. Catch missing or invalid `.env` variables *before* they break your app.

---

## ✨ Features

- ✅ Validate required environment variables
- 🔢 Type checking for `string`, `number`, and `boolean`
- ⚠️ Support for optional keys and default values
- 🎟️ Support for allowed values (enum)
- 🎨 Color-coded terminal output for easy debugging
- ⚙️ Works great in development and CI/CD pipelines

---

## 🚀 Getting Started

### 1. Install

```bash
git clone https://github.com/chidinma21/env-lint.git
cd env-lint
go build -o env-lint
```

### 2. Prepare Your Files
You’ll need:

- A .env file with your environment variables

- A schema.json file describing expected keys, types, and whether they’re required

#### Example .env

```bash
PORT=3000
APP_NAME=env-lint
```

#### Example schema.json
```bash
{
  "PORT":       { "type": "number",  "required": true },
  "DEBUG_MODE": { "type": "boolean", "required": false },
  "APP_NAME":   { "type": "string",  "required": true }
}
```

### 3. Run the Validator
```bash
./env-lint validate --env .env --schema schema.json
```

Or using shorthand flags:

```bash
./env-lint validate -e .env -s schema.json
```

#### 📦 Sample Output
```bash
🚀 .env file loaded successfully
🚀 Schema file loaded successfully

🔍 Validating environment variables...

❌ APP_NAME                Missing required key  
🟡 DEBUG_MODE              Optional key missing (ok)

━━━━━━━━━━━━━━━━━━━━━━━  
❌ Validation failed. Please fix the errors above.
```

```bash
🚀 .env file loaded successfully
🚀 Schema file loaded successfully

🔍 Validating environment variables...

━━━━━━━━━━━━━━━━━━━━━━━  
✅ All checks passed. Your .env config looks great!
```

## 🔭 Roadmap
- Support default values for optional keys

- Enum / allowed-value validation

- Multiple .env file support (e.g., .env.production)

- Fail-fast mode (--fail-fast)

- JSON output for CI integrations

- npm package wrapper

## 🤝 Contributing
Contributions, issues, and feature requests are welcome!
Please:

- Fork the repository

- Create a branch: git checkout -b feature/your-feature

- Make your changes and add tests

- Commit: git commit -m "Add some feature"

- Push: git push origin feature/your-feature

- Open a Pull Request

- Ensure code follows Go formatting (gofmt) and include tests where applicable.

## 📄 License
MIT © 2025 Chidinma Onuora

Made with ❤️ in Go.