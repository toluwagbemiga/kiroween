'use client';

import { useQuery, useMutation } from '@apollo/client';
import { gql } from '@apollo/client';
import { DashboardLayout, PageContainer } from '@/components/layout';
import { Card, CardHeader, CardTitle, CardContent, Badge } from '@/components/ui';
import { withAuth } from '@/components/auth';

const GET_ADMIN_DATA = gql`
  query GetAdminData {
    users(limit: 100, offset: 0) {
      nodes {
        id
        email
        name
        teamId
        roles {
          name
        }
        createdAt
      }
      totalCount
    }
    roles {
      id
      name
      description
      permissions
      isSystem
    }
  }
`;

const ASSIGN_ROLE = gql`
  mutation AssignRole($userId: ID!, $roleId: ID!) {
    assignRole(userId: $userId, roleId: $roleId) {
      id
      email
    }
  }
`;

function AdminDashboard() {
  const { data, loading, refetch } = useQuery(GET_ADMIN_DATA);
  const [assignRole] = useMutation(ASSIGN_ROLE);

  const handleAssignRole = async (userId: string, roleId: string) => {
    try {
      await assignRole({ variables: { userId, roleId } });
      refetch();
    } catch (error) {
      console.error('Error assigning role:', error);
    }
  };

  if (loading) {
    return (
      <DashboardLayout title="Admin Dashboard">
        <PageContainer>
          <div className="text-center py-12 text-gray-400">Loading admin data...</div>
        </PageContainer>
      </DashboardLayout>
    );
  }

  const users = data?.users?.nodes || [];
  const roles = data?.roles || [];
  const totalUsers = data?.users?.totalCount || 0;

  return (
    <DashboardLayout
      title="Admin Dashboard"
      description="Manage users, roles, and system settings"
    >
      <PageContainer>
        <div className="space-y-6">
          {/* Stats */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <Card variant="glass">
              <CardContent className="p-6">
                <h3 className="text-sm font-medium text-gray-400">Total Users</h3>
                <p className="text-3xl font-bold text-white mt-2">{totalUsers}</p>
              </CardContent>
            </Card>
            <Card variant="glass">
              <CardContent className="p-6">
                <h3 className="text-sm font-medium text-gray-400">Total Roles</h3>
                <p className="text-3xl font-bold text-white mt-2">{roles.length}</p>
              </CardContent>
            </Card>
            <Card variant="glass">
              <CardContent className="p-6">
                <h3 className="text-sm font-medium text-gray-400">System Status</h3>
                <Badge variant="success" className="mt-2">Operational</Badge>
              </CardContent>
            </Card>
          </div>

          {/* Users Table */}
          <Card variant="glass">
            <CardHeader>
              <CardTitle>All Users</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-700">
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">Email</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">Name</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">Roles</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">Team ID</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">Created</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map((user: any) => (
                      <tr key={user.id} className="border-b border-gray-800 hover:bg-gray-800/50">
                        <td className="py-3 px-4 text-sm text-white">{user.email}</td>
                        <td className="py-3 px-4 text-sm text-gray-300">{user.name}</td>
                        <td className="py-3 px-4 text-sm">
                          <div className="flex gap-1 flex-wrap">
                            {user.roles?.map((role: any) => (
                              <Badge key={role.name} variant="default" size="sm">
                                {role.name}
                              </Badge>
                            ))}
                          </div>
                        </td>
                        <td className="py-3 px-4 text-sm text-gray-400">
                          {user.teamId || 'No team'}
                        </td>
                        <td className="py-3 px-4 text-sm text-gray-400">
                          {new Date(user.createdAt).toLocaleDateString()}
                        </td>
                        <td className="py-3 px-4 text-sm">
                          <select
                            onChange={(e) => handleAssignRole(user.id, e.target.value)}
                            className="bg-gray-800 text-white text-xs rounded px-2 py-1 border border-gray-700"
                            defaultValue=""
                          >
                            <option value="" disabled>Assign Role</option>
                            {roles.map((role: any) => (
                              <option key={role.id} value={role.id}>
                                {role.name}
                              </option>
                            ))}
                          </select>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>

          {/* Roles */}
          <Card variant="glass">
            <CardHeader>
              <CardTitle>System Roles</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                {roles.map((role: any) => (
                  <div key={role.id} className="border border-gray-700 rounded-lg p-4">
                    <div className="flex items-center justify-between mb-2">
                      <h3 className="font-semibold text-white">{role.name}</h3>
                      {role.isSystem && (
                        <Badge variant="info" size="sm">System</Badge>
                      )}
                    </div>
                    <p className="text-sm text-gray-400 mb-3">{role.description}</p>
                    <div className="flex flex-wrap gap-1">
                      {role.permissions?.slice(0, 5).map((perm: string) => (
                        <Badge key={perm} variant="default" size="sm">
                          {perm}
                        </Badge>
                      ))}
                      {role.permissions?.length > 5 && (
                        <Badge variant="outline" size="sm">
                          +{role.permissions.length - 5} more
                        </Badge>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </PageContainer>
    </DashboardLayout>
  );
}

export default withAuth(AdminDashboard);
