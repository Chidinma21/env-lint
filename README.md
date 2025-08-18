# env-lint

> A fast and flexible CLI tool to validate your `.env` files against a JSON or YAML schema.

Tame your environment configs. Catch missing or invalid `.env` variables *before* they break your app.

---

## ✨ Features

- ✅ Schema validation for `.env` files (JSON and YAML)
- ♺ JSON/YAML schema generation from an existing `.env` file
- 🔢 Type checking for `string`, `number`, and `boolean`
- ⚠️ Support for optional keys and default values
- 🎟️ Support for allowed values (enums), patterns (RegEx), min/max, minLength/maxLength, length
- 🟡 Support for warning suppression - keeping output clean in CI
- ❕ Optional fail-fast mode - stop on first error
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

#### Supported Schema Rules
| Rule          | Type           | Description                                                         |
| ------------- | -------------- | ------------------------------------------------------------------- |
| `type`        | string         | The variable type: `string`, `number`, or `boolean`. **(required)** |
| `required`    | bool           | If true, the key must exist in `.env`.                              |
| `default`     | any            | Value to use if key is missing.                                     |
| `allowed`     | []any          | List of valid values (enum).                                        |
| `pattern`     | string (RegEx) | Value must match this regex pattern.                                |
| `length`      | int            | Exact length of the string.                                         |
| `minLength`   | int            | Minimum string length.                                              |
| `maxLength`   | int            | Maximum string length.                                              |
| `min`         | float          | Minimum numeric value.                                              |
| `max`         | float          | Maximum numeric value.                                              |
| `customError` | string         | Custom error message when validation fails.                         |


#### Example .env

```bash
PORT=3000
APP_NAME=env-lint
DEBUG_MODE=true
MAX_USERS=50
API_KEY=abc123
```

#### Example schema.json
```bash
{
  "PORT": {
    "type": "number",
    "required": true,
    "min": 1000,
    "max": 9999
  },
  "DEBUG_MODE": {
    "type": "boolean",
    "required": false,
    "default": false
  },
  "APP_NAME": {
    "type": "string",
    "required": true,
    "minLength": 3,
    "maxLength": 20
  },
  "MAX_USERS": {
    "type": "number",
    "required": false,
    "default": 10,
    "min": 1,
    "max": 100
  },
  "API_KEY": {
    "type": "string",
    "required": true,
    "pattern": "^[a-zA-Z0-9]+$",
    "length": 6,
    "customError": "API_KEY must be alphanumeric and exactly 6 characters long"
  },
  "ENV": {
    "type": "string",
    "allowed": ["development", "staging", "production"],
    "default": "development"
  }
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

❌ APP_NAME                Too short (minLength = 3)
❌ API_KEY                 API_KEY must be alphanumeric and exactly 6 characters long
❌ ENV                     Invalid value: "local" (allowed: development, staging, production)
🟡 DEBUG_MODE              Optional key missing (ok)
🟡 MAX_USERS               Using default value: 10

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

### 🔍 Validate Command
```bash
./env-lint validate [flags]
```

#### Available Flags:

- `-e, --env` `string`: 
Path to the `.env` file (default: `.env`)

- `-s, --schema` `string`: 
Path to the schema file (default: schema.json)

- `-w, --suppress-warnings` `boolean`: 
Suppress non-critical warnings in output

- `-t, --strict-mode` `boolean`: 
Fail if `.env` contains keys not defined in schema

- `-f, --fail-fast` `boolean`: 
Stop validation after the first error

### 🔍 Generate Schema
```bash
./env-lint generate-schema [flags]
```

#### Available Flags:

- `-e, --env` `string`: 
Path to the `.env` file to generate schema from (default: `.env`)

- `-f, --format` `string`: 
Output format: `json`, `yaml`, or `yml` (default: `json`)

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