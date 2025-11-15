import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Input, Card, CardHeader, CardTitle, CardContent } from '@/components/ui';
import { EnvelopeIcon, LockClosedIcon } from '@heroicons/react/24/outline';

export interface LoginFormProps {
  onSuccess?: () => void;
  onRegisterClick?: () => void;
}

export const LoginForm: React.FC<LoginFormProps> = ({ onSuccess, onRegisterClick }) => {
  const { login, loading, error } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [formError, setFormError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    // Validation
    if (!email || !password) {
      setFormError('Please fill in all fields');
      return;
    }

    try {
      await login(email, password);
      onSuccess?.();
    } catch (err: any) {
      setFormError(err.message || 'Login failed');
    }
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Welcome Back</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {(formError || error) && (
            <div
              className="p-3 rounded-lg bg-red-500/20 border border-red-500/30 text-red-300 text-sm"
              role="alert"
            >
              {formError || error}
            </div>
          )}

          <Input
            type="email"
            label="Email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}

            required
            autoComplete="email"
            disabled={loading}
          />

          <Input
            type="password"
            label="Password"
            placeholder="Enter your password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}

            required
            autoComplete="current-password"
            disabled={loading}
          />

          <div className="flex items-center justify-between text-sm">
            <label className="flex items-center text-gray-300">
              <input
                type="checkbox"
                className="mr-2 rounded border-white/20 bg-white/10"
              />
              Remember me
            </label>
            <a
              href="/forgot-password"
              className="text-primary-400 hover:text-primary-300 transition-colors"
            >
              Forgot password?
            </a>
          </div>

          <Button
            type="submit"
            variant="primary"
            size="lg"
            className="w-full"
            isLoading={loading}
            disabled={loading}
          >
            Sign In
          </Button>

          {onRegisterClick && (
            <div className="text-center text-sm text-gray-300">
              Don't have an account?{' '}
              <button
                type="button"
                onClick={onRegisterClick}
                className="text-primary-400 hover:text-primary-300 transition-colors font-medium"
              >
                Sign up
              </button>
            </div>
          )}
        </form>
      </CardContent>
    </Card>
  );
};
