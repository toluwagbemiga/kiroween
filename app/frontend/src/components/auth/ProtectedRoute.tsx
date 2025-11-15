import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { LoadingPage } from '@/components/ui';

export interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: string;
  redirectTo?: string;
}

/**
 * ProtectedRoute component that wraps pages requiring authentication
 * Redirects to login if user is not authenticated
 * Optionally checks for required role
 */
export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requiredRole,
  redirectTo = '/login',
}) => {
  const { user, loading, isAuthenticated, hasRole } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      // Not authenticated - redirect to login
      if (!isAuthenticated) {
        router.push(redirectTo);
        return;
      }

      // Authenticated but missing required role
      if (requiredRole && !hasRole(requiredRole)) {
        router.push('/unauthorized');
      }
    }
  }, [loading, isAuthenticated, requiredRole, hasRole, router, redirectTo]);

  // Show loading state while checking authentication
  if (loading) {
    return <LoadingPage />;
  }

  // Not authenticated
  if (!isAuthenticated) {
    return null;
  }

  // Authenticated but missing required role
  if (requiredRole && !hasRole(requiredRole)) {
    return null;
  }

  // Authenticated and authorized
  return <>{children}</>;
};

/**
 * Higher-order component to wrap pages with authentication
 */
export function withAuth<P extends object>(
  Component: React.ComponentType<P>,
  options?: { requiredRole?: string; redirectTo?: string }
) {
  return function AuthenticatedComponent(props: P) {
    return (
      <ProtectedRoute
        requiredRole={options?.requiredRole}
        redirectTo={options?.redirectTo}
      >
        <Component {...props} />
      </ProtectedRoute>
    );
  };
}
