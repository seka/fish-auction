import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Stack, HStack, Box } from './stack';

describe('Stack Atoms', () => {
    describe('Box', () => {
        it('renders correctly', () => {
            render(<Box>Box Content</Box>);
            expect(screen.getByText('Box Content')).toBeInTheDocument();
        });
    });

    describe('Stack', () => {
        it('renders children correctly', () => {
            render(
                <Stack>
                    <div>Item 1</div>
                    <div>Item 2</div>
                </Stack>
            );
            expect(screen.getByText('Item 1')).toBeInTheDocument();
            expect(screen.getByText('Item 2')).toBeInTheDocument();
        });
    });

    describe('HStack', () => {
        it('renders children correctly', () => {
            render(
                <HStack>
                    <div>Item A</div>
                    <div>Item B</div>
                </HStack>
            );
            expect(screen.getByText('Item A')).toBeInTheDocument();
            expect(screen.getByText('Item B')).toBeInTheDocument();
        });
    });
});
