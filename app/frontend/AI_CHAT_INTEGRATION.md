# AI Chat Widget Integration Guide

## Overview

The LLM gateway service has been successfully integrated into the frontend with an AI-powered chat widget. This provides users with instant, intelligent support powered by GPT-4.

## Architecture

### Backend (GraphQL Gateway & LLM Service)

**GraphQL Mutation:**
```graphql
mutation CallPrompt($name: String!, $variables: JSON!) {
  callPrompt(name: $name, variables: $variables) {
    content
    model
    tokensUsed
    cost
  }
}
```

**Implementation:**
- Location: `app/gateway/graphql-api-gateway/internal/resolvers/mutation.resolvers.go`
- Calls LLM gateway service via gRPC
- Passes user context and variables to prompt template
- Returns AI-generated response

### LLM Gateway Service

**Prompt Template:**
- Location: `prompts/v1/support-chatbot.md`
- Model: GPT-4
- Temperature: 0.7
- Max Tokens: 500
- Variables: `user_message`, `user_name`, `user_email`

### Frontend Integration

#### ChatWidget Component (`src/components/ChatWidget.tsx`)

Floating chat widget with AI-powered responses:

```typescript
import { ChatWidget } from '@/components/ChatWidget';

// Already added to root layout - available everywhere!
```

**Features:**
- Floating chat bubble in bottom-right corner
- Expandable chat window
- Message history
- Real-time AI responses
- Loading states
- Error handling
- Keyboard shortcuts (Enter to send)
- Auto-scroll to latest message
- Glassmorphism design

## Usage

### For Users

1. **Open Chat**: Click the floating chat bubble in the bottom-right corner
2. **Ask Question**: Type your question in the input field
3. **Send**: Press Enter or click the send button
4. **Get Response**: AI assistant responds instantly
5. **Continue Conversation**: Ask follow-up questions

### For Developers

The chat widget is automatically available on all pages. No additional setup needed!

#### Customizing the Prompt

Edit `prompts/v1/support-chatbot.md`:

```markdown
---
name: v1/support-chatbot
description: AI support chatbot
model: gpt-4
temperature: 0.7
max_tokens: 500
variables:
  - user_message
  - user_name
---

Your custom prompt here...

User's Question:
{{user_message}}
```

#### Using Different Prompts

```typescript
const { data } = await callPrompt({
  variables: {
    name: 'v1/custom-prompt',  // Change prompt name
    variables: {
      custom_var: 'value',
    },
  },
});
```

#### Creating Custom Chat Components

```typescript
import { gql, useMutation } from '@apollo/client';

const CALL_PROMPT_MUTATION = gql`
  mutation CallPrompt($name: String!, $variables: JSON!) {
    callPrompt(name: $name, variables: $variables) {
      content
      model
      tokensUsed
      cost
    }
  }
`;

function CustomChat() {
  const [callPrompt, { loading, data }] = useMutation(CALL_PROMPT_MUTATION);

  const askAI = async (question: string) => {
    const { data } = await callPrompt({
      variables: {
        name: 'v1/support-chatbot',
        variables: {
          user_message: question,
          user_name: 'John',
        },
      },
    });

    console.log('AI Response:', data.callPrompt.content);
  };

  return <button onClick={() => askAI('How do I upgrade?')}>Ask AI</button>;
}
```

## Prompt Templates

### Creating New Prompts

1. **Create prompt file** in `prompts/` directory:

```markdown
---
name: v1/my-prompt
description: My custom prompt
model: gpt-4
temperature: 0.7
max_tokens: 300
variables:
  - input_text
  - context
---

You are a helpful assistant.

Context: {{context}}
User Input: {{input_text}}

Please respond helpfully.
```

2. **Use in frontend**:

```typescript
await callPrompt({
  variables: {
    name: 'v1/my-prompt',
    variables: {
      input_text: 'Hello',
      context: 'Support chat',
    },
  },
});
```

### Available Prompts

| Prompt Name | Description | Variables |
|-------------|-------------|-----------|
| `v1/support-chatbot` | General support assistant | `user_message`, `user_name`, `user_email` |

### Prompt Best Practices

1. **Be Specific**: Clear instructions produce better results
2. **Set Context**: Provide relevant background information
3. **Use Variables**: Make prompts reusable with template variables
4. **Set Limits**: Use `max_tokens` to control response length
5. **Adjust Temperature**: 
   - 0.0-0.3: Focused, deterministic
   - 0.4-0.7: Balanced (recommended)
   - 0.8-1.0: Creative, varied

## Features

### Message History

Messages are stored in component state:

```typescript
interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}
```

### Loading States

Shows animated dots while waiting for AI response:

```typescript
{loading && (
  <div className="flex justify-start">
    <Loading size="sm" variant="dots" />
  </div>
)}
```

### Error Handling

Gracefully handles errors with user-friendly messages:

```typescript
catch (error) {
  const errorMessage: Message = {
    role: 'assistant',
    content: 'Sorry, I encountered an error. Please try again.',
  };
  setMessages((prev) => [...prev, errorMessage]);
}
```

### Keyboard Shortcuts

- **Enter**: Send message
- **Shift + Enter**: New line (not implemented yet)

## Customization

### Styling

The chat widget uses Tailwind CSS and can be customized:

```typescript
// Change position
className="fixed bottom-4 right-4"  // Bottom-right (default)
className="fixed bottom-4 left-4"   // Bottom-left
className="fixed top-4 right-4"     // Top-right

// Change size
className="w-96 h-[600px]"  // Default
className="w-80 h-[500px]"  // Smaller
className="w-[500px] h-[700px]"  // Larger

// Change colors
className="bg-white/10"  // Default glassmorphism
className="bg-gray-900"  // Solid dark
className="bg-white"     // Solid light
```

### Initial Message

Change the welcome message:

```typescript
const [messages, setMessages] = useState<Message[]>([
  {
    id: '1',
    role: 'assistant',
    content: 'Welcome! How can I assist you today?',  // Custom message
    timestamp: new Date(),
  },
]);
```

### Disable on Certain Pages

```typescript
// In specific page component
import { usePathname } from 'next/navigation';

export default function MyPage() {
  const pathname = usePathname();
  const showChat = pathname !== '/login';  // Hide on login page

  return (
    <div>
      {/* Page content */}
      {showChat && <ChatWidget />}
    </div>
  );
}
```

## Backend Configuration

### Environment Variables

**LLM Gateway Service:**
```bash
# OpenAI API Key
OPENAI_API_KEY=sk-...

# Model Configuration
DEFAULT_MODEL=gpt-4
DEFAULT_TEMPERATURE=0.7
DEFAULT_MAX_TOKENS=500

# Prompts Directory
PROMPTS_DIR=./prompts
```

### Cost Tracking

The LLM gateway tracks usage:

```typescript
{
  content: "AI response",
  model: "gpt-4",
  tokensUsed: 150,
  cost: 0.0045  // In USD
}
```

### Rate Limiting

Consider implementing rate limiting:

```go
// In LLM gateway service
if userCallsToday > 100 {
    return errors.New("daily limit exceeded")
}
```

## Performance

### Response Times

- **Average**: 2-5 seconds
- **Depends on**:
  - Model (GPT-4 slower than GPT-3.5)
  - Prompt length
  - Response length
  - OpenAI API load

### Optimization Tips

1. **Use GPT-3.5 for simple queries**: Faster and cheaper
2. **Limit max_tokens**: Shorter responses = faster
3. **Cache common responses**: Store FAQ answers
4. **Stream responses**: Show partial responses (future enhancement)

## Security

### Authentication

- User must be logged in to use chat
- JWT token required for GraphQL mutation
- User context passed to LLM service

### Data Privacy

- User messages are sent to OpenAI
- No sensitive data should be in prompts
- Consider data retention policies
- Comply with privacy regulations

### Best Practices

1. **Don't send PII**: Avoid sending sensitive personal information
2. **Sanitize inputs**: Validate and clean user messages
3. **Rate limit**: Prevent abuse
4. **Monitor costs**: Track OpenAI API usage
5. **Log conversations**: For quality and compliance

## Troubleshooting

### Chat Widget Not Appearing

**Problem:** Chat bubble doesn't show

**Solutions:**
1. Check if `<ChatWidget />` is in layout.tsx
2. Verify user is authenticated
3. Check browser console for errors
4. Ensure z-index isn't conflicting

### No AI Response

**Problem:** Message sent but no response

**Solutions:**
1. Check OpenAI API key in LLM service
2. Verify prompt file exists: `prompts/v1/support-chatbot.md`
3. Check LLM gateway service is running
4. Look at backend logs for errors
5. Verify GraphQL mutation is working

### Slow Responses

**Problem:** AI takes too long to respond

**Solutions:**
1. Use GPT-3.5 instead of GPT-4
2. Reduce `max_tokens` in prompt
3. Simplify prompt template
4. Check OpenAI API status
5. Consider caching common responses

### Error Messages

**Problem:** "Sorry, I encountered an error"

**Solutions:**
1. Check browser console for error details
2. Verify backend services are running
3. Check OpenAI API quota/billing
4. Review backend logs
5. Test GraphQL mutation directly

## Testing

### Manual Testing

1. **Start services**:
   ```bash
   # Terminal 1: LLM Gateway
   cd app/services/llm-gateway-service
   go run cmd/main.go

   # Terminal 2: GraphQL Gateway
   cd app/gateway/graphql-api-gateway
   go run cmd/main.go

   # Terminal 3: Frontend
   cd app/frontend
   npm run dev
   ```

2. **Test chat**:
   - Login to application
   - Click chat bubble
   - Send test message
   - Verify AI response

### Automated Testing

```typescript
// Mock the mutation
jest.mock('@apollo/client', () => ({
  useMutation: jest.fn(() => [
    jest.fn().mockResolvedValue({
      data: {
        callPrompt: {
          content: 'Mocked AI response',
          model: 'gpt-4',
          tokensUsed: 50,
          cost: 0.001,
        },
      },
    }),
    { loading: false },
  ]),
}));
```

## Future Enhancements

Consider adding:

- **Streaming Responses**: Show AI typing in real-time
- **Message Persistence**: Save chat history to database
- **File Uploads**: Allow users to share screenshots
- **Voice Input**: Speech-to-text for messages
- **Multi-language**: Support multiple languages
- **Suggested Questions**: Quick action buttons
- **Feedback**: Thumbs up/down on responses
- **Handoff to Human**: Escalate to support team
- **Context Awareness**: Remember previous conversations
- **Rich Responses**: Markdown, code blocks, links

## Cost Management

### Estimating Costs

**GPT-4:**
- Input: $0.03 per 1K tokens
- Output: $0.06 per 1K tokens
- Average chat: ~200 tokens = $0.01

**GPT-3.5-Turbo:**
- Input: $0.0015 per 1K tokens
- Output: $0.002 per 1K tokens
- Average chat: ~200 tokens = $0.0004

### Reducing Costs

1. **Use GPT-3.5** for simple queries
2. **Implement caching** for common questions
3. **Set token limits** on responses
4. **Rate limit users** to prevent abuse
5. **Monitor usage** and set alerts

## Analytics

Track chat usage:

```typescript
import { useAnalytics } from '@/lib/analytics';

const { trackEvent } = useAnalytics();

// Track chat opened
trackEvent('chat_opened');

// Track message sent
trackEvent('chat_message_sent', {
  messageLength: message.length,
});

// Track AI response received
trackEvent('chat_response_received', {
  tokensUsed: data.callPrompt.tokensUsed,
  cost: data.callPrompt.cost,
});
```

## Support

For issues or questions:
1. Check this documentation
2. Review backend logs
3. Test GraphQL mutation directly
4. Check OpenAI API status
5. Contact development team

The AI chat widget provides intelligent, instant support to your users, enhancing their experience and reducing support burden!
