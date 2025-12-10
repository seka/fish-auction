'use client';

import { useBuyerPage } from './_hooks/useBuyerPage';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { css } from 'styled-system/css';

export default function AdminBuyersPage() {
    const { state, form, actions, t } = useBuyerPage();

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                {t('Admin.Buyers.title')}
            </Text>

            {state.message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" className={css({ color: 'blue.700' })} p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">{t('Common.notification')}</Text>
                    <Text>{state.message}</Text>
                </Box>
            )}

            <Box display="grid" gridTemplateColumns={{ base: '1fr', md: '3fr 1fr' }} gap="8" className={css({ md: { gridTemplateColumns: '1fr 2fr' } })}>
                {/* Form Section */}
                <Box className={css({ md: { gridColumn: '1 / 2' } })}>
                    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
                        <HStack mb="6">
                            <Box w="2" h="6" bg="green.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'green.900' })} fontWeight="bold">
                                {t('Admin.Buyers.register_title')}
                            </Text>
                        </HStack>
                        <form onSubmit={actions.onSubmit}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Buyers.name')}
                                    </Text>
                                    <Input
                                        type="text"
                                        {...form.register('name')}
                                        placeholder={t('Admin.Buyers.placeholder_name')}
                                        error={!!form.errors.name}
                                        className={css({ _focus: { borderColor: 'green.500', ringColor: 'green.500' } })}
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
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">{t('Admin.Buyers.list_title')}</Text>
                        </Box>
                        {state.isLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.loading)}</Box>
                        ) : state.buyers.length === 0 ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.no_data)}</Box>
                        ) : (
                            <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
                                {state.buyers.map((buyer) => (
                                    <Box as="li" key={buyer.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
                                        <HStack justify="between" align="center">
                                            <Box>
                                                <Text as="h3" fontSize="lg" fontWeight="bold" className={css({ color: 'green.900' })}>{buyer.name}</Text>
                                                <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="1">ID: {buyer.id}</Text>
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
