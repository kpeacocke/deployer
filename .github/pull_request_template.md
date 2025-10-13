name: Pull Request
description: Submit a pull request
title: ""
labels: []
assignees: []

body:
  - type: markdown
    attributes:
      value: |
        Thank you for contributing to gh-deployer! Please fill out this template to help us review your pull request.

  - type: dropdown
    id: pr-type
    attributes:
      label: Type of Change
      description: What type of change does this PR introduce?
      options:
        - Bug fix (non-breaking change which fixes an issue)
        - New feature (non-breaking change which adds functionality)
        - Breaking change (fix or feature that would cause existing functionality to not work as expected)
        - Documentation update
        - Code refactoring
        - Performance improvement
        - Test improvement
        - CI/CD improvement
        - Other
    validations:
      required: true

  - type: textarea
    id: description
    attributes:
      label: Description
      description: Please describe your changes in detail
      placeholder: What does this PR do?
    validations:
      required: true

  - type: textarea
    id: motivation
    attributes:
      label: Motivation and Context
      description: Why is this change required? What problem does it solve?
      placeholder: This change is needed because...

  - type: textarea
    id: related-issues
    attributes:
      label: Related Issues
      description: Link any related issues here
      placeholder: Fixes #123, Closes #456

  - type: textarea
    id: testing
    attributes:
      label: Testing
      description: How has this been tested? Please describe your testing approach.
      placeholder: |
        - [ ] Unit tests pass
        - [ ] Integration tests pass
        - [ ] Manual testing performed
        - [ ] Tested on target platform (Raspberry Pi, etc.)

  - type: checkboxes
    id: checklist
    attributes:
      label: Checklist
      description: Please confirm the following
      options:
        - label: My code follows the code style of this project
          required: true
        - label: I have performed a self-review of my code
          required: true
        - label: I have commented my code, particularly in hard-to-understand areas
          required: true
        - label: I have made corresponding changes to the documentation
          required: false
        - label: My changes generate no new warnings
          required: true
        - label: I have added tests that prove my fix is effective or that my feature works
          required: false
        - label: New and existing unit tests pass locally with my changes
          required: true
        - label: Any dependent changes have been merged and published in downstream modules
          required: false

  - type: textarea
    id: screenshots
    attributes:
      label: Screenshots (if applicable)
      description: Add screenshots to help explain your changes

  - type: textarea
    id: additional-notes
    attributes:
      label: Additional Notes
      description: Any additional information that reviewers should know
