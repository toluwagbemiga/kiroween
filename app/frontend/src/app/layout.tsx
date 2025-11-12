import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import { Providers } from './providers';
import { AnalyticsTracker } from './analytics-tracker';
import { ChatWidget } from '@/components/ChatWidget';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Haunted SaaS - Modern SaaS Platform',
  description: 'A modern SaaS platform with AI-powered features',
  keywords: ['SaaS', 'AI', 'Platform', 'Analytics', 'Billing'],
  authors: [{ name: 'Haunted SaaS Team' }],
  viewport: 'width=device-width, initial-scale=1',
  themeColor: '#d946ef',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <Providers>
          <AnalyticsTracker />
          {children}
          <ChatWidget />
        </Providers>
      </body>
    </html>
  );
}
