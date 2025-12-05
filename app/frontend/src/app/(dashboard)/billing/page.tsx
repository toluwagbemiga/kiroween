'use client';

import { useQuery, useMutation } from '@apollo/client';
import { gql } from '@apollo/client';
import { useState } from 'react';

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
  const [selectedPlan, setSelectedPlan] = useState<string | null>(null);

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
    return <div className="p-6">Loading plans...</div>;
  }

  const subscription = data?.mySubscription;
  const plans = data?.plans || [];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Billing & Subscription</h1>
        <p className="text-gray-600 mt-2">
          Manage your subscription and billing
        </p>
      </div>

      {/* Current Subscription */}
      {subscription && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Current Subscription</h2>
          <div className="flex justify-between items-center">
            <div>
              <p className="font-medium">{subscription.plan.name}</p>
              <p className="text-sm text-gray-600">
                Status: <span className="capitalize">{subscription.status}</span>
              </p>
              <p className="text-sm text-gray-600">
                Renews: {new Date(subscription.currentPeriodEnd).toLocaleDateString()}
              </p>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold">${subscription.plan.price}</p>
              <p className="text-sm text-gray-600">per month</p>
            </div>
          </div>
        </div>
      )}

      {/* Available Plans */}
      <div>
        <h2 className="text-2xl font-semibold mb-4">Available Plans</h2>
        <div className="grid md:grid-cols-3 gap-6">
          {plans.map((plan: any) => (
            <div
              key={plan.id}
              className={`border rounded-lg p-6 ${
                selectedPlan === plan.id ? 'border-blue-500 ring-2 ring-blue-200' : 'border-gray-200'
              }`}
            >
              <h3 className="text-xl font-bold mb-2">{plan.name}</h3>
              <div className="mb-4">
                <span className="text-3xl font-bold">${plan.price}</span>
                <span className="text-gray-600">/{plan.interval}</span>
              </div>
              <ul className="space-y-2 mb-6">
                {plan.features.map((feature: string, idx: number) => (
                  <li key={idx} className="flex items-start">
                    <span className="text-green-500 mr-2">✓</span>
                    <span className="text-sm">{feature}</span>
                  </li>
                ))}
              </ul>
              <button
                onClick={() => handleSubscribe(plan.id)}
                className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
              >
                {subscription?.plan.name === plan.name ? 'Current Plan' : 'Subscribe'}
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Billing Portal */}
      {subscription && (
        <div className="bg-gray-50 rounded-lg p-6">
          <h3 className="font-semibold mb-2">Manage Billing</h3>
          <p className="text-sm text-gray-600 mb-4">
            Update payment method, view invoices, or cancel subscription
          </p>
          <button className="bg-gray-800 text-white px-4 py-2 rounded-md hover:bg-gray-900">
            Open Billing Portal
          </button>
        </div>
      )}
    </div>
  );
}
