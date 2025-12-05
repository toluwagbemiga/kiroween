'use client';

import React, { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client';
import { Modal, Button, Loading } from '@/components/ui';
import { useToast } from '@/components/ui';
import {
  CANCEL_SUBSCRIPTION_MUTATION,
  UPDATE_SUBSCRIPTION_MUTATION,
} from '@/lib/graphql/mutations/billing';
import { BILLING_PORTAL_URL_QUERY, MY_SUBSCRIPTION_QUERY } from '@/lib/graphql/queries/billing';
import {
  CancelSubscriptionMutation,
  UpdateSubscriptionMutation,
  UpdateSubscriptionMutationVariables,
  BillingPortalUrlQuery,
} from '@/lib/graphql/generated/graphql';
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';

export interface SubscriptionManagementModalProps {
  isOpen: boolean;
  onClose: () => void;
  subscription: {
    id: string;
    status: string;
    cancelAtPeriodEnd: boolean;
    plan: {
      name: string;
    };
  };
}

export const SubscriptionManagementModal: React.FC<
  SubscriptionManagementModalProps
> = ({ isOpen, onClose, subscription }) => {
  const { success, error: showError } = useToast();
  const [showCancelConfirm, setShowCancelConfirm] = useState(false);

  // Query for billing portal URL
  const { data: portalData, loading: portalLoading } = useQuery<BillingPortalUrlQuery>(
    BILLING_PORTAL_URL_QUERY,
    {
      skip: !isOpen,
    }
  );

  // Cancel subscription mutation
  const [cancelSubscription, { loading: cancelLoading }] = useMutation<CancelSubscriptionMutation>(
    CANCEL_SUBSCRIPTION_MUTATION,
    {
      refetchQueries: [{ query: MY_SUBSCRIPTION_QUERY }],
      onCompleted: () => {
        success('Subscription will be canceled at the end of the billing period', 'Success');
        setShowCancelConfirm(false);
        onClose();
      },
      onError: (error) => {
        showError(error.message || 'Failed to cancel subscription', 'Error');
      },
    }
  );

  const handleCancelSubscription = async () => {
    await cancelSubscription();
  };

  const handleManageBilling = () => {
    if (portalData?.billingPortalUrl) {
      window.location.href = portalData.billingPortalUrl;
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="md"
      title="Manage Subscription"
      description="Update your subscription settings"
    >
      <div className="space-y-6">
        {/* Current Plan Info */}
        <div className="bg-white/5 rounded-lg p-4">
          <div className="text-sm text-gray-400 mb-1">Current Plan</div>
          <div className="text-lg font-semibold text-white">
            {subscription.plan.name}
          </div>
          <div className="text-sm text-gray-400 mt-1">
            Status: <span className="capitalize">{subscription.status}</span>
          </div>
        </div>

        {/* Actions */}
        <div className="space-y-3">
          {/* Billing Portal */}
          <Button
            variant="primary"
            className="w-full"
            onClick={handleManageBilling}
            disabled={portalLoading}
            isLoading={portalLoading}
          >
            Manage Billing & Payment Methods
          </Button>

          {/* Cancel Subscription */}
          {!subscription.cancelAtPeriodEnd && (
            <>
              {!showCancelConfirm ? (
                <Button
                  variant="danger"
                  className="w-full"
                  onClick={() => setShowCancelConfirm(true)}
                >
                  Cancel Subscription
                </Button>
              ) : (
                <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4">
                  <div className="flex items-start gap-3 mb-4">
                    <ExclamationTriangleIcon className="h-6 w-6 text-red-400 flex-shrink-0" />
                    <div>
                      <h4 className="text-white font-medium mb-1">
                        Cancel Subscription?
                      </h4>
                      <p className="text-sm text-gray-300">
                        Your subscription will remain active until the end of your
                        current billing period. You won't be charged again.
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Button
                      variant="danger"
                      onClick={handleCancelSubscription}
                      isLoading={cancelLoading}
                      disabled={cancelLoading}
                    >
                      Confirm Cancellation
                    </Button>
                    <Button
                      variant="ghost"
                      onClick={() => setShowCancelConfirm(false)}
                      disabled={cancelLoading}
                    >
                      Keep Subscription
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}

          {subscription.cancelAtPeriodEnd && (
            <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-4">
              <p className="text-sm text-yellow-300">
                Your subscription is scheduled to be canceled at the end of the
                current billing period.
              </p>
            </div>
          )}
        </div>

        {/* Close Button */}
        <div className="pt-4 border-t border-white/10">
          <Button variant="secondary" className="w-full" onClick={onClose}>
            Close
          </Button>
        </div>
      </div>
    </Modal>
  );
};
