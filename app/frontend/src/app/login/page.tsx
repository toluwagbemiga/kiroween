'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { LoginForm, RegisterForm } from '@/components/auth';
import { SparklesIcon } from '@heroicons/react/24/outline';
import { useAnalytics } from '@/lib/analytics';

export default function LoginPage() {
  const [showRegister, setShowRegister] = useState(false);
  const router = useRouter();
  const { trackEvent } = useAnalytics();

  const handleSuccess = () => {
    // Track successful login
    trackEvent('user_logged_in', {
      method: showRegister ? 'register' : 'login',
    });
    router.push('/dashboard');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center p-4">
      {/* Background decoration */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary-500/20 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl" />
      </div>

      {/* Content */}
      <div className="relative z-10 w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center h-16 w-16 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-700 mb-4">
            <SparklesIcon className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Haunted SaaS</h1>
          <p className="text-gray-300">
            {showRegister ? 'Create your account' : 'Welcome back'}
          </p>
        </div>

        {/* Form */}
        {showRegister ? (
          <RegisterForm
            onSuccess={handleSuccess}
            onLoginClick={() => setShowRegister(false)}
          />
        ) : (
          <LoginForm
            onSuccess={handleSuccess}
            onRegisterClick={() => setShowRegister(true)}
          />
        )}
      </div>
    </div>
  );
}
