# adguard-home demo — rendered with VHS (https://github.com/charmbracelet/vhs).
# No account, signup, keyring or external network: only invented local API data.
# Regenerate from the repository root:
#   go build -o bin/adguard-home ./cmd/adguard-home && go build -o bin/demoserver ./tools/demoserver && vhs .demo/demo.tape
# The shell exits on a failed command rather than recording an error as a result.
# This CLI has no --jq flag; jq filters its actual JSON output instead.

Output assets/demo.gif

Set Shell bash
Set FontSize 15
Set Width 1200
Set Height 760
Set Padding 18
Set Theme "Dracula"
Set TypingSpeed 28ms

Hide
Type "set -e -o pipefail"
Enter
Type "export PATH=$PWD/bin:$PATH XDG_CONFIG_HOME=$(mktemp -d)"
Enter
Type "export ADGUARD_URL=http://127.0.0.1:8644 ADGUARD_USERNAME=demo ADGUARD_PASSWORD=demo-password"
Enter
Type `demoserver 127.0.0.1:8644 >"$XDG_CONFIG_HOME/demoserver.log" 2>&1 & demo_server_pid=$!`
Enter
Type `trap 'kill "$demo_server_pid" 2>/dev/null || true; rm -rf "$XDG_CONFIG_HOME"' EXIT`
Enter
Type `sleep 1; kill -0 "$demo_server_pid"`
Enter
Type `alias adguard-home='env -i PATH="$PATH" XDG_CONFIG_HOME="$XDG_CONFIG_HOME" ADGUARD_URL="$ADGUARD_URL" ADGUARD_USERNAME="$ADGUARD_USERNAME" ADGUARD_PASSWORD="$ADGUARD_PASSWORD" adguard-home'`
Enter
Type "clear"
Enter
Show

Sleep 500ms

Type "adguard-home status"
Sleep 250ms
Enter
Sleep 2400ms

Type "adguard-home rewrites list"
Sleep 250ms
Enter
Sleep 2400ms

Type "clear"
Enter
Sleep 300ms

Type "adguard-home stats -o json | jq '{queries: .num_dns_queries, blocked: .num_blocked_filtering}'"
Sleep 250ms
Enter
Sleep 2400ms

Type "adguard-home rewrites add lab.example.test 192.0.2.30"
Sleep 250ms
Enter
Sleep 2400ms

Type "adguard-home rewrites list"
Sleep 250ms
Enter
Sleep 2400ms
