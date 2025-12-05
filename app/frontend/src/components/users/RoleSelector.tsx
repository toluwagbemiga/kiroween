'use client';

import React, { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { Badge, Button, Loading } from '@/components/ui';
import { ROLES_QUERY } from '@/lib/graphql/queries/users';
import { ASSIGN_ROLE_MUTATION, REMOVE_ROLE_MUTATION } from '@/lib/graphql/mutations/users';
import { useToast } from '@/components/ui';
import {
  RolesQuery,
  AssignRoleMutation,
  AssignRoleMutationVariables,
  RemoveRoleMutation,
  RemoveRoleMutationVariables,
} from '@/lib/graphql/generated/graphql';
import { PlusIcon, XMarkIcon } from '@heroicons/react/24/outline';

export interface RoleSelectorProps {
  userId: string;
  currentRoles: Array<{
    id: string;
    name: string;
    description?: string | null;
  }>;
  onRolesChanged?: () => void;
}

export const RoleSelector: React.FC<RoleSelectorProps> = ({
  userId,
  currentRoles,
  onRolesChanged,
}) => {
  const { success, error: showError } = useToast();
  const [isAddingRole, setIsAddingRole] = useState(false);

  // Fetch all available roles
  const { data: rolesData, loading: rolesLoading } = useQuery<RolesQuery>(ROLES_QUERY);

  // Assign role mutation
  const [assignRole, { loading: assignLoading }] = useMutation<
    AssignRoleMutation,
    AssignRoleMutationVariables
  >(ASSIGN_ROLE_MUTATION, {
    onCompleted: () => {
      success('Role assigned successfully', 'Success');
      setIsAddingRole(false);
      onRolesChanged?.();
    },
    onError: (error) => {
      showError(error.message || 'Failed to assign role', 'Error');
    },
  });

  // Remove role mutation
  const [removeRole, { loading: removeLoading }] = useMutation<
    RemoveRoleMutation,
    RemoveRoleMutationVariables
  >(REMOVE_ROLE_MUTATION, {
    onCompleted: () => {
      success('Role removed successfully', 'Success');
      onRolesChanged?.();
    },
    onError: (error) => {
      showError(error.message || 'Failed to remove role', 'Error');
    },
  });

  const handleAssignRole = async (roleId: string) => {
    await assignRole({
      variables: {
        userId,
        roleId,
      },
    });
  };

  const handleRemoveRole = async (roleId: string) => {
    await removeRole({
      variables: {
        userId,
        roleId,
      },
    });
  };

  const availableRoles = rolesData?.roles.filter(
    (role) => !currentRoles.some((cr) => cr.id === role.id)
  ) || [];

  const isLoading = assignLoading || removeLoading;

  return (
    <div className="space-y-4">
      {/* Current Roles */}
      <div>
        <h3 className="text-sm font-medium text-gray-300 mb-2">Current Roles</h3>
        <div className="flex flex-wrap gap-2">
          {currentRoles.length === 0 ? (
            <p className="text-sm text-gray-400">No roles assigned</p>
          ) : (
            currentRoles.map((role) => (
              <div
                key={role.id}
                className="flex items-center gap-2 px-3 py-1.5 bg-white/5 border border-white/10 rounded-lg"
              >
                <Badge variant="info" size="sm">
                  {role.name}
                </Badge>
                <button
                  onClick={() => handleRemoveRole(role.id)}
                  disabled={isLoading}
                  className="text-gray-400 hover:text-red-400 transition-colors disabled:opacity-50"
                  title="Remove role"
                >
                  <XMarkIcon className="h-4 w-4" />
                </button>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Add Role Section */}
      <div>
        {!isAddingRole ? (
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setIsAddingRole(true)}
            disabled={availableRoles.length === 0}
          >
            <PlusIcon className="h-4 w-4 mr-2" />
            Add Role
          </Button>
        ) : (
          <div className="space-y-2">
            <h3 className="text-sm font-medium text-gray-300">Available Roles</h3>
            {rolesLoading ? (
              <Loading size="sm" />
            ) : availableRoles.length === 0 ? (
              <p className="text-sm text-gray-400">No more roles available</p>
            ) : (
              <div className="flex flex-wrap gap-2">
                {availableRoles.map((role) => (
                  <button
                    key={role.id}
                    onClick={() => handleAssignRole(role.id)}
                    disabled={isLoading}
                    className="px-3 py-1.5 bg-white/5 border border-white/10 rounded-lg text-sm text-white hover:bg-white/10 transition-colors disabled:opacity-50"
                    title={role.description || undefined}
                  >
                    {role.name}
                  </button>
                ))}
              </div>
            )}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsAddingRole(false)}
            >
              Cancel
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};
