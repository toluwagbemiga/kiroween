'use client';

import { useEffect, Suspense } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';
import { useAnalytics } from '@/lib/analytics';

/**
 * Internal tracker component that uses useSearchParams
 */
const AnalyticsTrackerInternal: React.FC = () => {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const { trackPageView, isEnabled } = useAnalytics();

  useEffect(() => {
    if (isEnabled && pathname) {
      // Track page view with pathname
      trackPageView(pathname);
    }
  }, [pathname, searchParams, trackPageView, isEnabled]);

  // This component doesn't render anything
  return null;
};

/**
 * AnalyticsTracker component that tracks page views on route changes
 * This component should be placed in the root layout
 * Wrapped in Suspense to support static export
 */
export const AnalyticsTracker: React.FC = () => {
  return (
    <Suspense fallback={null}>
      <AnalyticsTrackerInternal />
    </Suspense>
  );
};
