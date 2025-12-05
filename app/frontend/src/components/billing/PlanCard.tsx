'use client';

import React from 'react';
import { Card, CardContent, Badge, Button } from '@/components/ui';
import { CheckIcon } from '@heroicons/react/24/outline';
import { CheckoutButton } from './CheckoutButton';

export interface PlanCardProps {
  plan: {
    id: string;
    name: string;
    description?: string | null;
    price: number;
    currency: string;
    interval: string;
    features: string[];
    isActive: boolean;
  };
  isCurrentPlan?: boolean;
  onSelectPlan?: (planId: string) => void;
  loading?: boolean;
  useCheckout?: boolean;
}

export const PlanCard: React.FC<PlanCardProps> = ({
  plan,
  isCurrentPlan = false,
  onSelectPlan,
  loading = false,
  useCheckout = false,
}) => {
  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency.toUpperCase(),
      minimumFractionDigits: 0,
    }).format(price / 100); // Assuming price is in cents
  };

  return (
    <Card
      className={`relative ${
        isCurrentPlan
          ? 'border-2 border-primary-500 shadow-lg shadow-primary-500/20'
          : ''
      }`}
    >
      {isCurrentPlan && (
        <div className="absolute -top-3 left-1/2 -translate-x-1/2">
          <Badge variant="success" size="sm">
            Current Plan
          </Badge>
        </div>
      )}

      <CardContent className="p-6 space-y-6">
        {/* Plan Header */}
        <div className="text-center">
          <h3 className="text-2xl font-bold text-white mb-2">{plan.name}</h3>
          {plan.description && (
            <p className="text-sm text-gray-400">{plan.description}</p>
          )}
        </div>

        {/* Pricing */}
        <div className="text-center py-4">
          <div className="flex items-baseline justify-center gap-1">
            <span className="text-4xl font-bold text-white">
              {formatPrice(plan.price, plan.currency)}
            </span>
            <span className="text-gray-400">/{plan.interval}</span>
          </div>
        </div>

        {/* Features */}
        <div className="space-y-3">
          {plan.features.map((feature, index) => (
            <div key={index} className="flex items-start gap-3">
              <CheckIcon className="h-5 w-5 text-green-400 flex-shrink-0 mt-0.5" />
              <span className="text-sm text-gray-300">{feature}</span>
            </div>
          ))}
        </div>

        {/* Action Button */}
        <div className="pt-4">
          {isCurrentPlan ? (
            <Button variant="secondary" className="w-full" disabled>
              Current Plan
            </Button>
          ) : useCheckout ? (
            <CheckoutButton
              planId={plan.id}
              planName={plan.name}
              variant="primary"
              className="w-full"
            >
              {plan.isActive ? 'Subscribe Now' : 'Not Available'}
            </CheckoutButton>
          ) : (
            <Button
              variant="primary"
              className="w-full"
              onClick={() => onSelectPlan?.(plan.id)}
              disabled={loading || !plan.isActive}
              isLoading={loading}
            >
              {plan.isActive ? 'Select Plan' : 'Not Available'}
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
