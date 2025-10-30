workflow "Platform Testing Matrix" {
    on = "pull_request"
    
    // Test across multiple operating systems
    foreach os in ["ubuntu-latest", "windows-latest", "macos-latest"] {
        job "test_${os}" {
            runs-on = "${os}"
            step "Checkout" uses = "actions/checkout@v3"
            step "Setup Go" uses = "actions/setup-go@v3" 
            step "Run tests" run = "go test ./..."
            step "Build for ${os}" run = "go build -o bin/app-${os} ."
        }
    }
    
    // Test multiple Go versions
    for version in 18..21 {
        job "go_${version}_test" {
            runs-on = "ubuntu-latest"
            step "Setup Go ${version}" uses = "actions/setup-go@v3"
            step "Test with Go 1.${version}" run = "go test -v ./..."
        }
    }
    
    // Performance testing with different loads
    repeat 3 {
        job "performance_test" {
            runs-on = "ubuntu-latest"
            step "Load test iteration" run = "ab -n 1000 -c 10 http://localhost:8080/"
            step "Collect metrics" run = "echo 'Performance test completed'"
        }
    }
}
