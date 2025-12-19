'use client';

import { useVenuePage } from './_hooks/useVenuePage';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { css } from 'styled-system/css';
import { styled } from 'styled-system/jsx';
import { EmptyState } from '../../_components/atoms/EmptyState';

// Textarea component with similar styling to Input
const Textarea = styled('textarea', {
    base: {
        display: 'block',
        width: 'full',
        px: '3',
        py: '2',
        bg: 'white',
        color: 'gray.900',
        border: '1px solid',
        borderColor: 'gray.300',
        borderRadius: 'md',
        fontSize: 'sm',
        outline: 'none',
        transition: 'border-color 0.2s',
        resize: 'vertical',
        _focus: {
            borderColor: 'indigo.500',
            ring: '1px',
            ringColor: 'indigo.500',
        },
        _disabled: {
            bg: 'gray.50',
            cursor: 'not-allowed',
        },
    }
});

export default function AdminVenuesPage() {
    const { state, form, actions, t } = useVenuePage();

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                {t('Admin.Venues.title')}
            </Text>

            {state.message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" color="blue.700" p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">{t('Common.notification')}</Text>
                    <Text>{state.message}</Text>
                </Box>
            )}

            <Box display="grid" gridTemplateColumns={{ base: '1fr', md: '3fr 1fr' }} gap="8" className={css({ md: { gridTemplateColumns: '1fr 2fr' } })}>
                {/* Form Section */}
                <Box className={css({ md: { gridColumn: '1 / 2' } })}>
                    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
                        <HStack mb="6">
                            <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'indigo.900' })} fontWeight="bold">
                                {state.editingVenue ? t('Admin.Venues.edit_title') : t('Admin.Venues.register_title')}
                            </Text>
                        </HStack>
                        <form onSubmit={actions.onSubmit}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Venues.name')}
                                    </Text>
                                    <Input
                                        type="text"
                                        {...form.register('name')}
                                        placeholder={t('Admin.Venues.placeholder_name')}
                                        error={!!form.errors.name}
                                    />
                                    {form.errors.name && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.name.message}</Text>
                                    )}
                                </Box>
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Venues.location')}
                                    </Text>
                                    <Input
                                        type="text"
                                        {...form.register('location')}
                                        placeholder={t('Admin.Venues.placeholder_location')}
                                    />
                                </Box>
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Venues.description')}
                                    </Text>
                                    <Textarea
                                        {...form.register('description')}
                                        rows={3}
                                        placeholder={t('Admin.Venues.placeholder_description')}
                                    />
                                </Box>

                                <HStack spacing="2">
                                    <Button
                                        type="submit"
                                        disabled={state.isCreating || state.isUpdating}
                                        width="full"
                                        className={css({ flex: '1' })}
                                        variant="primary"
                                    >
                                        {state.editingVenue ? (state.isUpdating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.update)) : (state.isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register))}
                                    </Button>
                                    {state.editingVenue && (
                                        <Button
                                            type="button"
                                            onClick={actions.onCancelEdit}
                                            variant="outline"
                                        >
                                            {t(COMMON_TEXT_KEYS.cancel)}
                                        </Button>
                                    )}
                                </HStack>
                            </Stack>
                        </form>
                    </Card>
                </Box>

                {/* List Section */}
                <Box className={css({ md: { gridColumn: '2 / 3' } })}>
                    <Card padding="none" overflow="hidden">
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">{t('Admin.Venues.list_title')}</Text>
                        </Box>
                        {state.isLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.loading)}</Box>
                        ) : state.venues.length === 0 ? (
                            <EmptyState
                                message={t(COMMON_TEXT_KEYS.no_data)}
                                icon={<span role="img" aria-label="venue">📍</span>}
                            />
                        ) : (
                            <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
                                {state.venues.map((venue) => (
                                    <Box as="li" key={venue.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
                                        <HStack justify="between" align="start">
                                            <Box>
                                                <Text as="h3" fontSize="lg" fontWeight="bold" className={css({ color: 'indigo.900' })}>{venue.name}</Text>
                                                {venue.location && (
                                                    <Text fontSize="sm" className={css({ color: 'gray.700' })} mt="1" display="flex" alignItems="center">
                                                        <span className={css({ mr: '2' })}>📍</span>
                                                        {venue.location}
                                                    </Text>
                                                )}
                                                {venue.description && (
                                                    <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="2">{venue.description}</Text>
                                                )}
                                            </Box>
                                            <HStack spacing="2">
                                                <Button size="sm" variant="outline" onClick={() => actions.onEdit(venue)}>
                                                    {t(COMMON_TEXT_KEYS.edit)}
                                                </Button>
                                                <Button size="sm" className={css({ bg: 'red.50', color: 'red.600', _hover: { bg: 'red.100' } })} onClick={() => actions.onDelete(venue.id)}>
                                                    {t(COMMON_TEXT_KEYS.delete)}
                                                </Button>
                                            </HStack>
                                        </HStack>
                                    </Box>
                                ))}
                            </Stack>
                        )}
                    </Card>
                </Box>
            </Box>
        </Box>
    );
}
