import React from 'react';
import { cn } from '@/lib/utils';

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
}

const Badge = React.forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant = 'default', size = 'md', children, ...props }, ref) => {
    const baseStyles = [
      'inline-flex items-center justify-center rounded-full font-medium',
      'backdrop-blur-sm border transition-all duration-200',
    ];

    const variants = {
      default: [
        'bg-white/10 text-white border-white/20',
      ],
      success: [
        'bg-green-500/20 text-green-300 border-green-500/30',
      ],
      warning: [
        'bg-yellow-500/20 text-yellow-300 border-yellow-500/30',
      ],
      error: [
        'bg-red-500/20 text-red-300 border-red-500/30',
      ],
      info: [
        'bg-blue-500/20 text-blue-300 border-blue-500/30',
      ],
      outline: [
        'bg-transparent text-white border-white/40',
      ],
    };

    const sizes = {
      sm: 'px-2 py-0.5 text-xs',
      md: 'px-2.5 py-1 text-sm',
      lg: 'px-3 py-1.5 text-base',
    };

    return (
      <div
        ref={ref}
        className={cn(
          baseStyles,
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

Badge.displayName = 'Badge';

// Status Badge Component for common statuses
export interface StatusBadgeProps extends Omit<BadgeProps, 'variant'> {
  status: 'active' | 'inactive' | 'pending' | 'completed' | 'failed';
}

const StatusBadge = React.forwardRef<HTMLDivElement, StatusBadgeProps>(
  ({ status, ...props }, ref) => {
    const statusMap = {
      active: { variant: 'success' as const, label: 'Active' },
      inactive: { variant: 'default' as const, label: 'Inactive' },
      pending: { variant: 'warning' as const, label: 'Pending' },
      completed: { variant: 'success' as const, label: 'Completed' },
      failed: { variant: 'error' as const, label: 'Failed' },
    };

    const { variant, label } = statusMap[status];

    return (
      <Badge ref={ref} variant={variant} {...props}>
        {label}
      </Badge>
    );
  }
);

StatusBadge.displayName = 'StatusBadge';

export { Badge, StatusBadge };
