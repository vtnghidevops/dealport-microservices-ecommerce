# E-commerce Microservices DevSecOps Pipeline

## 🏗️ **Architecture Overview**

This is a **mono-repository** containing 9 microservices for an e-commerce platform with a comprehensive DevSecOps pipeline.

### **Microservices Architecture**

```
┌─────────────────┐    ┌──────────────────────────────────────┐
│   HTTP Client   │───▶│           Broker Service             │
└─────────────────┘    │        (HTTP Gateway)               │
                       └──────────────────┬───────────────────┘
                                          │ gRPC
                       ┌──────────────────▼───────────────────┐
                       │          gRPC Services               │
                       │  ┌─────────────────────────────────┐ │
                       │  │ • Authentication Service        │ │
                       │  │ • User Service                  │ │
                       │  │ • Product Service               │ │
                       │  │ • Cart Service                  │ │
                       │  │ • Checkout Service              │ │
                       │  │ • Mail Service                  │ │
                       │  │ • Logger Service                │ │
                       │  │ • Listener Service              │ │
                       │  └─────────────────────────────────┘ │
                       └──────────────────────────────────────┘
```

## 🚀 **DevSecOps Pipeline Stages**

### **10-Stage CI/CD Pipeline**

```yaml
stages:
  - unit-test # Go unit testing
  - build # Docker image building
  - scan-code # SAST security scanning
  - push-image # Harbor registry push
  - scan-image # Container image scanning
  - deploy-dev # Development deployment
  - deploy-staging # Staging deployment
  - dast-scan # Dynamic security testing
  - performance-test # Load/stress testing
  - deploy-prod # Production deployment
  - push-artifacts # MinIO artifact storage
```

## 🔒 **Security Testing Coverage**

### **SAST (Static Application Security Testing)**

- **GitLeaks**: Secret detection in source code
- **Trivy FS**: Filesystem vulnerability scanning
- **Snyk**: Dependency & code vulnerability scanning
- **SonarQube**: Code quality & security analysis

### **DAST (Dynamic Application Security Testing)**

- **OWASP ZAP Baseline**: Passive security scanning
- **OWASP ZAP Full Scan**: Comprehensive active testing
- **Security Tests**: SQL injection, XSS, CSRF, authentication issues

### **Container Security**

- **Trivy Image**: Container image vulnerability scanning
- **SBOM Generation**: Software Bill of Materials

## 📊 **Performance Testing**

### **K6 Performance Tests**

- **Smoke Test**: Basic functionality validation
- **Load Test**: Normal traffic simulation
- **Stress Test**: High traffic & breaking point testing

## 🗄️ **Artifact Management**

### **MinIO Storage Structure**

```
minio/bucket/project-name/YYYYMMDD_branch_commit/
├── sbom/                    # Software Bill of Materials
├── security-scans/          # SAST scan results
│   ├── secrets/            # GitLeaks reports
│   ├── trivy-image/        # Container scan reports
│   ├── trivy-fs/           # Filesystem scan reports
│   └── snyk/               # Dependency scan reports
├── dast-reports/           # DAST scan results
│   ├── zap-baseline-report.html
│   ├── zap-full-scan-report.html
│   └── dast-summary.txt
├── performance-tests/       # K6 test results
├── test-reports/           # Unit test reports
├── code-quality/           # SonarQube reports
└── manifest.json           # Artifact inventory
```

## 🏷️ **Deployment Strategy**

### **Git Tag-Based Deployment**

- `v1.0.0-dev`: Deploy to **Development**
- `v1.0.0-staging`: Deploy to **Staging**
- `v1.0.0-prod`: Deploy to **Production**

### **Environment Progression**

```
Development → Staging → Production
     ↓           ↓          ↓
  Auto Deploy  Manual    Manual
```

## 🛠️ **Technology Stack**

### **Backend**

- **Language**: Go (Golang)
- **Architecture**: gRPC microservices
- **Gateway**: HTTP-to-gRPC broker service
- **Database**: PostgreSQL/MongoDB
- **Message Queue**: RabbitMQ/Kafka

### **DevOps Tools**

- **CI/CD**: GitLab CI/CD
- **Container Registry**: Harbor
- **Orchestration**: Kubernetes
- **Monitoring**: Prometheus + Grafana
- **Artifact Storage**: MinIO
- **Security Scanning**: OWASP ZAP, Trivy, Snyk, GitLeaks

## 🚦 **Pipeline Exit Codes**

### **Security Scan Results**

- **0**: ✅ No issues - Safe to deploy
- **1**: ⚠️ Medium risk - Review recommended
- **2**: 🚫 Critical issues - **DEPLOYMENT BLOCKED**

## 📋 **Getting Started**

### **Prerequisites**

```bash
# Required tools
- Docker & Docker Compose
- Kubernetes cluster
- GitLab Runner
- Harbor registry access
- MinIO storage access
```

### **Environment Variables**

```yaml
# Container Registry
HARBOR_URL: "harbor.example.com"
HARBOR_USERNAME: "robot$username"
HARBOR_PASSWD: "password"

# Deployment
GITHUB_TOKEN: "github_token"
GITHUB_EMAIL: "user@example.com"
DEV_URL: "https://dev.example.com"
STAGING_URL: "https://staging.example.com"
PROD_URL: "https://prod.example.com"

# Security
SNYK_TOKEN: "snyk_token"
SONAR_HOST_URL: "https://sonar.example.com"

# Artifact Storage
MINIO_URL: "https://minio.example.com"
MINIO_ACCESS_KEY: "access_key"
MINIO_SECRET_KEY: "secret_key"
MINIO_BUCKET: "artifacts"
```

### **Running the Pipeline**

```bash
# 1. Create and push a tag
git tag v1.0.0-dev
git push origin v1.0.0-dev

# 2. Pipeline will automatically:
#    - Run unit tests
#    - Build Docker images
#    - Perform security scans
#    - Deploy to development
#    - Generate artifacts

# 3. For staging deployment
git tag v1.0.0-staging
git push origin v1.0.0-staging
# Manual approval required for staging

# 4. For production deployment
git tag v1.0.0-prod
git push origin v1.0.0-prod
# Manual approval required for production
```

## 📈 **Monitoring & Reporting**

### **Security Reports**

- **Real-time**: GitLab Security Dashboard
- **Historical**: MinIO artifact storage
- **Alerts**: Critical vulnerabilities block deployment

### **Performance Metrics**

- **Response Time**: P95, P99 percentiles
- **Throughput**: Requests per second
- **Error Rate**: 4xx/5xx responses
- **Resource Usage**: CPU, Memory, Network

### **Quality Gates**

- **Unit Test Coverage**: >80%
- **Security Scan**: No critical vulnerabilities
- **Performance**: Response time <500ms
- **Code Quality**: SonarQube quality gate passed

## 🔧 **Maintenance**

### **Regular Tasks**

- Update security scanning tools monthly
- Review and update dependency versions
- Monitor artifact storage usage
- Performance baseline updates

### **Troubleshooting**

- Check GitLab CI/CD logs for pipeline failures
- Review security scan reports for vulnerability details
- Monitor Kubernetes cluster health
- Verify MinIO artifact uploads

## 📞 **Support**

### **Team Contacts**

- **DevOps Team**: devops@company.com
- **Security Team**: security@company.com
- **Development Team**: dev@company.com

### **Documentation**

- **Pipeline Configuration**: `ci/templates/.gitlab-ci.yml`
- **Security Policies**: `docs/security/`
- **Deployment Guides**: `docs/deployment/`
- **Performance Baselines**: `docs/performance/`

---

## 🏆 **Key Features**

✅ **Comprehensive Security**: SAST + DAST + Container scanning  
✅ **Performance Testing**: Load, stress, and smoke tests  
✅ **Artifact Management**: Long-term storage and traceability  
✅ **Multi-Environment**: Dev → Staging → Production  
✅ **Quality Gates**: Automated blocking of vulnerable deployments  
✅ **Microservices Ready**: gRPC + HTTP gateway architecture  
✅ **Cloud Native**: Kubernetes + Harbor + MinIO  
✅ **Monitoring**: Full observability stack

**Built with ❤️ for secure, scalable e-commerce microservices**
