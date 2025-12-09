import { styled } from '@/styled-system/jsx';

// Modal backdrop and container
export const ModalBackdrop = styled('div', {
    base: {
        position: 'fixed',
        inset: '0',
        bg: 'blackAlpha.600',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: '50',
        p: '4',
    }
});

export const ModalContent = styled('div', {
    base: {
        bg: 'white',
        borderRadius: 'xl',
        shadow: '2xl',
        maxW: '2xl',
        w: 'full',
        maxH: '90vh',
        overflowY: 'auto',
    }
});
