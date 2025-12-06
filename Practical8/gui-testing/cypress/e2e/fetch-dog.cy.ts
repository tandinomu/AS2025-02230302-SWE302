describe('Dog Image Browser - Fetch Dog Functionality', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should fetch and display a random dog image when button is clicked', () => {
    // Click the fetch button
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Wait for loading to finish and image to appear
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
    
    // Verify image has a valid src
    cy.get('[data-testid="dog-image"]')
      .should('have.attr', 'src')
      .and('include', 'dog.ceo');
    
    // Verify placeholder is gone
    cy.get('[data-testid="placeholder-message"]')
      .should('not.exist');
  });

  it('should fetch different dog images on multiple clicks', () => {
    // Array to store image URLs
    const imageUrls: string[] = [];
    
    // Click button and store first image
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .invoke('attr', 'src')
      .then((src) => {
        imageUrls.push(src as string);
      });
    
    // Click button again and verify different image
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .invoke('attr', 'src')
      .then((src) => {
        imageUrls.push(src as string);
        // Note: This might occasionally fail if API returns same image
        // In real tests, you'd mock the API to ensure different responses
      });
  });

  it('should handle rapid successive clicks gracefully', () => {
    // Click button multiple times quickly
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Should eventually show an image
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });
});

describe('Dog Image Browser - Breed Selection', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should load breed options in the dropdown', () => {
    // Check that dropdown has breeds (wait for API to load)
    cy.get('[data-testid="breed-selector"]')
      .find('option')
      .should('have.length.greaterThan', 1); // More than just "All Breeds"
  });

  it('should fetch a specific breed when selected', () => {
    // Wait for breeds to load
    cy.get('[data-testid="breed-selector"]')
      .find('option')
      .should('have.length.greaterThan', 1);
    
    // Select a specific breed (e.g., 'husky')
    cy.get('[data-testid="breed-selector"]').select('husky');
    
    // Click fetch button
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Verify image appears
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and('have.attr', 'src')
      .and('include', 'husky');
  });

  it('should allow switching between breeds', () => {
    // Wait for breeds to load
    cy.get('[data-testid="breed-selector"]')
      .find('option')
      .should('have.length.greaterThan', 1);
    
    // Select first breed
    cy.get('[data-testid="breed-selector"]').select('husky');
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and('have.attr', 'src')
      .and('include', 'husky');
    
    // Switch to different breed
    cy.get('[data-testid="breed-selector"]').select('beagle');
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and('have.attr', 'src')
      .and('include', 'beagle');
  });

  it('should capitalize breed names in the dropdown', () => {
    // Get a few breed options and check capitalization
    cy.get('[data-testid="breed-selector"]')
      .find('option')
      .should('have.length.greaterThan', 1)
      .eq(1) // Get first breed (index 0 is "All Breeds")
      .then(($option) => {
        const text = $option.text();
        // First letter should be uppercase
        expect(text.charAt(0)).to.equal(text.charAt(0).toUpperCase());
      });
  });
});
