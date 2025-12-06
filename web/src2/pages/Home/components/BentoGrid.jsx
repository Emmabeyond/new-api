import React from 'react';

const BentoGrid = ({ children, className = '' }) => {
  return (
    <div
      className={`
        grid
        grid-cols-1
        md:grid-cols-2
        lg:grid-cols-4
        gap-3 md:gap-4
        w-full
        max-w-6xl
        mx-auto
        px-4
        ${className}
      `}
    >
      {children}
    </div>
  );
};

export default BentoGrid;
