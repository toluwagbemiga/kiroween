# Prompt-as-Code Directory

This directory contains LLM prompt templates that are loaded dynamically by the LLM Gateway Service.

## Structure

Prompts are organized by feature area:
- `onboarding/` - User onboarding flows
- `notifications/` - Notification message generation
- `analytics/` - Analytics insights and summaries
- `support/` - Customer support responses

## Format

Prompts support YAML frontmatter for metadata:

```markdown
---
description: Generate a welcome email for new users
required_vars: [user_name, user_email, team_name]
default_model: gpt-4
---

You are a friendly customer success manager...

Write a personalized welcome email for {{.user_name}} ({{.user_email}})...
```

## Variable Syntax

Use `{{.variable_name}}` for simple variables (note the dot prefix required by Go templates) and `{{.object.property}}` for nested values.

## Hot Reloading

The LLM Gateway Service watches this directory and automatically reloads prompts when files change.
