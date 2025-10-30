# GitHub Actions DSL Example
name: Test Workflow with Advanced Features
on: 
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  NODE_VERSION: '18'
  BUILD_TYPE: production

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node: ['16', '18', '20']
        os: [ubuntu-latest, windows-latest]
    
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v3
        
      - name: Setup Node.js ${{ matrix.node }}
        uses: actions/setup-node@v3
        with:
          node-version: ${{ matrix.node }}
          cache: 'npm'
          
      - name: Install Dependencies
        run: |
          npm ci
          npm run build
        
      - name: Run Tests
        run: npm test
        env:
          CI: true
          
      - name: Conditional Step
        if: ${{ matrix.node == '18' && runner.os == 'Linux' }}
        run: echo "Running on Node 18 Linux"
        
  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
      - name: Deploy Application
        uses: custom/deploy-action@v1
        with:
          environment: ${{ env.BUILD_TYPE }}
          api-key: ${{ secrets.DEPLOY_KEY }}
