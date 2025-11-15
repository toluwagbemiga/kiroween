'use client';

import React, { useState, useRef, useEffect } from 'react';
import { gql, useMutation } from '@apollo/client';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Loading } from '@/components/ui';
import {
  ChatBubbleLeftRightIcon,
  XMarkIcon,
  PaperAirplaneIcon,
  SparklesIcon,
} from '@heroicons/react/24/outline';
import { cn } from '@/lib/utils';

// GraphQL Mutation
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

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

/**
 * ChatWidget - AI-powered support chat widget
 * Floating chat bubble in the corner of the screen
 */
export const ChatWidget: React.FC = () => {
  const { user } = useAuth();
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      role: 'assistant',
      content: 'Hi! I\'m your AI assistant. How can I help you today?',
      timestamp: new Date(),
    },
  ]);
  const [inputValue, setInputValue] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const [callPrompt, { loading }] = useMutation(CALL_PROMPT_MUTATION);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Focus input when chat opens
  useEffect(() => {
    if (isOpen) {
      inputRef.current?.focus();
    }
  }, [isOpen]);

  const handleSendMessage = async () => {
    if (!inputValue.trim() || loading) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: inputValue,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputValue('');

    try {
      const { data } = await callPrompt({
        variables: {
          name: 'v1/support-chatbot',
          variables: {
            user_message: inputValue,
            user_name: user?.name || 'User',
            user_email: user?.email || '',
          },
        },
      });

      if (data?.callPrompt) {
        const assistantMessage: Message = {
          id: (Date.now() + 1).toString(),
          role: 'assistant',
          content: data.callPrompt.content,
          timestamp: new Date(),
        };

        setMessages((prev) => [...prev, assistantMessage]);
      }
    } catch (error) {
      console.error('Failed to get AI response:', error);
      
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: 'Sorry, I encountered an error. Please try again or contact support.',
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, errorMessage]);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <>
      {/* Chat Window */}
      {isOpen && (
        <div className="fixed bottom-20 right-4 z-50 w-96 h-[600px] flex flex-col rounded-xl bg-white/10 backdrop-blur-md border border-white/20 shadow-2xl animate-scale-in">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-white/10">
            <div className="flex items-center space-x-2">
              <div className="h-8 w-8 rounded-full bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center">
                <SparklesIcon className="h-5 w-5 text-white" />
              </div>
              <div>
                <h3 className="text-sm font-semibold text-white">AI Assistant</h3>
                <p className="text-xs text-gray-400">Always here to help</p>
              </div>
            </div>
            <button
              onClick={() => setIsOpen(false)}
              className="p-1 rounded-lg text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
              aria-label="Close chat"
            >
              <XMarkIcon className="h-5 w-5" />
            </button>
          </div>

          {/* Messages */}
          <div className="flex-1 overflow-y-auto p-4 space-y-4 scrollbar-thin">
            {messages.map((message) => (
              <div
                key={message.id}
                className={cn(
                  'flex',
                  message.role === 'user' ? 'justify-end' : 'justify-start'
                )}
              >
                <div
                  className={cn(
                    'max-w-[80%] rounded-lg p-3 text-sm',
                    message.role === 'user'
                      ? 'bg-primary-600 text-white'
                      : 'bg-white/10 text-white border border-white/20'
                  )}
                >
                  {message.content}
                </div>
              </div>
            ))}
            
            {loading && (
              <div className="flex justify-start">
                <div className="bg-white/10 border border-white/20 rounded-lg p-3">
                  <Loading size="sm" />
                </div>
              </div>
            )}
            
            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <div className="p-4 border-t border-white/10">
            <div className="flex items-center space-x-2">
              <input
                ref={inputRef}
                type="text"
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Type your message..."
                disabled={loading}
                className={cn(
                  'flex-1 px-4 py-2 rounded-lg',
                  'bg-white/10 backdrop-blur-sm border border-white/20',
                  'text-white placeholder:text-gray-400',
                  'focus:outline-none focus:ring-2 focus:ring-primary-500',
                  'disabled:opacity-50 disabled:cursor-not-allowed'
                )}
              />
              <Button
                onClick={handleSendMessage}
                disabled={!inputValue.trim() || loading}
                variant="primary"
                size="md"
                className="p-2"
                aria-label="Send message"
              >
                <PaperAirplaneIcon className="h-5 w-5" />
              </Button>
            </div>
            <p className="mt-2 text-xs text-gray-400 text-center">
              Powered by AI • Responses may vary
            </p>
          </div>
        </div>
      )}

      {/* Chat Bubble Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={cn(
          'fixed bottom-4 right-4 z-50',
          'h-14 w-14 rounded-full',
          'bg-gradient-to-br from-primary-600 to-primary-700',
          'text-white shadow-lg hover:shadow-xl',
          'transition-all duration-200 hover:scale-110',
          'flex items-center justify-center',
          isOpen && 'scale-0'
        )}
        aria-label="Open chat"
      >
        <ChatBubbleLeftRightIcon className="h-6 w-6" />
      </button>
    </>
  );
};
