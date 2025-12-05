'use client';

import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { UserTable } from '@/components/users/UserTable';
import { UserFormModal } from '@/components/users/UserFormModal';
import { RoleManagementModal } from '@/components/users/RoleManagementModal';
import { Button, Loading } from '@/components/ui';
import { EmptyState } from '@/components/layout/EmptyState';
import { USERS_LIST_QUERY } from '@/lib/graphql/queries/users';
import { UsersListQuery, UsersListQueryVariables } from '@/lib/graphql/generated/graphql';
import { UserPlusIcon } from '@heroicons/react/24/outline';

export default function UsersPage() {
  const [page, setPage] = useState(0);
  const pageSize = 20;
  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const [isRoleModalOpen, setIsRoleModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<any>(null);

  const { data, loading, error, refetch } = useQuery<UsersListQuery, UsersListQueryVariables>(
    USERS_LIST_QUERY,
    {
      variables: {
        limit: pageSize,
        offset: page * pageSize,
      },
      fetchPolicy: 'cache-and-network',
    }
  );

  const handleEditUser = (user: any) => {
    setSelectedUser(user);
    setIsFormModalOpen(true);
  };

  const handleManageRoles = (user: any) => {
    setSelectedUser(user);
    setIsRoleModalOpen(true);
  };

  const handleCreateUser = () => {
    setSelectedUser(null);
    setIsFormModalOpen(true);
  };

  const handleCloseFormModal = () => {
    setIsFormModalOpen(false);
    setSelectedUser(null);
  };

  const handleCloseRoleModal = () => {
    setIsRoleModalOpen(false);
    setSelectedUser(null);
  };

  if (error) {
    return (
      <DashboardLayout
        title="Users"
        description="Manage user accounts and permissions"
      >
        <EmptyState
          title="Error loading users"
          description={error.message}
          action={{
            label: 'Try Again',
            onClick: () => refetch(),
          }}
        />
      </DashboardLayout>
    );
  }

  const users = data?.users.nodes || [];
  const totalCount = data?.users.totalCount || 0;
  const hasUsers = users.length > 0;

  return (
    <DashboardLayout
      title="Users"
      description="Manage user accounts and permissions"
      actions={
        <Button
          variant="primary"
          onClick={handleCreateUser}
        >
          <UserPlusIcon className="h-5 w-5 mr-2" />
          Add User
        </Button>
      }
    >
      {loading && !data ? (
        <div className="flex items-center justify-center py-12">
          <Loading size="lg" />
        </div>
      ) : !hasUsers ? (
        <EmptyState
          icon={<UserPlusIcon className="h-12 w-12" />}
          title="No users yet"
          description="Get started by creating your first user account"
          action={{
            label: 'Add User',
            onClick: handleCreateUser,
          }}
        />
      ) : (
        <div className="space-y-6">
          {/* Stats */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-white/5 backdrop-blur-xl rounded-lg border border-white/10 p-6">
              <div className="text-sm text-gray-400">Total Users</div>
              <div className="text-3xl font-bold text-white mt-2">
                {totalCount}
              </div>
            </div>
            <div className="bg-white/5 backdrop-blur-xl rounded-lg border border-white/10 p-6">
              <div className="text-sm text-gray-400">Active Users</div>
              <div className="text-3xl font-bold text-white mt-2">
                {users.length}
              </div>
            </div>
            <div className="bg-white/5 backdrop-blur-xl rounded-lg border border-white/10 p-6">
              <div className="text-sm text-gray-400">Roles Assigned</div>
              <div className="text-3xl font-bold text-white mt-2">
                {users.reduce((acc, user) => acc + user.roles.length, 0)}
              </div>
            </div>
          </div>

          {/* User Table */}
          <UserTable
            users={users}
            loading={loading}
            onEditUser={handleEditUser}
            onManageRoles={handleManageRoles}
          />

          {/* Pagination */}
          {totalCount > pageSize && (
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-400">
                Showing {page * pageSize + 1} to{' '}
                {Math.min((page + 1) * pageSize, totalCount)} of {totalCount}{' '}
                users
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setPage((p) => Math.max(0, p - 1))}
                  disabled={page === 0}
                >
                  Previous
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setPage((p) => p + 1)}
                  disabled={(page + 1) * pageSize >= totalCount}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* User Form Modal */}
      <UserFormModal
        isOpen={isFormModalOpen}
        onClose={handleCloseFormModal}
        user={selectedUser}
      />

      {/* Role Management Modal */}
      <RoleManagementModal
        isOpen={isRoleModalOpen}
        onClose={handleCloseRoleModal}
        userId={selectedUser?.id}
        userName={selectedUser?.name}
      />
    </DashboardLayout>
  );
}
