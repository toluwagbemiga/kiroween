'use client';

import React from 'react';
import { Button } from '@/components/ui';
import { subDays, startOfDay, endOfDay } from 'date-fns';

export interface DateRange {
  startDate: Date;
  endDate: Date;
}

export interface DateRangePickerProps {
  value: DateRange;
  onChange: (range: DateRange) => void;
}

const presets = [
  { label: 'Last 7 days', days: 7 },
  { label: 'Last 30 days', days: 30 },
  { label: 'Last 90 days', days: 90 },
];

export const DateRangePicker: React.FC<DateRangePickerProps> = ({
  value,
  onChange,
}) => {
  const handlePresetClick = (days: number) => {
    const endDate = endOfDay(new Date());
    const startDate = startOfDay(subDays(endDate, days));
    onChange({ startDate, endDate });
  };

  return (
    <div className="flex items-center gap-2">
      {presets.map((preset) => (
        <Button
          key={preset.days}
          variant="secondary"
          size="sm"
          onClick={() => handlePresetClick(preset.days)}
        >
          {preset.label}
        </Button>
      ))}
    </div>
  );
};
