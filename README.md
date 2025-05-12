# oneHome CLI (`oh`)

A command‐line interface for interacting with the [oneHome API](https://onehome.dogado.de/docs/api/).

---

## 🔧 Installation

Download the `oh` binary from the [releases](https://github.com/group-one-platform/oh/releases) page.

You can also install it with Go (requires Go 1.23+):

```bash
go install github.com/group-one-platform/oh@latest
```

Or clone this repository and build manually:

```bash
git clone https://github.com/group-one-platform/oh.git
cd oh
make build
```

The resulting `oh` executable will be in your `$GOPATH/bin` (or `./oh` if you used `go build`).

You can also build for all platforms (Windows, Linux, Mac):

```bash
make build-all
```

---

### 📑 Shell Completion

`oh` can generate shell‐completion scripts for Bash, Zsh, Fish and PowerShell. Adapt for your distribution:

```bash
# Bash
oh completion bash > /etc/bash_completion.d/oh

# Zsh
oh completion zsh > "${fpath[1]}/_oh"

# Fish
oh completion fish > ~/.config/fish/completions/oh.fish

# PowerShell
oh completion powershell > ~/oh.ps1
```

Or, for a one‐off session, you can eval the output directly:

```bash
# In Bash
eval "$(oh completion bash)"
# Or Zsh
eval "$(oh completion zsh)"
```
---

After installing into your shell’s completion directory (and reloading/restarting your shell), typing:

```bash
oh <TAB>
```

will auto‐complete commands, flags, and even argument values (e.g. your vps execute actions and even server ids, flavours, images etc).

## 🔐 Authentication

Before making any API calls, you need to provide your API token.  The token is stored in `~/.oh.yaml`.

To enter your token interactively:

```bash
oh token
```

Or pass the token directly on the command line:

```bash
oh token YOUR_API_TOKEN_HERE
```

You should see:

```
🔐 Token stored!
```

### Alternative Token Input Methods

In addition to entering your token interactively or pass it as an argument, you can also supply it from a file or standard input.

#### From a File

If you have a file containing just your token (for example `token.txt`), you can use:

```bash
oh token --file ./token.txt
```

This will read the first line of the file as your API token.

#### From STDIN

You can pipe the token in directly from another command or environment variable:

```bash
cat ./token.txt | oh token --file -
```

or

```bash
echo $ONEHOME_TOKEN | oh token --file -
```

Here `-` tells `oh token` to read the token from standard input.

---

## ⚙️ Configuration

By default, `oh` looks for a config file at:

- `$HOME/.oh.yaml`

And reads environment variables matching its keys.  The only default is:

```yaml
base_url: "https://onehome.dogado.de/api/v1/"
```

If you have access to the API in other environments, feel free to update the `base_url` accordingly.

You can override this in your config file or via the `--config` flag:

```bash
oh --config /path/to/your-config.yaml vps list
```

---

## 💡 Global Flags

These flags are available on *all* commands:

- `--json`            Print JSON (Go‐encoded)
- `-j`, `--jq`        Pipe JSON through `jq` (optional filter)
- `--config <file>`   Path to config file (default `$HOME/.oh.yaml`)

Example:

```bash
# List VPSes as raw JSON:
oh vps list --json

# List and then filter with jq:
oh vps list --jq '.[] | select(.Status=="active")'

# Output with jq, no filter
oh vps list -j

# Show just the first server in a key: value form:
oh vps get 42
```

---

## 📚 Commands

Use `oh <command> --help` to see full usage and subcommands.  Below are some common workflows.

### VPS Management (`oh vps`)

- **List all VPSes**

  ```bash
  oh vps list
  ```

- **Get a single VPS** (form view):

  ```bash
  oh vps get 42
  ```

- **Execute an action** (e.g. soft-reboot, power-off):

  ```bash
  oh vps execute 42 soft-reboot
  ```

- **Manage flavours**
  - List available flavours for server 42:
    ```bash
    oh vps flavour list 42
    ```
  - Change flavour:
    ```bash
    oh vps flavour set 42 --flavour=33
    ```

- **Manage networks**
  - List *all* available networks:
    ```bash
    oh vps network list-available
    ```
  - List networks attached to server 42:
    ```bash
    oh vps network list 42
    ```
  - Attach a network:
    ```bash
    oh vps network attach 42 --network-id=abc123 --ipv4=192.0.2.5
    ```
  - Detach a network:
    ```bash
    oh vps network detach 42 --network-id=abc123
    ```

- **Manage images**
  - List images:
    ```bash
    oh vps image list
    ```
  - Get image details:
    ```bash
    oh vps image get 100
    ```

- **Manage products**
  - List products (with plan JSON):
    ```bash
    oh vps product list
    ```

- **Order a new VPS**

  You can pass the order payload as a JSON file or raw string:

  ```bash
  oh vps order -f order.json
  # or:
  cat order.json | oh vps order -f -
  # or:
  oh vps order '{"name":"foo","product_id":12,"flavour_id":3}'
  ```

  See `oh vps order -h` for more information about the order payload.
---

## 🛠️ Contributing

Feel free to open issues or pull requests—happy to accept improvements, bug fixes, and new commands!

---

## 📄 License

This project is released under the Apache 2.0 License.  See [LICENSE](LICENSE) for details.
