import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ModalBackdrop, ModalContent } from './modal';

describe('Modal Molecules', () => {
    it('renders correctly', () => {
        render(
            <ModalBackdrop data-testid="backdrop">
                <ModalContent>Modal Content</ModalContent>
            </ModalBackdrop>
        );
        expect(screen.getByTestId('backdrop')).toBeInTheDocument();
        expect(screen.getByText('Modal Content')).toBeInTheDocument();
    });
});
