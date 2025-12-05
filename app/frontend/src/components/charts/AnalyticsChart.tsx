'use client';

import React from 'react';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface AnalyticsChartProps {
  data: Array<{
    name: string;
    value: number;
    [key: string]: any;
  }>;
  type?: 'line' | 'bar';
  dataKey?: string;
  xAxisKey?: string;
  height?: number;
}

/**
 * Custom tooltip for charts
 */
const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-gray-800/95 backdrop-blur-sm border border-gray-700 rounded-lg p-3 shadow-lg">
        <p className="text-sm font-medium text-white mb-1">{label}</p>
        {payload.map((entry: any, index: number) => (
          <p key={index} className="text-xs text-gray-300">
            <span style={{ color: entry.color }}>{entry.name}: </span>
            <span className="font-semibold">{entry.value?.toLocaleString()}</span>
          </p>
        ))}
      </div>
    );
  }
  return null;
};

/**
 * Analytics Chart component for visualizing time-series data
 */
export const AnalyticsChart: React.FC<AnalyticsChartProps> = ({
  data,
  type = 'line',
  dataKey = 'value',
  xAxisKey = 'name',
  height = 300,
}) => {
  const ChartComponent = type === 'line' ? LineChart : BarChart;
  const DataComponent = type === 'line' ? Line : Bar;

  return (
    <ResponsiveContainer width="100%" height={height}>
      <ChartComponent
        data={data}
        margin={{ top: 5, right: 20, left: 0, bottom: 5 }}
      >
        <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
        <XAxis
          dataKey={xAxisKey}
          stroke="#9CA3AF"
          style={{ fontSize: '12px' }}
          tick={{ fill: '#9CA3AF' }}
        />
        <YAxis
          stroke="#9CA3AF"
          style={{ fontSize: '12px' }}
          tick={{ fill: '#9CA3AF' }}
        />
        <Tooltip content={<CustomTooltip />} />
        <Legend
          wrapperStyle={{ fontSize: '12px', color: '#9CA3AF' }}
          iconType="circle"
        />
        <DataComponent
          type={type === 'line' ? 'monotone' : undefined}
          dataKey={dataKey}
          stroke="#8B5CF6"
          fill="#8B5CF6"
          strokeWidth={2}
          name="Events"
        />
      </ChartComponent>
    </ResponsiveContainer>
  );
};
