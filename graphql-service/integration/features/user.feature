Feature: user feedback management
  In order enter my feedback
  As a  user of the system
  I need to be able to save my details and feedback

  Scenario: User enters his details and feedback
    Given "John Doe" is a user
    When he add his details and feedback as below:
      | firstName | lastName | email       | jobTitle      | feedback |
      | John      | Doe      |  john@doe.com | Developer |  I like it |
    Then there should be a user called "John" saved in the system with feedback "I like it"
