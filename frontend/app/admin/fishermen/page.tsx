'use client';

import { useFishermanPage } from './_hooks/useFishermanPage';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { EmptyState } from '../../_components/atoms/EmptyState';

export default function AdminFishermenPage() {
    const { state, form, actions, t } = useFishermanPage();

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                {t('Admin.Fishermen.title')}
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
                                {t('Admin.Fishermen.register_title')}
                            </Text>
                        </HStack>
                        <form onSubmit={actions.onSubmit}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Fishermen.name')}
                                    </Text>
                                    <Input
                                        type="text"
                                        {...form.register('name')}
                                        placeholder={t('Admin.Fishermen.placeholder_name')}
                                        error={!!form.errors.name}
                                    />
                                    {form.errors.name && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.name.message}</Text>
                                    )}
                                </Box>

                                <Button
                                    type="submit"
                                    disabled={state.isCreating}
                                    width="full"
                                    className={css({ flex: '1' })}
                                    variant="primary"
                                >
                                    {state.isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register)}
                                </Button>
                            </Stack>
                        </form>
                    </Card>
                </Box>

                {/* List Section */}
                <Box className={css({ md: { gridColumn: '2 / 3' } })}>
                    <Card padding="none" overflow="hidden">
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">{t('Admin.Fishermen.list_title')}</Text>
                        </Box>
                        {state.isLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.loading)}</Box>
                        ) : state.fishermen.length === 0 ? (
                            <EmptyState
                                message={t(COMMON_TEXT_KEYS.no_data)}
                                icon={<span role="img" aria-label="fisherman">🎣</span>}
                            />
                        ) : (
                            <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
                                {state.fishermen.map((fisherman) => (
                                    <Box as="li" key={fisherman.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
                                        <HStack justify="between" align="center">
                                            <Box>
                                                <Text as="h3" fontSize="lg" fontWeight="bold" className={css({ color: 'indigo.900' })}>{fisherman.name}</Text>
                                                <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="1">ID: {fisherman.id}</Text>
                                            </Box>
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
