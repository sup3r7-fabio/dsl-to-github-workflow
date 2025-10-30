workflow "Multi-Environment CI/CD Pipeline" {
    on = "push"
    
    // Simple for loop - creates test_1, test_2, test_3
    for i in 1..3 {
        job "test_${i}" {
            runs-on = "ubuntu-latest"
            step "Setup ${i}" uses = "actions/setup-node@v3"
            step "Test ${i}" run = "npm test -- --suite=${i}"
        }
    }
    
    // Foreach loop - creates deploy_dev, deploy_staging, deploy_prod
    foreach env in ["dev", "staging", "prod"] {
        job "deploy_${env}" {
            runs-on = "ubuntu-latest" 
            step "Deploy to ${env}" run = "kubectl apply -f k8s/${env}/"
            step "Smoke test ${env}" run = "curl -f https://${env}.myapp.com/health"
        }
    }
    
    // Repeat loop - creates build_1, build_2, build_3, build_4, build_5
    repeat 5 {
        job "build_matrix" {
            runs-on = "ubuntu-latest"
            step "Build iteration" run = "make build-${iteration_0}"
            step "Archive artifacts" uses = "actions/upload-artifact@v3"
        }
    }
    
    // For loop with range - creates security_scan_18, security_scan_20, security_scan_22 
    for version in range(18..22) {
        job "security_scan_${version}" {
            runs-on = "ubuntu-latest"
            step "Setup Node ${version}" uses = "actions/setup-node@v3"
            step "Security scan" run = "npm audit --node-version=${version}"
        }
    }
}
