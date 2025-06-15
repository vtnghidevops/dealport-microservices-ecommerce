# E-commerce Microservices Performance Testing Suite

This directory contains a comprehensive performance testing suite for the e-commerce microservices project using k6. The suite includes smoke tests, load tests, stress tests, and pipeline tests to ensure system performance and reliability.

## 📋 Table of Contents

- [Overview](#overview)
- [Test Types](#test-types)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Test Configuration](#test-configuration)
- [Results Analysis](#results-analysis)
- [CI/CD Integration](#cicd-integration)
- [Troubleshooting](#troubleshooting)

## 🎯 Overview

The performance testing suite validates the e-commerce microservices architecture under various load conditions:

- **Smoke Tests**: Basic functionality validation with minimal load
- **Load Tests**: Normal expected load simulation
- **Stress Tests**: Beyond-capacity testing to find breaking points
- **Pipeline Tests**: Quick validation for CI/CD pipelines

### Architecture Coverage

The tests cover all major services and endpoints:

- Authentication Service (login, registration, token validation)
- User Service (profile management, wishlist)
- Product Service (catalog browsing, search, filtering)
- Cart Service (add/remove items, coupon application)
- Checkout Service (order creation, payment processing)

## 🧪 Test Types

### 1. Smoke Test (`smoke-test.js`)

- **Purpose**: Verify basic functionality with minimal load
- **Load**: 1-2 virtual users for 40 seconds
- **Thresholds**:
  - 95% of requests under 2 seconds
  - Error rate under 10%
  - 90% of checks pass

### 2. Load Test (`load-test.js`)

- **Purpose**: Simulate normal expected load
- **Load**: Gradually increase to 50 users over 16 minutes
- **Thresholds**:
  - 95% of requests under 3 seconds
  - Error rate under 5%
  - 95% of checks pass

### 3. Stress Test (`stress-test.js`)

- **Purpose**: Test system beyond normal capacity
- **Load**: Up to 150 users with spike testing
- **Thresholds**:
  - 95% of requests under 5 seconds
  - Error rate under 10%
  - 90% of checks pass

### 4. Pipeline Test (`pipeline-test.js`)

- **Purpose**: Quick performance validation for CI/CD
- **Load**: Up to 10 users for 4 minutes
- **Thresholds**:
  - 95% of requests under 2.5 seconds
  - Error rate under 2%
  - 98% of checks pass

## 🔧 Prerequisites

### Required Software

- [k6](https://k6.io/docs/getting-started/installation/) (v0.40.0 or later)
- Node.js (for package management, optional)
- Docker and Docker Compose (for local testing)

### System Requirements

- **For Load Testing**: 4GB RAM, 2 CPU cores
- **For Stress Testing**: 8GB RAM, 4 CPU cores
- **Network**: Stable internet connection for remote environments

## 📦 Installation

### 1. Install k6

#### macOS (using Homebrew)

```bash
brew install k6
```

#### Windows (using Chocolatey)

```powershell
choco install k6
```

#### Linux (using package manager)

```bash
# Ubuntu/Debian
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# CentOS/RHEL
sudo yum install https://dl.k6.io/rpm/repo.rpm
sudo yum install k6
```

### 2. Setup Test Environment

#### Local Environment

```bash
# Start all services
docker-compose up -d

# Wait for services to be ready
sleep 60

# Verify services are running
curl http://localhost:58080/api/v1/health
```

#### Remote Environments

Update the base URLs in `config/test-config.js` for your staging/production environments.

## 🚀 Usage

### Command Line Interface

#### Linux/macOS

```bash
# Make script executable
chmod +x run-tests.sh

# Run all tests on local environment
./run-tests.sh

# Run specific test type
./run-tests.sh -t smoke
./run-tests.sh -t load
./run-tests.sh -t stress
./run-tests.sh -t pipeline

# Run on different environment
./run-tests.sh -e staging -t load

# Skip environment checks
./run-tests.sh -s -t smoke
```

#### Windows (PowerShell)

```powershell
# Run all tests on local environment
.\run-tests.ps1

# Run specific test type
.\run-tests.ps1 -TestType smoke
.\run-tests.ps1 -TestType load
.\run-tests.ps1 -TestType stress
.\run-tests.ps1 -TestType pipeline

# Run on different environment
.\run-tests.ps1 -Environment staging -TestType load

# Skip environment checks
.\run-tests.ps1 -SkipSetup -TestType smoke
```

### Direct k6 Execution

```bash
# Set environment and run test directly
ENVIRONMENT=local k6 run tests/smoke-test.js

# With JSON output for analysis
ENVIRONMENT=local k6 run --out json=results/smoke-results.json tests/smoke-test.js

# With custom thresholds
ENVIRONMENT=local k6 run --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" tests/load-test.js
```

### NPM Scripts (if using package.json)

```bash
npm run test:smoke
npm run test:load
npm run test:stress
npm run test:pipeline
npm run test:all
```

## ⚙️ Test Configuration

### Environment Configuration (`config/test-config.js`)

```javascript
export const config = {
  baseUrls: {
    local: "http://localhost:58080",
    staging: "https://staging-api.ecommerce.com",
    production: "https://api.ecommerce.com",
  },

  testUsers: {
    admin: { email: "admin@test.com", password: "admin123456" },
    customer: { email: "customer@test.com", password: "customer123456" },
  },

  thresholds: {
    smoke: {
      http_req_duration: ["p(95)<2000"],
      http_req_failed: ["rate<0.1"],
      checks: ["rate>0.9"],
    },
    // ... other configurations
  },
};
```

### Customizing Tests

#### Modifying Load Patterns

Edit the `options.stages` in each test file:

```javascript
export const options = {
  stages: [
    { duration: "2m", target: 10 }, // Ramp up
    { duration: "5m", target: 50 }, // Stay at load
    { duration: "2m", target: 0 }, // Ramp down
  ],
  // ...
};
```

#### Adding Custom Scenarios

Create new scenarios in `utils/test-scenarios.js`:

```javascript
export function customScenario(token) {
  // Your custom test logic here
  const response = http.get(`${getBaseUrl()}/api/v1/custom-endpoint`);
  check(response, {
    "custom check": (r) => r.status === 200,
  });
}
```

## 📊 Results Analysis

### Understanding k6 Output

#### Key Metrics

- **http_req_duration**: Response time statistics
- **http_req_failed**: Error rate percentage
- **http_reqs**: Total number of requests
- **vus**: Virtual users (concurrent users)
- **checks**: Percentage of successful checks

#### Sample Output

```
✓ products list loaded
✓ products response time OK
✓ categories loaded
✓ get cart successful

checks.........................: 98.50% ✓ 394   ✗ 6
data_received..................: 2.1 MB 35 kB/s
data_sent......................: 156 kB 2.6 kB/s
http_req_blocked...............: avg=1.2ms    min=0s      med=0s      max=45ms     p(90)=0s      p(95)=0s
http_req_connecting............: avg=0.4ms    min=0s      med=0s      max=15ms     p(90)=0s      p(95)=0s
http_req_duration..............: avg=245ms    min=12ms    med=156ms   max=2.1s     p(90)=456ms   p(95)=678ms
http_req_failed................: 2.50%  ✓ 10   ✗ 390
http_req_receiving.............: avg=1.2ms    min=0.1ms   med=0.8ms   max=12ms     p(90)=2.1ms   p(95)=3.2ms
http_req_sending...............: avg=0.1ms    min=0s      med=0s      max=2ms      p(90)=0s      p(95)=0.1ms
http_req_waiting...............: avg=244ms    min=11ms    med=155ms   max=2.1s     p(90)=455ms   p(95)=677ms
http_reqs......................: 400    6.67/s
iteration_duration.............: avg=8.9s     min=4.2s    med=8.1s    max=15.6s    p(90)=12.8s   p(95)=14.2s
iterations.....................: 60     1/s
vus............................: 1      min=1   max=10
vus_max........................: 10     min=10  max=10
```

### Performance Baselines

#### Acceptable Performance Targets

- **Response Time**: 95% of requests under 2 seconds for critical paths
- **Error Rate**: Less than 5% for normal operations
- **Throughput**: Minimum 100 requests/second for product browsing
- **Concurrent Users**: Support for 50+ concurrent users

#### Warning Indicators

- Response times consistently above 3 seconds
- Error rates above 10%
- Memory usage growing continuously
- Database connection pool exhaustion

### Results Storage

Results are stored in the `results/` directory with timestamps:

- `{test-type}_{environment}_{timestamp}.json` - Detailed k6 results
- `{test-type}_{environment}_{timestamp}_summary.txt` - Human-readable summary

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Performance Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  performance-test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Setup k6
        run: |
          sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6

      - name: Start services
        run: docker-compose up -d

      - name: Wait for services
        run: sleep 60

      - name: Run pipeline performance test
        run: |
          cd scripts/performance-tests
          chmod +x run-tests.sh
          ./run-tests.sh -t pipeline -s

      - name: Upload results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: performance-results
          path: scripts/performance-tests/results/
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any

    stages {
        stage('Setup') {
            steps {
                sh 'docker-compose up -d'
                sleep(60)
            }
        }

        stage('Performance Test') {
            steps {
                dir('scripts/performance-tests') {
                    sh 'chmod +x run-tests.sh'
                    sh './run-tests.sh -t pipeline -s'
                }
            }
        }

        stage('Results') {
            steps {
                archiveArtifacts artifacts: 'scripts/performance-tests/results/**/*', fingerprint: true

                // Fail build if performance thresholds not met
                script {
                    def exitCode = sh(script: 'grep -q "✓" scripts/performance-tests/results/*.json', returnStatus: true)
                    if (exitCode != 0) {
                        error("Performance tests failed - check results")
                    }
                }
            }
        }
    }

    post {
        always {
            sh 'docker-compose down'
        }
    }
}
```

### Performance Gates

Set up performance gates to prevent deployment of slow code:

```javascript
// In your test files, use strict thresholds for CI/CD
export const options = {
  thresholds: {
    http_req_duration: ["p(95)<2000"], // Fail if 95% > 2s
    http_req_failed: ["rate<0.02"], // Fail if error rate > 2%
    checks: ["rate>0.98"], // Fail if checks < 98%
  },
};
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Connection Refused Errors

```
ERRO[0001] Get "http://localhost:58080/api/v1/products": dial tcp [::1]:58080: connect: connection refused
```

**Solution:**

- Ensure all services are running: `docker-compose ps`
- Check service health: `curl http://localhost:58080/api/v1/health`
- Wait longer for services to start: `sleep 120`

#### 2. Authentication Failures

```
✗ login successful
✗ login response has token
```

**Solution:**

- Verify test user credentials in `config/test-config.js`
- Check if user registration is working
- Ensure authentication service is running

#### 3. High Error Rates

```
http_req_failed................: 25.00% ✓ 100  ✗ 300
```

**Solution:**

- Check service logs: `docker-compose logs [service-name]`
- Verify database connections
- Reduce load to identify breaking point
- Check resource usage: `docker stats`

#### 4. Slow Response Times

```
http_req_duration..............: avg=5.2s    p(95)=8.1s
```

**Solution:**

- Check database query performance
- Monitor resource usage (CPU, memory)
- Verify network connectivity
- Consider scaling services

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
# Run with debug output
k6 run --verbose tests/smoke-test.js

# Run with HTTP debug
k6 run --http-debug tests/smoke-test.js

# Run with custom log level
k6 run --log-output=stdout --logformat=json tests/smoke-test.js
```

### Performance Monitoring

Monitor system resources during tests:

```bash
# Monitor Docker containers
docker stats

# Monitor system resources
htop

# Monitor network
netstat -an | grep :58080

# Check database connections
docker-compose exec postgres-auth psql -U auth_user -d auth_db -c "SELECT count(*) FROM pg_stat_activity;"
```

## 📈 Best Practices

### Test Design

1. **Start Small**: Begin with smoke tests before running load tests
2. **Realistic Data**: Use realistic test data and user behaviors
3. **Gradual Ramp**: Always ramp up load gradually
4. **Recovery Time**: Allow system recovery between test runs
5. **Environment Isolation**: Use dedicated test environments

### Performance Targets

1. **Response Time**: Set realistic targets based on user expectations
2. **Error Rates**: Keep error rates low (< 5% for normal operations)
3. **Scalability**: Test beyond expected peak load
4. **Resource Usage**: Monitor CPU, memory, and database connections

### Continuous Improvement

1. **Baseline Tracking**: Maintain performance baselines
2. **Trend Analysis**: Track performance trends over time
3. **Regression Detection**: Catch performance regressions early
4. **Capacity Planning**: Use results for infrastructure planning

## 📚 Additional Resources

- [k6 Documentation](https://k6.io/docs/)
- [Performance Testing Best Practices](https://k6.io/docs/testing-guides/test-types/)
- [k6 Thresholds Guide](https://k6.io/docs/using-k6/thresholds/)
- [Load Testing Patterns](https://k6.io/docs/testing-guides/load-testing-patterns/)

## 🤝 Contributing

To contribute to the performance testing suite:

1. Fork the repository
2. Create a feature branch
3. Add or modify tests
4. Update documentation
5. Submit a pull request

### Adding New Tests

1. Create test file in `tests/` directory
2. Add scenario functions in `utils/test-scenarios.js`
3. Update configuration in `config/test-config.js`
4. Add test to runner scripts
5. Update this README

---

For questions or support, please contact the DevOps team or create an issue in the project repository.
