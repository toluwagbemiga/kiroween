'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { PageContainer } from '@/components/layout';
import { Card, CardHeader, CardTitle, CardContent, Button } from '@/components/ui';

export default function SettingsPage() {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState('profile');

  return (
    <PageContainer>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-white">Settings</h1>
          <p className="text-gray-400 mt-2">
            Manage your account and preferences
          </p>
        </div>

        {/* Tabs */}
        <div className="border-b border-gray-700">
          <nav className="-mb-px flex space-x-8">
            {['profile', 'notifications', 'security', 'team'].map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`
                  py-4 px-1 border-b-2 font-medium text-sm capitalize
                  ${activeTab === tab
                    ? 'border-primary-500 text-primary-400'
                    : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-600'
                  }
                `}
              >
                {tab}
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <Card variant="glass">
          <CardContent className="p-6">
            {activeTab === 'profile' && (
              <div className="space-y-4">
                <h2 className="text-xl font-semibold text-white">Profile Settings</h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300">Email</label>
                    <input
                      type="email"
                      value={user?.email || ''}
                      disabled
                      className="mt-1 block w-full rounded-md bg-gray-800 border-gray-700 text-gray-400 px-3 py-2"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-300">Name</label>
                    <input
                      type="text"
                      defaultValue={user?.name || ''}
                      className="mt-1 block w-full rounded-md bg-gray-800 border-gray-700 text-white px-3 py-2"
                    />
                  </div>
                  <Button variant="primary">Save Changes</Button>
                </div>
              </div>
            )}

            {activeTab === 'notifications' && (
              <div className="space-y-4">
                <h2 className="text-xl font-semibold text-white">Notification Preferences</h2>
                <div className="space-y-3">
                  {['Email notifications', 'Push notifications', 'In-app notifications'].map((pref) => (
                    <label key={pref} className="flex items-center">
                      <input type="checkbox" defaultChecked className="rounded" />
                      <span className="ml-2 text-gray-300">{pref}</span>
                    </label>
                  ))}
                </div>
              </div>
            )}

            {activeTab === 'security' && (
              <div className="space-y-4">
                <h2 className="text-xl font-semibold text-white">Security Settings</h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300">Current Password</label>
                    <input type="password" className="mt-1 block w-full rounded-md bg-gray-800 border-gray-700 text-white px-3 py-2" />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-300">New Password</label>
                    <input type="password" className="mt-1 block w-full rounded-md bg-gray-800 border-gray-700 text-white px-3 py-2" />
                  </div>
                  <Button variant="primary">Update Password</Button>
                </div>
              </div>
            )}

            {activeTab === 'team' && (
              <div className="space-y-4">
                <h2 className="text-xl font-semibold text-white">Team Settings</h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300">Team ID</label>
                    <input
                      type="text"
                      value={user?.teamId || 'No team'}
                      disabled
                      className="mt-1 block w-full rounded-md bg-gray-800 border-gray-700 text-gray-400 px-3 py-2"
                    />
                  </div>
                  <div className="flex items-center space-x-2">
                    <input type="checkbox" id="multiTenant" className="rounded" />
                    <label htmlFor="multiTenant" className="text-sm text-gray-300">
                      Enable Multi-Tenancy (SaaS Mode)
                    </label>
                  </div>
                  <p className="text-sm text-gray-400">
                    Multi-tenancy allows your team to manage multiple isolated customer environments.
                  </p>
                  <Button variant="primary">Save Team Settings</Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </PageContainer>
  );
}
