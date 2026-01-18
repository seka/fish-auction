'use client';

import { Box, Stack, HStack, Text } from '@/src/core/ui';
import { useToast } from '@/src/hooks/useToast';
import { css } from '@/styled-system/css';
import { useRouter } from 'next/navigation';

export const ToastContainer = () => {
    const { toasts, removeToast } = useToast();
    console.log('ToastContainer rendering. Toasts count:', toasts.length);

    return (
        <Box
            position="fixed"
            top="8"
            right="8"
            zIndex="2147483647"
            display="flex"
            flexDirection="column"
            gap="2"
            pointerEvents="none"
            style={{ opacity: 1, visibility: 'visible' }}
        >
            {toasts.map((toast) => (
                <ToastItem key={toast.id} toast={toast} onRemove={() => removeToast(toast.id)} />
            ))}
        </Box>
    );
};

const ToastItem = ({ toast, onRemove }: { toast: any; onRemove: () => void }) => {
    const router = useRouter();

    const handleClick = () => {
        if (toast.url) {
            router.push(toast.url);
        }
        onRemove();
    };

    return (
        <Box
            pointerEvents="auto"
            bg="white"
            boxShadow="lg"
            rounded="md"
            overflow="hidden"
            maxW="sm"
            w="full"
            borderLeftWidth="4px"
            borderLeftColor={
                toast.type === 'error' ? 'red.500' :
                    toast.type === 'warning' ? 'yellow.500' :
                        toast.type === 'success' ? 'green.500' :
                            'blue.500'
            }
            className={css({
                animation: 'slideIn 0.3s ease-out forwards',
                cursor: toast.url ? 'pointer' : 'default',
                _hover: {
                    transform: toast.url ? 'translateY(-2px)' : 'none',
                    boxShadow: toast.url ? 'xl' : 'lg',
                }
            })}
            onClick={handleClick}
        >
            <HStack p="4" gap="4" alignItems="flex-start">
                <Box flex="1">
                    <Text fontWeight="bold" fontSize="sm" color="default">
                        {toast.title}
                    </Text>
                    <Text fontSize="sm" color="muted" mt="1">
                        {toast.body}
                    </Text>
                </Box>
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        onRemove();
                    }}
                    className={css({
                        color: 'gray.400',
                        _hover: { color: 'gray.600' },
                        fontSize: 'xl',
                        lineHeight: '1',
                        cursor: 'pointer'
                    })}
                >
                    &times;
                </button>
            </HStack>
        </Box>
    );
};
