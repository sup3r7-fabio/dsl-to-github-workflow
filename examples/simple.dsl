workflow "CI Pipeline" {
    on = "push"
    job "build" {
        runs-on = "ubuntu-latest"
        step "Checkout" uses = "actions/checkout@v3"
        step "Build" run = "make build"
    }
}


