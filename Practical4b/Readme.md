## Practical Report: OWASP ZAP DAST Integration

> **Practical 4b**: Implementing DAST using OWASP ZAP in GitHub Actions
[Repo Link](https://github.com/tandinomu/SWE302_Practical4_cicd-demo.git)

-----

### Executive Summary

**Dynamic Application Security Testing (DAST)** using **OWASP ZAP** was successfully integrated into a GitHub Actions pipeline. This established an automated workflow for runtime vulnerability scanning against a containerized Spring Boot application.

### Key Achievements

  * **DAST Pipeline Automation**: Fully automated security testing implemented in GitHub Actions using `zaproxy/action-baseline`.
  * **Vulnerability Detection**: **19 runtime security vulnerabilities** were successfully detected and reported by the ZAP scan.
  * **Rule Customization**: Custom ZAP rules were configured via a `.zap/rules.tsv` file to enforce strict security policies, including failing the build on high-risk issues (e.g., SQL Injection, XSS) and missing security headers.
  * **Containerized Testing**: The application was built and deployed as a Docker container, running on port **3000**, enabling isolated and reliable black-box testing.

-----

### Implementation Details

#### 1\. OWASP ZAP Configuration

Custom security rules were defined in `.zap/rules.tsv` to control which alerts fail the build. High-risk issues and critical missing security headers were set to `FAIL`.

**ZAP Rules Configuration (Excerpt):**

| Rule ID | Risk Level | Action | Issue Type |
| :--- | :--- | :--- | :--- |
| 40018 | HIGH | FAIL | SQL Injection |
| 40012 | HIGH | FAIL | XSS (Reflected) |
| 10020 | MEDIUM | FAIL | X-Frame-Options Missing |
| 10038 | MEDIUM | FAIL | Content Security Policy Missing |

#### 2\. GitHub Actions DAST Workflow

The workflow (`zap-dast.yml`) ensures the application is running before the scan begins.

**Workflow Steps (Summary):**

1.  **Build and Deploy**: The application is built with Maven, containerized with Docker, and run in the background on `http://localhost:3000`. A timeout is used to ensure the application is live and responsive.
2.  **Run ZAP Baseline Scan**: The `zaproxy/action-baseline@v0.12.0` action is executed against the application target, referencing the custom `.zap/rules.tsv` file.

**Workflow Excerpt:**

```yaml
jobs:
  zap-baseline-scan:
    steps:
      - name: Build and deploy application
        run: |
          mvn clean package -DskipTests
          docker run -d --name cicd-demo-app -p 3000:3000 cicd-demo:latest
          # Wait for application health check
      - name: Run ZAP Baseline Scan
        uses: zaproxy/action-baseline@v0.12.0
        with:
          target: 'http://localhost:3000'
          rules_file_name: '.zap/rules.tsv'
```

-----

### DAST Analysis Results

The ZAP scan successfully identified vulnerabilities introduced via the `DastVulnerableController.java` endpoints, confirming the application's runtime security posture.

| Vulnerability Type | Risk Level | Remediation Impact |
| :--- | :--- | :--- |
| **Reflected XSS** | High | User input reflection without encoding. |
| **Directory Traversal** | High | Path parameter (`/download?filename=`) lacks validation. |
| **Missing Security Headers** | Medium | Absence of X-Frame-Options, CSP, HSTS. |
| **Open Redirect** | Medium | Unvalidated URL redirection parameter. |
| **Insecure Cookie Settings** | Low/Medium | Missing Secure, HttpOnly, and SameSite flags. |

**DAST Analysis Summary:**

  * **Total Issues Found**: **19 security vulnerabilities** (`WARN-NEW: 19`).
  * **Risk Breakdown**: 3 High, 8 Medium, 8+ Low.
  * **Security Checks Passed**: 128 successful validations.
  * **Coverage**: 100% of exposed application endpoints were scanned.

-----

### Security Findings Analysis

The primary DAST findings demonstrated failures in secure coding practices and server configuration:

1.  **Reflected XSS (High)**: Detected in endpoints that reflect un-sanitized user input (`/search`). **Remediation**: Implement output encoding.
2.  **Directory Traversal (High)**: Detected in file access functions (`/download`). **Remediation**: Implement strict input validation and file path normalization.
3.  **Missing Security Headers (Medium)**: Confirmed the lack of security-critical headers, leaving the application vulnerable to Clickjacking and XSS. **Remediation**: Configure the application or server to issue **X-Frame-Options: DENY**, **Content-Security-Policy**, and **Strict-Transport-Security**.

### DAST vs SAST

The DAST analysis confirmed exploitable vulnerabilities and misconfigurations that **SAST (Static Application Security Testing)** typically cannot detect, such as:

| Security Issue | DAST Detection | SAST Detection |
| :--- | :--- | :--- |
| **Missing Security Headers** | Runtime analysis of HTTP response. | Limited, focuses on source code structure. |
| **Insecure Cookie Flags** | Analysis of actual cookie header attributes. | Limited visibility of runtime cookie context. |
| **Open Redirect Behavior** | Actual HTTP request/response validation. | Limited to code logic without runtime context. |

DAST provides a critical layer of security by simulating real-world attacks, complementing the source code analysis provided by SAST tools.

-----

### Conclusion

The OWASP ZAP DAST integration successfully established an **automated runtime security validation** mechanism within the CI/CD pipeline. The pipeline identified **19 vulnerabilities**, providing actionable reports necessary for remediation. This implementation ensures a robust, continuous security posture, moving the project towards a comprehensive DevSecOps model.