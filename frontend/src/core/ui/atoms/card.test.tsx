import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Card } from './card';

describe('Card', () => {
    it('renders with children', () => {
        render(<Card>Card Content</Card>);
        expect(screen.getByText('Card Content')).toBeInTheDocument();
    });

    it('renders with variants', () => {
        render(<Card variant="elevated">Elevated Card</Card>);
        expect(screen.getByText('Elevated Card')).toBeInTheDocument();
    });

    it('renders with padding variants', () => {
        render(<Card padding="lg">Padded Card</Card>);
        expect(screen.getByText('Padded Card')).toBeInTheDocument();
    });
});
