import { test, expect } from "@playwright/test";

// TC006-TC007: Timer Tests
test.describe("Timer Tests", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("http://localhost:5173");
    await page.waitForLoadState("networkidle");
  });

  test("TC006: Timer Countdown", async ({ page }) => {
    // Start the quiz
    await page.click("text=Start Quiz");
    await page.waitForSelector('[data-testid="timer"]');

    // Verify timer starts at 30
    await expect(page.locator('[data-testid="timer"]')).toContainText("30");

    // Wait 2 seconds and verify countdown
    await page.waitForTimeout(2000);
    const timerText = await page.locator('[data-testid="timer"]').textContent();
    const currentTime = parseInt(timerText?.match(/\d+/)?.[0] || "0");

    // Should be around 28 seconds (allowing for some variance)
    expect(currentTime).toBeGreaterThanOrEqual(27);
    expect(currentTime).toBeLessThanOrEqual(29);

    // Wait another 3 seconds and verify continued countdown
    await page.waitForTimeout(3000);
    const newTimerText = await page
      .locator('[data-testid="timer"]')
      .textContent();
    const newCurrentTime = parseInt(newTimerText?.match(/\d+/)?.[0] || "0");

    // Should be around 25 seconds
    expect(newCurrentTime).toBeGreaterThanOrEqual(24);
    expect(newCurrentTime).toBeLessThanOrEqual(26);
    expect(newCurrentTime).toBeLessThan(currentTime);
  });

  test("TC007: Timer Expiry", async ({ page }) => {
    // Increase test timeout for this specific test since we need to wait for timer
    test.setTimeout(60000);

    // Start the quiz
    await page.click("text=Start Quiz");
    await page.waitForSelector('[data-testid="timer"]', { timeout: 5000 });

    // Wait for timer to reach 0 and check for game over state
    // Check multiple possible end-game indicators
    await page.waitForFunction(
      () => {
        // Check for game-over element
        const gameOver = document.querySelector('[data-testid="game-over"]');
        if (gameOver) return true;

        // Check for final-score element
        const finalScore = document.querySelector('[data-testid="final-score"]');
        if (finalScore) return true;

        // Check for results text
        const resultsText = document.body.innerText;
        if (resultsText.includes("Quiz Complete") || 
            resultsText.includes("Time's Up") ||
            resultsText.includes("Game Over") ||
            resultsText.includes("Your Score")) {
          return true;
        }

        // Check if timer reached 0
        const timer = document.querySelector('[data-testid="timer"]');
        if (timer) {
          const time = parseInt(timer.textContent?.match(/\d+/)?.[0] || "30");
          if (time === 0) {
            // Give it a moment to transition
            return true;
          }
        }

        return false;
      },
      { timeout: 35000 } // 35 seconds to account for the 30 second timer plus transition
    );

    // Wait a bit for any transitions to complete
    await page.waitForTimeout(1000);

    // Try to find game over indicators - use more flexible selectors
    const gameOverVisible = await page.locator('[data-testid="game-over"]').isVisible().catch(() => false);
    const finalScoreVisible = await page.locator('[data-testid="final-score"]').isVisible().catch(() => false);
    const startQuizVisible = await page.locator('text=Start Quiz').isVisible().catch(() => false);
    const playAgainVisible = await page.locator('text=Play Again').isVisible().catch(() => false);

    // At least one end-game indicator should be visible
    const gameEnded = gameOverVisible || finalScoreVisible || startQuizVisible || playAgainVisible;
    expect(gameEnded).toBeTruthy();

    // Verify timer is no longer counting (either at 0 or not visible)
    const timerExists = await page.locator('[data-testid="timer"]').isVisible().catch(() => false);
    if (timerExists) {
      const finalTimerText = await page.locator('[data-testid="timer"]').textContent();
      const finalTime = parseInt(finalTimerText?.match(/\d+/)?.[0] || "0");
      expect(finalTime).toBeLessThanOrEqual(1); // Should be 0 or 1
    }
  });
});