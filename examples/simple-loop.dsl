workflow "Simple Loop Test" {
    on = "push"
    
    for i in 1..3 {
        job "test_${i}" {
            runs-on = "ubuntu-latest"
            step "Setup ${i}" uses = "actions/setup-node@v3"
        }
    }
}
