workflow "Production Deploy" {
    on = "push"
    
    job "test" {
        runs-on = "ubuntu-latest"
        step "Checkout" uses = "actions/checkout@v3"
        step "Setup Go" {
            uses = "actions/setup-go@v3"  
            with {
                go-version = "1.20"
            }
        }
        step "Run Tests" run = "go test -v ./..."
    }
    
    job "build" {
        runs-on = "ubuntu-latest"
        step "Checkout" uses = "actions/checkout@v3"
        step "Build Binary" run = "go build -o bin/app ."
    }
    
    job "deploy" {
        runs-on = "ubuntu-latest"
        step "Deploy Application" {
            run = "echo 'Deploying to production...'"
            env {
                DEPLOY_ENV = "production"
                API_KEY = "${{ secrets.API_KEY }}"
            }
        }
    }
}
