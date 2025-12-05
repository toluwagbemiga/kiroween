'use client';

import { useQuery, useMutation } from '@apollo/client';
import { gql } from '@apollo/client';
import { useState } from 'react';
import { PageContainer } from '@/components/layout';
import { Card, CardHeader, CardTitle, CardContent, Button, Badge } from '@/components/ui';

const GET_PLANS = gql`
  query GetPlans {
    plans {
      id
      name
      description
      price
      currency
      interval
      features
      isActive
    }
    mySubscription {
      id
      status
      currentPeriodEnd
      cancelAtPeriodEnd
      plan {
        name
        price
      }
    }
  }
`;

const CREATE_CHECKOUT = gql`
  mutation CreateCheckout($planId: ID!) {
    createSubscriptionCheckout(planId: $planId) {
      url
      sessionId
    }
  }
`;

export default function BillingPage() {
  const { data, loading } = useQuery(GET_PLANS);
  const [createCheckout] = useMutation(CREATE_CHECKOUT);

  const handleSubscribe = async (planId: string) => {
    try {
      const result = await createCheckout({ variables: { planId } });
      if (result.data?.createSubscriptionCheckout?.url) {
        window.location.href = result.data.createSubscriptionCheckout.url;
      }
    } catch (error) {
      console.error('Checkout error:', error);
    }
  };

  if (loading) {
    return (
      <PageContainer>
        <div className="text-center py-12 text-gray-400">Loading plans...</div>
      </PageContainer>
    );
  }

  const subscription = data?.mySubscription;
  const plans = data?.plans || [];

  return (
    <PageContainer>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-white">Billing & Subscription</h1>
          <p className="text-gray-400 mt-2">
            Manage your subscription and billing
          </p>
        </div>

        {/* Current Subscription */}
        {subscription && (
          <Card variant="glass">
            <CardHeader>
              <CardTitle>Current Subscription</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex justify-between items-center">
                <div>
                  <p className="font-medium text-white">{subscription.plan.name}</p>
                  <p className="text-sm text-gray-400 mt-1">
                    Status: <Badge variant={subscription.status === 'active' ? 'success' : 'warning'}>
                      {subscription.status}
                    </Badge>
                  </p>
                  <p className="text-sm text-gray-400 mt-1">
                    Renews: {new Date(subscription.currentPeriodEnd).toLocaleDateString()}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-2xl font-bold text-white">${subscription.plan.price}</p>
                  <p className="text-sm text-gray-400">per month</p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Available Plans */}
        <div>
          <h2 className="text-2xl font-semibold mb-4 text-white">Available Plans</h2>
          <div className="grid md:grid-cols-3 gap-6">
            {plans.map((plan: any) => (
              <Card key={plan.id} variant="glass" hover>
                <CardContent className="p-6">
                  <h3 className="text-xl font-bold mb-2 text-white">{plan.name}</h3>
                  <div className="mb-4">
                    <span className="text-3xl font-bold text-white">${plan.price}</span>
                    <span className="text-gray-400">/{plan.interval}</span>
                  </div>
                  <ul className="space-y-2 mb-6">
                    {plan.features.map((feature: string, idx: number) => (
                      <li key={idx} className="flex items-start">
                        <span className="text-green-400 mr-2">✓</span>
                        <span className="text-sm text-gray-300">{feature}</span>
                      </li>
                    ))}
                  </ul>
                  <Button
                    variant="primary"
                    onClick={() => handleSubscribe(plan.id)}
                    className="w-full"
                    disabled={subscription?.plan.name === plan.name}
                  >
                    {subscription?.plan.name === plan.name ? 'Current Plan' : 'Subscribe'}
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Billing Portal */}
        {subscription && (
          <Card variant="glass">
            <CardContent className="p-6">
              <h3 className="font-semibold mb-2 text-white">Manage Billing</h3>
              <p className="text-sm text-gray-400 mb-4">
                Update payment method, view invoices, or cancel subscription
              </p>
              <Button variant="secondary">Open Billing Portal</Button>
            </CardContent>
          </Card>
        )}
      </div>
    </PageContainer>
  );
}
