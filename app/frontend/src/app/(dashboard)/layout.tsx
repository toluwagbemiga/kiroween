'use client';

import { DashboardLayout } from '@/components/layout';
import { withAuth } from '@/components/auth';

function DashboardGroupLayout({ children }: { children: React.ReactNode }) {
  return (
    <DashboardLayout>
      {children}
    </DashboardLayout>
  );
}

export default withAuth(DashboardGroupLayout);
