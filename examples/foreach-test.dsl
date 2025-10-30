workflow "Foreach Loop Test" {
    on = "push"
    
    foreach env in ["dev", "staging", "prod"] {
        job "deploy_${env}" {
            runs-on = "ubuntu-latest"
            step "Deploy to ${env}" run = "kubectl apply -f k8s/${env}/"
        }
    }
}
