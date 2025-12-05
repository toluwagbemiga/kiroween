'use client';

import React from 'react';
import { useQuery } from '@apollo/client';
import { Modal, Loading } from '@/components/ui';
import { RoleSelector } from './RoleSelector';
import { PermissionMatrix } from './PermissionMatrix';
import { USER_QUERY } from '@/lib/graphql/queries/users';
import { UserQuery, UserQueryVariables } from '@/lib/graphql/generated/graphql';

export interface RoleManagementModalProps {
  isOpen: boolean;
  onClose: () => void;
  userId: string;
  userName?: string | null;
}

export const RoleManagementModal: React.FC<RoleManagementModalProps> = ({
  isOpen,
  onClose,
  userId,
  userName,
}) => {
  const { data, loading, refetch } = useQuery<UserQuery, UserQueryVariables>(
    USER_QUERY,
    {
      variables: { id: userId },
      skip: !isOpen,
      fetchPolicy: 'cache-and-network',
    }
  );

  const user = data?.user;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="xl"
      title={`Manage Roles - ${userName || 'User'}`}
      description="Assign roles and view permissions for this user"
    >
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loading size="lg" />
        </div>
      ) : !user ? (
        <div className="text-center py-12 text-gray-400">
          User not found
        </div>
      ) : (
        <div className="space-y-6">
          {/* Role Selector */}
          <RoleSelector
            userId={userId}
            currentRoles={user.roles}
            onRolesChanged={() => refetch()}
          />

          {/* Permission Matrix */}
          {user.roles.length > 0 && (
            <PermissionMatrix
              roles={user.roles}
              userPermissions={user.permissions}
            />
          )}
        </div>
      )}
    </Modal>
  );
};
