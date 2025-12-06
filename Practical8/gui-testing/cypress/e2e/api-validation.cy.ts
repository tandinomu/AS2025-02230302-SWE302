describe('Dog Image Browser - API Response Validation', () => {
  it('should validate breeds API response structure', () => {
    cy.request('/api/dogs/breeds').then((response) => {
      // Check status code
      expect(response.status).to.eq(200);
      
      // Check response body structure
      expect(response.body).to.have.property('status', 'success');
      expect(response.body).to.have.property('message');
      expect(response.body.message).to.be.an('object');
      
      // Check that message contains breed data
      const breeds = Object.keys(response.body.message);
      expect(breeds).to.have.length.greaterThan(0);
    });
  });

  it('should validate random dog API response structure', () => {
    cy.request('/api/dogs').then((response) => {
      // Check status code
      expect(response.status).to.eq(200);
      
      // Check response body structure
      expect(response.body).to.have.property('status', 'success');
      expect(response.body).to.have.property('message');
      
      // Message should be a string (URL)
      expect(response.body.message).to.be.a('string');
      expect(response.body.message).to.include('https://');
      expect(response.body.message).to.include('dog.ceo');
    });
  });

  it('should validate specific breed API response', () => {
    cy.request('/api/dogs?breed=husky').then((response) => {
      // Check status code
      expect(response.status).to.eq(200);
      
      // Check response body structure
      expect(response.body).to.have.property('status', 'success');
      expect(response.body).to.have.property('message');
      
      // Message can be string or array
      if (Array.isArray(response.body.message)) {
        expect(response.body.message).to.have.length.greaterThan(0);
        expect(response.body.message[0]).to.include('husky');
      } else {
        expect(response.body.message).to.be.a('string');
        expect(response.body.message).to.include('husky');
      }
    });
  });
});
