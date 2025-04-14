# Create and open ZED spec

## Usage in zed editor

1. Create task inside `~/.config/zed/tasks.json`
2. Put some task like this:
```json
{
  "label": "create-open-test",
  "command": "~/path-to-zed_create_rspec/zed_create_rspec",
  "hide": "always",
  "reveal": "never",
  "args": ["\"$ZED_RELATIVE_FILE\""],
  "tags": ["ruby-test-create"]
}
```
3. Add to keymap file  `~/.config/zed/keymap.json`
```json
{
  "alt-cmd-.": ["task::Spawn", { "task_name": "create-open-test" }]
}
```


## Development:

### How to run ruby version manuall:
```shell
./app.rb app/controllers/some_controller.rb
./app.rb spec/controllers/some_controller_spec.rb
```

### How to build and run go version:
```shell
go build -o zed_create_rspec zed_create_rspec.go
./zed_create_rspec app/controllers/some_controller.rb
./zed_create_rspec spec/controllers/some_controller_spec.rb
```
