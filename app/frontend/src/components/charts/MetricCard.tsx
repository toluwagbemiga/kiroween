'use client';

import React from 'react';
import { Card, CardContent, Badge } from '@/components/ui';
import { ArrowUpIcon, ArrowDownIcon } from '@heroicons/react/24/outline';

interface MetricCardProps {
  title: string;
  value: string | number;
  change?: string;
  trend?: 'up' | 'down' | 'neutral';
  icon?: React.ComponentType<{ className?: string }>;
  description?: string;
  loading?: boolean;
}

/**
 * MetricCard component for displaying KPI metrics
 */
export const MetricCard: React.FC<MetricCardProps> = ({
  title,
  value,
  change,
  trend = 'neutral',
  icon: Icon,
  description,
  loading = false,
}) => {
  const getTrendColor = () => {
    switch (trend) {
      case 'up':
        return 'text-green-400';
      case 'down':
        return 'text-red-400';
      default:
        return 'text-gray-400';
    }
  };

  const getTrendIcon = () => {
    if (trend === 'up') return ArrowUpIcon;
    if (trend === 'down') return ArrowDownIcon;
    return null;
  };

  const TrendIcon = getTrendIcon();

  if (loading) {
    return (
      <Card variant="glass">
        <CardContent className="p-6">
          <div className="animate-pulse space-y-3">
            <div className="h-4 bg-gray-700 rounded w-1/2"></div>
            <div className="h-8 bg-gray-700 rounded w-3/4"></div>
            <div className="h-4 bg-gray-700 rounded w-1/3"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card variant="glass" hover>
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <p className="text-sm font-medium text-gray-300 mb-2">{title}</p>
            <p className="text-3xl font-bold text-white mb-2">
              {typeof value === 'number' ? value.toLocaleString() : value}
            </p>
            
            {change && (
              <div className="flex items-center space-x-1">
                {TrendIcon && (
                  <TrendIcon className={`h-4 w-4 ${getTrendColor()}`} />
                )}
                <span className={`text-sm font-medium ${getTrendColor()}`}>
                  {change}
                </span>
              </div>
            )}
            
            {description && (
              <p className="text-xs text-gray-400 mt-2">{description}</p>
            )}
          </div>
          
          {Icon && (
            <div className="p-3 rounded-lg bg-primary-500/20">
              <Icon className="h-6 w-6 text-primary-400" />
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
