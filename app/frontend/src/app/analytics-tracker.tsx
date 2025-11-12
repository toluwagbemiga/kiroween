'use client';

import { useEffect } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';
import { useAnalytics } from '@/lib/analytics';

/**
 * AnalyticsTracker component that tracks page views on route changes
 * This component should be placed in the root layout
 */
export const AnalyticsTracker: React.FC = () => {
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
