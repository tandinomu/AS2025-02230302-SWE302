## Performance Testing Report: Practical 7 with k6

This report summarizes the defined procedure and criteria for the performance testing exercise for the Dog CEO API application using the k6 load testing tool, covering four distinct test scenarios.

---

### I. Procedure Summary

The practical involved setting up **k6** locally and configuring it for cloud execution via **Grafana Cloud**. The **Next.js** application, which integrated the **Dog CEO API**, served as the target system. Four mandatory performance test scenarios—**Average-Load**, **Spike-Load**, **Stress**, and **Soak**—were defined and planned for execution against the application to assess its performance under various conditions.

The tests were designed to utilize **Virtual Users (VUs)** and structured **Stages** to simulate different traffic patterns, incorporating **Checks** for response verification and **Thresholds** for defining overall test pass/fail criteria.

---

### II. Test Scenarios and Criteria

The following table defines the specific criteria and objectives for the four mandatory test scenarios. These criteria were to be clearly defined in the k6 test options for each scenario.

| Scenario | Duration | Expected VUs/Load | Key Thresholds (Success Criteria) | Objective |
| :--- | :--- | :--- | :--- | :--- |
| **Average-Load** | Custom (Varies) | Moderate/Constant VUs (e.g., 10-20) | $p(95) < 500\text{ms}$ Response Time, Error Rate $< 1\%$ | Measure baseline performance and throughput under typical user load. |
| **Spike-Load** | 1 minute | Extremely High VUs (Sudden Ramp-up) | $p(90) < 1000\text{ms}$ Response Time during spike, Error Rate $< 5\%$ | Test system resilience and recovery behavior during sudden, intense traffic surges. |
| **Stress** | 5 minutes | VUs near system capacity | $p(99) < 2000\text{ms}$ Response Time, Stability Check | Determine the breaking point and maximum operational capacity of the system. |
| **Soak** | 30 minutes | Low/Constant VUs (Long duration) | Minimal degradation in **avg** response time, Error Rate $< 1\%$ | Identify issues related to resource exhaustion, such as memory leaks, over an extended period. |

---

### III. Execution Status

Execution was planned for both local k6 CLI and cloud-based Grafana UI. Due to the absence of output data, the analysis of actual performance metrics (e.g., specific response times, error rates, and pass/fail status against thresholds) cannot be provided in this report.

---

### IV. Conclusion

The testing framework and criteria for assessing the performance of the Dog CEO API application against **Average-Load**, **Spike-Load**, **Stress**, and **Soak** scenarios have been successfully designed and documented. The implementation relies on k6 to enforce defined thresholds for response time and error rate, providing a formal evaluation of the application's stability and speed under various load profiles.