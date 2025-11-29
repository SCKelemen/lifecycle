Output lifecycle.gif
Set Width 1200
Set Height 800
Set Theme "Monokai Pro"
Set FontSize 18
Set Shell "bash"

# Set TERM for color support
Type "export TERM=xterm-256color"
Enter
Sleep 200ms

# Build the example
Type "cd example && go build -o main main.go"
Enter
Sleep 2s

# Run the lifecycle example
Type "./main"
Enter
Sleep 3s

