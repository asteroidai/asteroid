import React from 'react';
import { format } from 'timeago.js';

interface CreatedAgoProps {
  datetime: string;
  label?: string;
  className?: string;
}

export const CreatedAgo: React.FC<CreatedAgoProps> = ({ datetime, label = 'Created', className = '' }) => {
  const formattedDate = format(datetime);

  return (
    <time
      dateTime={datetime}
      className={`text-sm text-gray-500 ${className}`}
      title={new Date(datetime).toLocaleString()}
    >
      {label} {formattedDate}
    </time>
  );
};
