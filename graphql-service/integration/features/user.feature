Feature: user management
  In order use the system
  As a  user of the system
  I need to be able to save my details

  Scenario: User signs up
    Given "John Doe" is a user
    When he sign up with details below:
      | firstName | lastName | email       | jobTitle      |
      | John      | Doe      |  john@doe.com | Developer |
    Then there should be a user called "John"
