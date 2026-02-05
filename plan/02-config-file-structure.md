# Step 02: Config & File Structure

## Goal
Implement the config system and automatic `.tasks/` directory initialization.

## Tasks

### 2.1 Create Config Package (`internal/config/config.go`)

```go
type Config struct {
    Board BoardConfig
}

type BoardConfig struct {
    Columns []string
}
```

Functions:
- `Load(path string) (*Config, error)` - Load config from TOML file
- `GetDefaultConfig() *Config` - Returns default config with ["TODO", "DOING", "DONE"]
- `Save(path string, config *Config) error` - Write config to file

### 2.2 Create Initialization Logic

On app startup:
1. Check if `.tasks/` directory exists
2. If not, create it
3. Check if `.tasks/config.toml` exists
4. If not, create with default columns
5. Print message: "Created default config in .tasks/config.toml"

### 2.3 Default config.toml Content
```toml
[board]
columns = ["TODO", "DOING", "DONE"]
```

## Testing This Step

```bash
# Remove .tasks if exists
rm -rf .tasks

# Run app
./doings

# Check files created
ls -la .tasks/
cat .tasks/config.toml
```

**Expected output**:
- Message about config creation
- `.tasks/` directory exists
- `config.toml` contains default columns

## Files to Create/Modify
- `internal/config/config.go` - Config handling
- `main.go` - Add initialization call

## Definition of Done
- [ ] Running app in fresh directory creates `.tasks/`
- [ ] `config.toml` is created with default columns
- [ ] User sees message about config creation
- [ ] Running again doesn't recreate/overwrite existing config
