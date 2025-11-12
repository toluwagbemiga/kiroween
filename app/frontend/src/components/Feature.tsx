import React from 'react';
import { useFeatureFlags } from '@/contexts/FeatureFlagContext';

export interface FeatureProps {
  name: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

/**
 * Feature component that conditionally renders children based on feature flag
 * 
 * @example
 * <Feature name="bento-grid">
 *   <BentoGridLayout />
 * </Feature>
 * 
 * @example With fallback
 * <Feature name="new-feature" fallback={<OldFeature />}>
 *   <NewFeature />
 * </Feature>
 */
export const Feature: React.FC<FeatureProps> = ({ name, children, fallback = null }) => {
  const { isFeatureEnabled, loading } = useFeatureFlags();

  // While loading, don't render anything (or you could show a loading state)
  if (loading) {
    return null;
  }

  // Check if feature is enabled
  if (isFeatureEnabled(name)) {
    return <>{children}</>;
  }

  // Feature is not enabled, render fallback
  return <>{fallback}</>;
};

/**
 * Hook to check if a feature is enabled
 * Useful for conditional logic outside of JSX
 * 
 * @example
 * const isNewFeatureEnabled = useFeature('new-feature');
 * if (isNewFeatureEnabled) {
 *   // Do something
 * }
 */
export const useFeature = (featureName: string): boolean => {
  const { isFeatureEnabled } = useFeatureFlags();
  return isFeatureEnabled(featureName);
};
