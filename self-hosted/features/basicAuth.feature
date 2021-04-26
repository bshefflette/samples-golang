Feature: Basic Login with Username/Password

  Scenario: Login with username/Password
    Given I am an annymous user
    And I navigate to /login
    When I fill in my username
    And I fill in my Password
    And I submit the login form
    Then I should see my profile page details
