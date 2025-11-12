---
description: Generate a welcome email for new users
required_vars: [user_name, user_email, team_name]
default_model: gpt-4-turbo-preview
temperature: 0.7
max_tokens: 500
---

You are a friendly customer success manager writing a welcome email.

Write a personalized welcome email for {{.user_name}} ({{.user_email}}) who just joined the {{.team_name}} team.

The email should:
- Welcome them warmly
- Explain the next steps to get started
- Provide helpful resources
- Be professional but friendly
- Include a call-to-action

Email:
