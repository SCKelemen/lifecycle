# Troubleshooting Color Output

If you're not seeing colors in the output, here are some things to check:

## Terminal Support

The colors require a terminal that supports ANSI color codes. Most modern terminals do, but some environments don't:

1. **Check your terminal type:**
   ```bash
   echo $TERM
   ```
   If it shows `dumb`, colors won't work. Try setting:
   ```bash
   export TERM=xterm-256color
   go run ./example
   ```

2. **Check if colors are supported:**
   ```bash
   tput colors
   ```
   Should return a number > 0 (like 256). If it returns -1, colors aren't supported.

3. **Check for NO_COLOR:**
   ```bash
   echo $NO_COLOR
   ```
   If this is set, colors are disabled. Unset it:
   ```bash
   unset NO_COLOR
   ```

## Testing Colors Directly

Test if your terminal supports colors:

```bash
# Test ANSI colors
echo -e "\033[32mThis should be green\033[0m"

# Test with lipgloss
cd /Users/samuel.kelemen/Code/github.com/SCKelemen/lifecycle
go run -c 'package main; import ("fmt"; "github.com/charmbracelet/lipgloss"); func main() { s := lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")); fmt.Println(s.Render("This should be blue")) }'
```

## Running in Different Terminals

- **iTerm2 / Terminal.app (macOS)**: Should work by default
- **VS Code Integrated Terminal**: Should work, but check terminal type
- **SSH sessions**: May need `TERM=xterm-256color`
- **CI/CD pipelines**: Usually don't support colors (this is expected)

## Force Colors (for testing)

If your terminal supports colors but they're not showing, you can test by setting:

```bash
export TERM=xterm-256color
export CLICOLOR=1
export CLICOLOR_FORCE=1
go run ./example
```

## Expected Behavior

- **With color support**: You'll see colored service names, API names, event types, and status codes
- **Without color support**: You'll see plain text (this is normal in some environments)

The code will work in both cases - colors are just a visual enhancement.

