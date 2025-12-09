import { styled } from '@/styled-system/jsx';

export const Table = styled('table', {
    base: {
        minW: 'full',
        divideY: '1px',
        divideColor: 'gray.200'
    }
});

export const Thead = styled('thead', {
    base: {
        bg: 'gray.50'
    }
});

export const Tbody = styled('tbody', {
    base: {
        bg: 'white',
        divideY: '1px',
        divideColor: 'gray.200'
    }
});

export const Tr = styled('tr', {
    base: {
        cursor: 'pointer',
        _hover: {
            bg: 'gray.50',
            transition: 'colors'
        }
    }
});

export const Th = styled('th', {
    base: {
        px: '6',
        py: '4',
        textAlign: 'left',
        fontSize: 'xs',
        fontWeight: 'bold',
        color: 'gray.500',
        textTransform: 'uppercase',
        letterSpacing: 'wider'
    }
});

export const Td = styled('td', {
    base: {
        px: '6',
        py: '4',
        whiteSpace: 'nowrap'
    }
});
