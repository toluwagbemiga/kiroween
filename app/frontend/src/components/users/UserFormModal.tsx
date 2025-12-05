'use client';

import React from 'react';
import { useMutation } from '@apollo/client';
import { Modal } from '@/components/ui';
import { UserForm } from './UserForm';
import { UPDATE_USER_MUTATION } from '@/lib/graphql/mutations/users';
import { USERS_LIST_QUERY } from '@/lib/graphql/queries/users';
import { useToast } from '@/components/ui';
import {
  UpdateProfileMutation,
  UpdateProfileMutationVariables,
} from '@/lib/graphql/generated/graphql';

export interface UserFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  user?: {
    id: string;
    name?: string | null;
    email: string;
  };
}

export const UserFormModal: React.FC<UserFormModalProps> = ({
  isOpen,
  onClose,
  user,
}) => {
  const { success, error: showError } = useToast();
  const isEditing = !!user;

  const [updateUser, { loading }] = useMutation<
    UpdateProfileMutation,
    UpdateProfileMutationVariables
  >(UPDATE_USER_MUTATION, {
    refetchQueries: [{ query: USERS_LIST_QUERY }],
    onCompleted: () => {
      success(
        isEditing
          ? 'User updated successfully'
          : 'User created successfully',
        'Success'
      );
      onClose();
    },
    onError: (error) => {
      showError(
        error.message || 'Failed to save user',
        'Error'
      );
    },
  });

  const handleSubmit = async (data: { name: string; email: string }) => {
    await updateUser({
      variables: {
        input: {
          name: data.name,
          email: data.email,
        },
      },
    });
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="md">
      <UserForm
        user={user}
        onSubmit={handleSubmit}
        onCancel={onClose}
        isLoading={loading}
      />
    </Modal>
  );
};
