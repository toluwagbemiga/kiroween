'use client';

import React from 'react';
import { Badge } from '@/components/ui';

export interface EventTableProps {
  events: Array<{
    eventName: string;
    count: number;
  }>;
}

export const EventTable: React.FC<EventTableProps> = ({ events }) => {
  if (events.length === 0) {
    return (
      <div className="text-center py-8 text-gray-400">
        No events recorded yet
      </div>
    );
  }

  return (
    <div className="overflow-x-auto rounded-lg border border-white/10 bg-white/5">
      <table className="w-full">
        <thead>
          <tr className="border-b border-white/10">
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase">
              Event Name
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-300 uppercase">
              Count
            </th>
          </tr>
        </thead>
        <tbody>
          {events.map((event, index) => (
            <tr
              key={event.eventName}
              className={`border-b border-white/5 ${
                index % 2 === 0 ? 'bg-white/5' : ''
              }`}
            >
              <td className="px-6 py-4">
                <Badge variant="info" size="sm">
                  {event.eventName}
                </Badge>
              </td>
              <td className="px-6 py-4 text-right text-white font-medium">
                {event.count.toLocaleString()}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
