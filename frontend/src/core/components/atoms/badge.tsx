import { ReactNode } from 'react';
import { Box } from './index';
import { css, cx } from 'styled-system/css';

interface BadgeProps {
  children: ReactNode;
  variant?: 'success' | 'warning' | 'error' | 'info' | 'neutral';
  className?: string;
}

export const Badge = ({ children, variant = 'neutral', className }: BadgeProps) => {
  const variantStyles = {
    success: css({
      bg: 'green.50',
      color: 'green.700',
      borderColor: 'green.200',
    }),
    warning: css({
      bg: 'yellow.50',
      color: 'yellow.700',
      borderColor: 'yellow.200',
    }),
    error: css({
      bg: 'red.50',
      color: 'red.700',
      borderColor: 'red.200',
    }),
    info: css({
      bg: 'blue.50',
      color: 'blue.700',
      borderColor: 'blue.200',
    }),
    neutral: css({
      bg: 'gray.50',
      color: 'gray.700',
      borderColor: 'gray.200',
    }),
  };

  return (
    <Box
      as="span"
      px="2.5"
      py="0.5"
      borderRadius="full"
      fontSize="xs"
      fontWeight="medium"
      border="1px solid"
      className={cx(variantStyles[variant], className)}
    >
      {children}
    </Box>
  );
};
