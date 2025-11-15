import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Input, Card, CardHeader, CardTitle, CardContent } from '@/components/ui';
import { EnvelopeIcon, LockClosedIcon, UserIcon } from '@heroicons/react/24/outline';

export interface RegisterFormProps {
  onSuccess?: () => void;
  onLoginClick?: () => void;
}

export const RegisterForm: React.FC<RegisterFormProps> = ({ onSuccess, onLoginClick }) => {
  const { register, loading, error } = useAuth();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [formError, setFormError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    // Validation
    if (!name || !email || !password || !confirmPassword) {
      setFormError('Please fill in all fields');
      return;
    }

    if (password !== confirmPassword) {
      setFormError('Passwords do not match');
      return;
    }

    if (password.length < 8) {
      setFormError('Password must be at least 8 characters');
      return;
    }

    try {
      await register(email, password, name);
      onSuccess?.();
    } catch (err: any) {
      setFormError(err.message || 'Registration failed');
    }
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Create Account</CardTitle>
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
            type="text"
            label="Full Name"
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}

            required
            autoComplete="name"
            disabled={loading}
          />

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
            placeholder="At least 8 characters"
            value={password}
            onChange={(e) => setPassword(e.target.value)}

            required
            autoComplete="new-password"
            disabled={loading}
            helperText="Must be at least 8 characters"
          />

          <Input
            type="password"
            label="Confirm Password"
            placeholder="Re-enter your password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}

            required
            autoComplete="new-password"
            disabled={loading}
          />

          <div className="text-xs text-gray-400">
            By creating an account, you agree to our{' '}
            <a href="/terms" className="text-primary-400 hover:text-primary-300">
              Terms of Service
            </a>{' '}
            and{' '}
            <a href="/privacy" className="text-primary-400 hover:text-primary-300">
              Privacy Policy
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
            Create Account
          </Button>

          {onLoginClick && (
            <div className="text-center text-sm text-gray-300">
              Already have an account?{' '}
              <button
                type="button"
                onClick={onLoginClick}
                className="text-primary-400 hover:text-primary-300 transition-colors font-medium"
              >
                Sign in
              </button>
            </div>
          )}
        </form>
      </CardContent>
    </Card>
  );
};
