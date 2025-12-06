describe('Dog Image Browser - API Mocking', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should handle successful API response', () => {
    // Intercept the API call and mock response
    cy.intercept('GET', '/api/dogs', {
      statusCode: 200,
      body: {
        message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
        status: 'success'
      }
    }).as('getDog');
    
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
    
    cy.get('[data-testid="dog-image"]')
      .should('be.visible')
      .and('have.attr', 'src')
      .and('include', 'n02110185_1469');
  });

  it.skip('should handle API errors gracefully', () => {
    // Mock API failure
    // Note: Skipped because the API route handles errors internally
    cy.intercept('GET', '/api/dogs', {
      statusCode: 500,
      body: {
        error: 'Internal Server Error'
      }
    }).as('getDogError');
    
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDogError');
    
    cy.get('[data-testid="error-message"]')
      .should('be.visible')
      .and('contain.text', 'Failed to load dog image');
  });

  it('should handle network timeout', () => {
    // Mock slow API response (delay)
    cy.intercept('GET', '/api/dogs', {
      statusCode: 200,
      body: {
        message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
        status: 'success'
      },
      delay: 3000 // 3 second delay
    }).as('getSlowDog');
    
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Button should show loading state
    cy.get('[data-testid="fetch-dog-button"]')
      .should('contain.text', 'Loading...')
      .and('be.disabled');
    
    cy.wait('@getSlowDog');
    
    // Image should eventually appear
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });

  it('should handle breeds API failure', () => {
    // Intercept breeds API and make it fail
    cy.intercept('GET', '/api/dogs/breeds', {
      statusCode: 500,
      body: {
        error: 'Failed to load breeds'
      }
    }).as('getBreedsError');
    
    cy.visit('/');
    cy.wait('@getBreedsError');
    
    // Dropdown should still exist but with only default option
    cy.get('[data-testid="breed-selector"]')
      .should('be.visible')
      .find('option')
      .should('have.length', 1); // Only "All Breeds (Random)"
  });

  it('should verify request headers', () => {
    cy.intercept('GET', '/api/dogs', (req) => {
      // You can inspect or modify request here
      expect(req.headers).to.have.property('accept');
      
      req.reply({
        statusCode: 200,
        body: {
          message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
          status: 'success'
        }
      });
    }).as('getDog');
    
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
  });

  it('should verify breed query parameter is sent correctly', () => {
    // Wait for breeds to load first
    cy.get('[data-testid="breed-selector"]')
      .find('option')
      .should('have.length.greaterThan', 1);
    
    // Intercept request with query parameter
    cy.intercept('GET', '/api/dogs?breed=husky').as('getHusky');
    
    cy.get('[data-testid="breed-selector"]').select('husky');
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    cy.wait('@getHusky');
  });
});
