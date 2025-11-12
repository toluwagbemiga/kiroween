import React, { createContext, useContext } from 'react';
import { useAuth } from './AuthContext';
import { gql, useQuery } from '@apollo/client';

// GraphQL Query
const ME_QUERY = gql`
  query Me {
    me {
      id
      email
      name
      enabledFeatures
      roles {
        id
        name
      }
    }
  }
`;

export interface FeatureFlagContextType {
  enabledFeatures: string[];
  isFeatureEnabled: (featureName: string) => boolean;
  loading: boolean;
}

const FeatureFlagContext = createContext<FeatureFlagContextType | undefined>(undefined);

/**
 * FeatureFlagProvider component that fetches and manages feature flags
 */
export const FeatureFlagProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();

  // Fetch user data including enabled features
  const { data, loading } = useQuery(ME_QUERY, {
    skip: !isAuthenticated,
    fetchPolicy: 'cache-and-network',
  });

  const enabledFeatures = data?.me?.enabledFeatures || [];

  const isFeatureEnabled = (featureName: string): boolean => {
    return enabledFeatures.includes(featureName);
  };

  const value: FeatureFlagContextType = {
    enabledFeatures,
    isFeatureEnabled,
    loading,
  };

  return (
    <FeatureFlagContext.Provider value={value}>
      {children}
    </FeatureFlagContext.Provider>
  );
};

/**
 * Custom hook to use feature flags
 */
export const useFeatureFlags = (): FeatureFlagContextType => {
  const context = useContext(FeatureFlagContext);
  if (context === undefined) {
    throw new Error('useFeatureFlags must be used within a FeatureFlagProvider');
  }
  return context;
};
