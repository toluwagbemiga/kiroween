'use client';

import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { PlanCard } from '@/components/billing/PlanCard';
import { CheckoutButton } from '@/components/billing/CheckoutButton';
import { SubscriptionManagementModal } from '@/components/billing/SubscriptionManagementModal';
import { Card, CardContent, Loading, Badge, Button } from '@/components/ui';
import { EmptyState } from '@/components/layout/EmptyState';
import { PLANS_QUERY, MY_SUBSCRIPTION_QUERY } from '@/lib/graphql/queries/billing';
import {
  PlansQuery,
  MySubscriptionQuery,
} from '@/lib/graphql/generated/graphql';
import { format } from 'date-fns';
import {
  CreditCardIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';

export default function BillingPage() {
  const [isManagementModalOpen, setIsManagementModalOpen] = useState(false);

  const { data: plansData, loading: plansLoading, error: plansError } = useQuery<PlansQuery>(
    PLANS_QUERY
  );

  const {
    data: subscriptionData,
    loading: subscriptionLoading,
    error: subscriptionError,
  } = useQuery<MySubscriptionQuery>(MY_SUBSCRIPTION_QUERY);

  const handleSelectPlan = (planId: string) => {
    // For plan changes, open management modal
    setIsManagementModalOpen(true);
  };

  const handleManageSubscription = () => {
    setIsManagementModalOpen(true);
  };

  if (plansError || subscriptionError) {
    return (
      <DashboardLayout
        title="Billing"
        description="Manage your subscription and billing"
      >
        <EmptyState
          title="Error loading billing information"
          description={plansError?.message || subscriptionError?.message}
        />
      </DashboardLayout>
    );
  }

  const plans = plansData?.plans || [];
  const subscription = subscriptionData?.mySubscription;
  const currentPlanId = subscription?.planId;

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return <Badge variant="success">Active</Badge>;
      case 'canceled':
      case 'cancelled':
        return <Badge variant="error">Canceled</Badge>;
      case 'past_due':
        return <Badge variant="warning">Past Due</Badge>;
      case 'trialing':
        return <Badge variant="info">Trial</Badge>;
      default:
        return <Badge variant="default">{status}</Badge>;
    }
  };

  return (
    <DashboardLayout
      title="Billing"
      description="Manage your subscription and billing"
    >
      <div className="space-y-6">
        {/* Current Subscription Card */}
        {subscriptionLoading ? (
          <Card>
            <CardContent className="p-6">
              <Loading size="lg" />
            </CardContent>
          </Card>
        ) : subscription ? (
          <Card>
            <CardContent className="p-6">
              <div className="flex items-start justify-between mb-6">
                <div>
                  <h2 className="text-xl font-semibold text-white mb-2">
                    Current Subscription
                  </h2>
                  <p className="text-sm text-gray-400">
                    Manage your subscription and billing details
                  </p>
                </div>
                {getStatusBadge(subscription.status)}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {/* Plan Name */}
                <div className="bg-white/5 rounded-lg p-4">
                  <div className="text-sm text-gray-400 mb-1">Plan</div>
                  <div className="text-lg font-semibold text-white">
                    {subscription.plan.name}
                  </div>
                </div>

                {/* Price */}
                <div className="bg-white/5 rounded-lg p-4">
                  <div className="text-sm text-gray-400 mb-1">Price</div>
                  <div className="text-lg font-semibold text-white">
                    {new Intl.NumberFormat('en-US', {
                      style: 'currency',
                      currency: subscription.plan.currency.toUpperCase(),
                      minimumFractionDigits: 0,
                    }).format(subscription.plan.price / 100)}
                    <span className="text-sm text-gray-400">
                      /{subscription.plan.interval}
                    </span>
                  </div>
                </div>

                {/* Current Period */}
                <div className="bg-white/5 rounded-lg p-4">
                  <div className="text-sm text-gray-400 mb-1">Current Period</div>
                  <div className="text-sm font-medium text-white">
                    {format(new Date(subscription.currentPeriodStart), 'MMM d')} -{' '}
                    {format(new Date(subscription.currentPeriodEnd), 'MMM d, yyyy')}
                  </div>
                </div>

                {/* Renewal Status */}
                <div className="bg-white/5 rounded-lg p-4">
                  <div className="text-sm text-gray-400 mb-1">Renewal</div>
                  <div className="flex items-center gap-2">
                    {subscription.cancelAtPeriodEnd ? (
                      <>
                        <ExclamationTriangleIcon className="h-5 w-5 text-yellow-400" />
                        <span className="text-sm font-medium text-yellow-400">
                          Cancels on{' '}
                          {format(new Date(subscription.currentPeriodEnd), 'MMM d')}
                        </span>
                      </>
                    ) : (
                      <>
                        <CheckCircleIcon className="h-5 w-5 text-green-400" />
                        <span className="text-sm font-medium text-green-400">
                          Auto-renews
                        </span>
                      </>
                    )}
                  </div>
                </div>
              </div>

              {/* Plan Features */}
              {subscription.plan.features.length > 0 && (
                <div className="mt-6">
                  <h3 className="text-sm font-medium text-gray-300 mb-3">
                    Plan Features
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                    {subscription.plan.features.map((feature, index) => (
                      <div key={index} className="flex items-center gap-2">
                        <CheckCircleIcon className="h-4 w-4 text-green-400" />
                        <span className="text-sm text-gray-300">{feature}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Actions */}
              <div className="mt-6 flex items-center gap-3">
                <Button
                  variant="primary"
                  onClick={handleManageSubscription}
                >
                  Manage Subscription
                </Button>
              </div>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardContent className="p-6">
              <EmptyState
                icon={<CreditCardIcon className="h-12 w-12" />}
                title="No active subscription"
                description="Choose a plan below to get started"
              />
            </CardContent>
          </Card>
        )}

        {/* Available Plans */}
        <div>
          <h2 className="text-xl font-semibold text-white mb-4">
            {subscription ? 'Change Plan' : 'Choose a Plan'}
          </h2>
          
          {plansLoading ? (
            <div className="flex items-center justify-center py-12">
              <Loading size="lg" />
            </div>
          ) : plans.length === 0 ? (
            <EmptyState
              title="No plans available"
              description="Please check back later"
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {plans.map((plan: any) => (
                <PlanCard
                  key={plan.id}
                  plan={plan}
                  isCurrentPlan={plan.id === currentPlanId}
                  onSelectPlan={subscription ? handleSelectPlan : undefined}
                  useCheckout={!subscription}
                />
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Subscription Management Modal */}
      {subscription && (
        <SubscriptionManagementModal
          isOpen={isManagementModalOpen}
          onClose={() => setIsManagementModalOpen(false)}
          subscription={subscription}
        />
      )}
    </DashboardLayout>
  );
}
