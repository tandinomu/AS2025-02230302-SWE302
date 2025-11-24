# Practical: End-to-End Tests using Playwright

This practical demonstrates the implementation of automated end-to-end testing using Playwright for a React-based coding quiz application, including cross-browser testing capabilities.

## Overview

This project showcases:
* Comprehensive end-to-end test coverage using Playwright
* Cross-browser testing (Chromium, Firefox, WebKit)
* React component testing and validation
* Automated test execution with detailed reporting

## Screenshots

### 1. Quiz Application - Home Page
The main landing page of the coding quiz application:

![Quizhomepage](./images/quizss.png)

### 2. Quiz Question Interface
Active quiz question displaying React-specific knowledge testing:

![QuizQ](./images/playing.png)

### 3. Playwright Test Execution
Terminal output showing successful test execution:

![TestResults](./images/testallpassed.png)

### 4. Playwright HTML Test Report
Detailed browser-based test report showing comprehensive test coverage:

![HTMLTestReport](./images/testcasespassedinchrome.png)


**Test suites executed:**
- Data Validation Tests (data-validation.spec.ts)
  - TC014: Question Data Integrity (Chromium, Firefox, WebKit)
  - TC015: Score Calculation (Chromium, Firefox, WebKit)
  - TC015b: Mixed Correct/Incorrect Score Calculation (Chromium, Firefox, WebKit)
  - TC014b: Timer Data Validation (Chromium, Firefox, WebKit)

All tests ran successfully across multiple browsers ensuring cross-platform compatibility.

---

This practical demonstrates a robust testing workflow with cross-browser compatibility verification and comprehensive test coverage for a React quiz application.