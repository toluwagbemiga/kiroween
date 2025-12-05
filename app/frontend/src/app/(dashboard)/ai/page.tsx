import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'AI Assistant | Haunted SaaS',
  description: 'AI-powered assistant',
};

export default function AIPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">AI Assistant</h1>
        <p className="text-gray-600 mt-2">
          Interact with AI-powered features
        </p>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🤖</div>
          <h2 className="text-2xl font-semibold mb-2">AI Features Coming Soon</h2>
          <p className="text-gray-600 max-w-md mx-auto">
            This section will include AI-powered prompts, LLM interactions, and intelligent automation features.
          </p>
        </div>
      </div>
    </div>
  );
}
