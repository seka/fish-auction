import { cva, type RecipeVariantProps } from 'styled-system/css';
import { styled } from 'styled-system/jsx';

const cardRecipe = cva({
    base: {
        display: 'block',
        bg: 'white',
        borderRadius: '1rem',
        overflow: 'hidden',
    },
    variants: {
        variant: {
            default: {
                border: '1px solid',
                borderColor: 'gray.100',
                shadow: 'xl',
            },
            elevated: {
                shadow: '2xl',
                transition: 'all 0.3s',
                _hover: {
                    shadow: '2xl',
                    transform: 'translateY(-0.25rem)',
                },
            },
            interactive: {
                border: '2px solid transparent',
                shadow: 'xl',
                transition: 'all 0.3s',
                cursor: 'pointer',
                _hover: {
                    shadow: '2xl',
                    transform: 'translateY(-0.25rem)',
                },
            },
        },
        padding: {
            none: {
                p: '0',
            },
            sm: {
                p: '4',
            },
            md: {
                p: '6',
            },
            lg: {
                p: '10',
            },
        },
    },
    defaultVariants: {
        variant: 'default',
        padding: 'md',
    },
});

export type CardProps = RecipeVariantProps<typeof cardRecipe> &
    React.HTMLAttributes<HTMLDivElement>;

export const Card = styled('div', cardRecipe);
