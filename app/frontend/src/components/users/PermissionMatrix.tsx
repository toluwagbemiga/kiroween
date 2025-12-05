'use client';

import React from 'react';
import { Badge } from '@/components/ui';
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline';

export interface PermissionMatrixProps {
  roles: Array<{
    id: string;
    name: string;
    permissions: string[];
  }>;
  userPermissions: string[];
}

export const PermissionMatrix: React.FC<PermissionMatrixProps> = ({
  roles,
  userPermissions,
}) => {
  // Get all unique permissions from roles
  const allPermissions = Array.from(
    new Set(roles.flatMap((role) => role.permissions))
  ).sort();

  if (allPermissions.length === 0) {
    return (
      <div className="text-sm text-gray-400">
        No permissions defined for these roles
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h3 className="text-sm font-medium text-gray-300">Permission Matrix</h3>
      
      <div className="overflow-x-auto rounded-lg border border-white/10 bg-white/5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-white/10">
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-300 uppercase">
                Permission
              </th>
              {roles.map((role) => (
                <th
                  key={role.id}
                  className="px-4 py-3 text-center text-xs font-medium text-gray-300 uppercase"
                >
                  <Badge variant="info" size="sm">
                    {role.name}
                  </Badge>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {allPermissions.map((permission, index) => {
              const hasPermission = userPermissions.includes(permission);
              
              return (
                <tr
                  key={permission}
                  className={index % 2 === 0 ? 'bg-white/5' : ''}
                >
                  <td className="px-4 py-3 text-white">
                    <div className="flex items-center gap-2">
                      {hasPermission ? (
                        <CheckCircleIcon className="h-4 w-4 text-green-400" />
                      ) : (
                        <XCircleIcon className="h-4 w-4 text-gray-500" />
                      )}
                      <span className="font-mono text-xs">{permission}</span>
                    </div>
                  </td>
                  {roles.map((role) => {
                    const roleHasPermission = role.permissions.includes(permission);
                    
                    return (
                      <td
                        key={`${role.id}-${permission}`}
                        className="px-4 py-3 text-center"
                      >
                        {roleHasPermission ? (
                          <CheckCircleIcon className="h-5 w-5 text-green-400 mx-auto" />
                        ) : (
                          <XCircleIcon className="h-5 w-5 text-gray-600 mx-auto" />
                        )}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Summary */}
      <div className="text-xs text-gray-400">
        <p>
          Total permissions: {userPermissions.length} of {allPermissions.length}
        </p>
      </div>
    </div>
  );
};
