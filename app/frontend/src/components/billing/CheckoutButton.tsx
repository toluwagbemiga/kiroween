'use client';

import React from 'react';
import { useMutation } from '@apollo/client';
import { Button } from '@/components/ui';
import { useToast } from '@/components/ui';
import { CREATE_SUBSCRIPTION_CHECKOUT_MUTATION } from '@/lib/graphql/mutations/billing';
import {
  CreateSubscriptionCheckoutMutation,
  CreateSubscriptionCheckoutMutationVariables,
} from '@/lib/graphql/generated/graphql';

export interface CheckoutButtonProps {
  planId: string;
  planName: string;
  children?: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  className?: string;
}

export const CheckoutButton: React.FC<CheckoutButtonProps> = ({
  planId,
  planName,
  children = 'Subscribe',
  variant = 'primary',
  className,
}) => {
  const { success, error: showError } = useToast();

  const [createCheckout, { loading }] = useMutation<
    CreateSubscriptionCheckoutMutation,
    CreateSubscriptionCheckoutMutationVariables
  >(CREATE_SUBSCRIPTION_CHECKOUT_MUTATION, {
    onCompleted: (data) => {
      // Redirect to Stripe Checkout
      if (data.createSubscriptionCheckout.url) {
        window.location.href = data.createSubscriptionCheckout.url;
      }
    },
    onError: (error) => {
      showError(
        error.message || 'Failed to create checkout session',
        'Checkout Error'
      );
    },
  });

  const handleCheckout = async () => {
    try {
      await createCheckout({
        variables: {
          planId,
        },
      });
    } catch (err) {
      console.error('Checkout error:', err);
    }
  };

  return (
    <Button
      variant={variant}
      onClick={handleCheckout}
      isLoading={loading}
      disabled={loading}
      className={className}
    >
      {children}
    </Button>
  );
};
