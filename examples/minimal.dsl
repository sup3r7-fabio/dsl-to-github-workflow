workflow "Minimal Test" {
    on = "push"
    job "simple" {
        runs-on = "ubuntu-latest"
        step "test" run = "echo hello"
    }
}
