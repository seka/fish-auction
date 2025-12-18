import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Text } from './text';

describe('Text', () => {
    it('renders correctly with default variant', () => {
        render(<Text>Default Text</Text>);
        const text = screen.getByText('Default Text');
        expect(text).toBeInTheDocument();
        expect(text.tagName).toBe('P');
    });

    it('renders as h1 when "as" prop is provided', () => {
        render(<Text as="h1">Heading 1</Text>);
        const text = screen.getByRole('heading', { level: 1 });
        expect(text).toBeInTheDocument();
        expect(text.tagName).toBe('H1');
    });

    it('applies variant styles correctly', () => {
        render(<Text variant="h2" as="h2">Heading 2</Text>);
        const text = screen.getByRole('heading', { level: 2 });
        expect(text).toBeInTheDocument();
    });

    it('applies color variant correctly', () => {
        render(<Text color="primary">Primary Text</Text>);
        const text = screen.getByText('Primary Text');
        expect(text).toBeInTheDocument();
    });
});
