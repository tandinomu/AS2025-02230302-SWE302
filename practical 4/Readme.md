## Practical Report: Setting up SAST with Snyk in GitHub Actions
[Repo Link](https://github.com/tandinomu/cicd-demo.git)

-----

### Procedure Summary

This practical established **Static Application Security Testing (SAST)** using **Snyk** integrated into a **GitHub Actions** CI/CD pipeline for a **Spring Boot/Maven** project (`cicd-demo`).

The key steps involved:

1.  **Prerequisites & Project Setup**: Cloning the `cicd-demo` repository and verifying the Java (v17+), Maven, and existing project structure/workflow (`maven.yml`).
2.  **Snyk Account & Token Configuration**: Creating a Snyk account, generating an **API Token**, and storing this token securely as a GitHub Repository Secret named `SNYK_TOKEN`.
3.  **Basic GitHub Actions Integration**: Utilizing the existing workflow's `security` job, which calls the `snyk/actions/maven@master` action, authenticated by the `SNYK_TOKEN`.
4.  **Enhanced Configuration**: Upgrading the workflow to include advanced arguments (e.g., `--severity-threshold=medium`), outputting results in **SARIF** format, and uploading them to **GitHub Code Scanning** via the `github/codeql-action/upload-sarif@v2` action.
5.  **Vulnerability Management (Hands-on)**: Demonstrating how to detect vulnerabilities (by introducing an intentionally old dependency), manage them using a `.snyk` ignore file, and ultimately fix them by updating the dependency version.
6.  **Advanced Strategies**: Implementing parallel scanning using a matrix strategy, setting up conditional scanning based on file changes (`dorny/paths-filter@v2`), and configuring weekly scheduled scans for continuous monitoring.

-----

### System and Tool Outputs

#### 1\. Project Setup Verification

| Command | Expected Output/Observation |
| :--- | :--- |
| `java -version` | Output should indicate version **17 or higher** |
| `mvn -version` | Output should show **Maven installation details** |
| `mvn clean compile` | **BUILD SUCCESS** |
| `mvn test` | **Tests Run: \> 0, Failures: 0** |

#### 2\. GitHub Secrets Configuration

| Name | Value | Status |
| :--- | :--- | :--- |
| `SNYK_TOKEN` | [Insert token value (Hidden)] | **Added/Verified** |

#### 3\. Basic Snyk Integration 

The initial workflow run using the basic Snyk configuration produced a security report


#### 4\. Sample Vulnerability Report (Interpreting Results)

A sample vulnerability detected by Snyk in a dependency:

```
✗ High severity vulnerability found in org.springframework:spring-core
  Description: Improper Input Validation
  Info: https://snyk.io/vuln/SNYK-JAVA-ORGSPRINGFRAMEWORK-1234567
  Remediation:
  ✓ Upgrade org.springframework.boot:spring-boot-starter-web to 3.1.3 or higher
```

#### 5\. Enhanced Snyk Workflow Configuration (Excerpt)

The enhanced configuration added SARIF upload for integration with GitHub Security Tab:

```yaml
- name: Run comprehensive Snyk scan
  uses: snyk/actions/maven@master
  env:
    SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
  with:
    args: --severity-threshold=medium --sarif-file-output=snyk.sarif

- name: Upload results to GitHub Security
  uses: github/codeql-action/upload-sarif@v2
  if: always()
  with:
    sarif_file: snyk.sarif
```

#### 6\. GitHub Security Tab Analysis

After the enhanced workflow ran:

| Tab | Result | Notes |
| :--- | :--- | :--- |
| **Security Tab** | **Code Scanning Alerts** | Snyk scan results displayed via SARIF upload. |

#### 7\. Vulnerability Management - `.snyk` file

To ignore a specific non-critical vulnerability temporarily:

```yaml
# .snyk file (Example)
version: v1.0.0
ignore:
  "SNYK-JAVA-COMFASTERXMLJACKSONCORE-1234567":
    - "*":
        reason: "Acceptable risk - not exploitable in our context"
        expires: "2024-12-31T23:59:59.999Z"
```

-----

### Conclusion and Takeaways

The practical successfully implemented an **automated SAST pipeline** using Snyk within GitHub Actions. This integration ensures that security vulnerabilities in the code and its dependencies are identified early in the development lifecycle. Key takeaways emphasize the importance of **automation**, setting a **Fail Fast Strategy** (e.g., using `--severity-threshold`), and **Continuous Monitoring** through scheduled and conditional scans.

## Evidence

**github workflow**
![s1](./images/Screenshot%202025-12-06%20at%204.28.44 PM.png)