'use client';

import { useItemPage } from './_hooks/useItemPage';
import { Box, Stack, HStack, Text, Card, Button, Input, Select } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';

export default function AdminItemsPage() {
    const { state, form, actions, t } = useItemPage();

    return (
        <Box maxW="6xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                {t('Admin.Items.title')}
            </Text>

            {state.message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" className={css({ color: 'blue.700' })} p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">{t('Common.notification')}</Text>
                    <Text>{state.message}</Text>
                </Box>
            )}

            <Card padding="lg">
                <HStack mb="6">
                    <Box w="2" h="6" bg="orange.500" mr="3" borderRadius="full" />
                    <Text as="h2" variant="h4" className={css({ color: 'orange.900' })} fontWeight="bold">
                        {t('Admin.Items.register_title')}
                    </Text>
                </HStack>
                <form onSubmit={actions.onSubmit}>
                    <Box display="grid" gridTemplateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap="6">
                        <Box className={css({ md: { gridColumn: 'span 2' } })}>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                {t('Admin.Items.auction')}
                            </Text>
                            <Select
                                {...form.register('auctionId')}
                            >
                                <option value="">{t('Admin.Items.placeholder_select_auction')}</option>
                                {state.auctions.map((auction) => (
                                    <option key={auction.id} value={auction.id}>
                                        {auction.auctionDate} {auction.startTime?.substring(0, 5)} - {auction.endTime?.substring(0, 5)} (ID: {auction.id})
                                    </option>
                                ))}
                            </Select>
                            {form.errors.auctionId && (
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.auctionId.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                {t('Admin.Items.fisherman')}
                            </Text>
                            <Select
                                {...form.register('fishermanId')}
                            >
                                <option value="">{t('Admin.Items.placeholder_select_fisherman')}</option>
                                {state.fishermen.map((fisherman) => (
                                    <option key={fisherman.id} value={fisherman.id}>
                                        {fisherman.name}
                                    </option>
                                ))}
                            </Select>
                            {form.errors.fishermanId && (
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.fishermanId.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                {t('Admin.Items.fish_type')}
                            </Text>
                            <Input
                                type="text"
                                {...form.register('fishType')}
                                placeholder={t('Admin.Items.placeholder_fish_type')}
                                error={!!form.errors.fishType}
                                className={css({ _focus: { borderColor: 'orange.500', ringColor: 'orange.500' } })}
                            />
                            {form.errors.fishType && (
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.fishType.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                {t('Admin.Items.quantity')}
                            </Text>
                            <Input
                                type="number"
                                {...form.register('quantity')}
                                placeholder={t('Admin.Items.placeholder_quantity')}
                                error={!!form.errors.quantity}
                                className={css({ _focus: { borderColor: 'orange.500', ringColor: 'orange.500' } })}
                            />
                            {form.errors.quantity && (
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.quantity.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                {t('Admin.Items.unit')}
                            </Text>
                            <Input
                                type="text"
                                {...form.register('unit')}
                                placeholder={t('Admin.Items.placeholder_unit')}
                                error={!!form.errors.unit}
                                className={css({ _focus: { borderColor: 'orange.500', ringColor: 'orange.500' } })}
                            />
                            {form.errors.unit && (
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.unit.message}</Text>
                            )}
                        </Box>
                        <Box className={css({ md: { gridColumn: 'span 2' }, pt: '4' })}>
                            <HStack spacing="2">
                                <Button
                                    type="submit"
                                    disabled={state.isCreating}
                                    width="full"
                                    className={css({ flex: '1' })}
                                    variant="primary"
                                >
                                    {state.isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register)}
                                </Button>
                            </HStack>
                        </Box>
                    </Box>
                </form>
            </Card>
        </Box>
    );
}
