import { useCallback } from 'react';
import { gql, useMutation } from '@apollo/client';
import { useAuth } from '@/contexts/AuthContext';

// GraphQL Mutations
const TRACK_EVENT_MUTATION = gql`
  mutation TrackEvent($input: TrackEventInput!) {
    trackEvent(input: $input)
  }
`;

const IDENTIFY_USER_MUTATION = gql`
  mutation IdentifyUser($properties: JSON!) {
    identifyUser(properties: $properties)
  }
`;

export interface TrackEventOptions {
  eventName: string;
  properties?: Record<string, any>;
}

export interface IdentifyUserOptions {
  properties: Record<string, any>;
}

/**
 * Custom hook for analytics tracking
 * Provides methods to track events and identify users
 */
export const useAnalytics = () => {
  const { user, isAuthenticated } = useAuth();
  const [trackEventMutation] = useMutation(TRACK_EVENT_MUTATION);
  const [identifyUserMutation] = useMutation(IDENTIFY_USER_MUTATION);

  /**
   * Track a custom event
   */
  const trackEvent = useCallback(
    async (eventName: string, properties?: Record<string, any>) => {
      // Only track if user is authenticated
      if (!isAuthenticated) {
        return;
      }

      try {
        // Add default properties
        const enrichedProperties = {
          ...properties,
          userId: user?.id,
          userEmail: user?.email,
          timestamp: new Date().toISOString(),
          url: typeof window !== 'undefined' ? window.location.href : undefined,
          pathname: typeof window !== 'undefined' ? window.location.pathname : undefined,
        };

        await trackEventMutation({
          variables: {
            input: {
              eventName,
              properties: enrichedProperties,
            },
          },
        });

        // Log in development
        if (process.env.NODE_ENV === 'development') {
          console.log('[Analytics] Event tracked:', eventName, enrichedProperties);
        }
      } catch (error) {
        console.error('[Analytics] Failed to track event:', error);
        // Don't throw - analytics failures shouldn't break the app
      }
    },
    [isAuthenticated, user, trackEventMutation]
  );

  /**
   * Identify user with additional properties
   */
  const identifyUser = useCallback(
    async (properties: Record<string, any>) => {
      if (!isAuthenticated) {
        return;
      }

      try {
        await identifyUserMutation({
          variables: {
            properties: {
              ...properties,
              userId: user?.id,
              email: user?.email,
              name: user?.name,
            },
          },
        });

        if (process.env.NODE_ENV === 'development') {
          console.log('[Analytics] User identified:', properties);
        }
      } catch (error) {
        console.error('[Analytics] Failed to identify user:', error);
      }
    },
    [isAuthenticated, user, identifyUserMutation]
  );

  /**
   * Track page view
   */
  const trackPageView = useCallback(
    (pageName?: string) => {
      const pathname = typeof window !== 'undefined' ? window.location.pathname : '';
      trackEvent('page_viewed', {
        pageName: pageName || pathname,
        referrer: typeof document !== 'undefined' ? document.referrer : undefined,
      });
    },
    [trackEvent]
  );

  /**
   * Track button click
   */
  const trackClick = useCallback(
    (buttonName: string, properties?: Record<string, any>) => {
      trackEvent('button_clicked', {
        buttonName,
        ...properties,
      });
    },
    [trackEvent]
  );

  /**
   * Track form submission
   */
  const trackFormSubmit = useCallback(
    (formName: string, properties?: Record<string, any>) => {
      trackEvent('form_submitted', {
        formName,
        ...properties,
      });
    },
    [trackEvent]
  );

  /**
   * Track error
   */
  const trackError = useCallback(
    (errorMessage: string, properties?: Record<string, any>) => {
      trackEvent('error_occurred', {
        errorMessage,
        ...properties,
      });
    },
    [trackEvent]
  );

  return {
    trackEvent,
    identifyUser,
    trackPageView,
    trackClick,
    trackFormSubmit,
    trackError,
    isEnabled: isAuthenticated,
  };
};
