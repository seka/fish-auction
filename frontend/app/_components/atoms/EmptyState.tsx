import { Box, Text, Button, Stack } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { ReactNode } from 'react';

type EmptyStateProps = {
    message: string;
    icon?: ReactNode;
    action?: {
        label: string;
        onClick: () => void;
    };
};

export const EmptyState = ({ message, icon, action }: EmptyStateProps) => {
    return (
        <Box
            p="12"
            textAlign="center"
            bg="gray.50"
            borderRadius="lg"
            border="1px dashed"
            borderColor="gray.300"
            className={css({ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' })}
        >
            {icon && (
                <Box mb="4" fontSize="4xl" className={css({ color: 'gray.400' })}>
                    {icon}
                </Box>
            )}
            <Text className={css({ color: 'gray.600', fontWeight: 'medium' })} mb={action ? "6" : "0"}>
                {message}
            </Text>
            {action && (
                <Button
                    variant="primary"
                    onClick={action.onClick}
                    size="sm"
                >
                    {action.label}
                </Button>
            )}
        </Box>
    );
};
