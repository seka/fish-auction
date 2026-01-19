'use client';

import { Suspense } from 'react';
import { useItemPage } from './_hooks/useItemPage';
import { Box, Stack, HStack, Text, Card, Button, Input, Select, Table, Thead, Tbody, Tr, Th, Td } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { EmptyState } from '@/app/_components/atoms/EmptyState';

function ItemsPageContent() {
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

            <Box display="grid" gridTemplateColumns={{ base: '1fr', md: '1fr 2fr' }} gap="8">
                {/* Form Section */}
                <Box>
                    <Card padding="lg">
                        <HStack mb="6">
                            <Box w="2" h="6" bg="orange.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'orange.900' })} fontWeight="bold">
                                {state.editingItem ? t('Admin.Items.edit_item') : t('Admin.Items.register_title')}
                            </Text>
                        </HStack>
                        <form onSubmit={actions.onSubmit}>
                            <Stack spacing="4">
                                <Box>
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
                                    />
                                    {form.errors.unit && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.unit.message}</Text>
                                    )}
                                </Box>
                                <HStack spacing="2" mt="4">
                                    <Button
                                        type="submit"
                                        disabled={state.isCreating || state.isUpdating}
                                        width="full"
                                        variant="primary"
                                    >
                                        {state.isCreating || state.isUpdating
                                            ? t(COMMON_TEXT_KEYS.loading)
                                            : state.editingItem
                                                ? t(COMMON_TEXT_KEYS.save)
                                                : t(COMMON_TEXT_KEYS.register)
                                        }
                                    </Button>
                                    {state.editingItem && (
                                        <Button
                                            type="button"
                                            onClick={actions.onCancelEdit}
                                            width="full"
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
                <Box>
                    <Card padding="none" overflow="hidden">
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200" bg="white" display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap="4">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">{t('Admin.Items.list_title')}</Text>
                            <HStack spacing="2">
                                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>{t('Admin.Items.filter_auction')}</Text>
                                <Select
                                    value={state.filterAuctionId || ''}
                                    onChange={(e) => actions.setFilterAuctionId(e.target.value ? Number(e.target.value) : undefined)}
                                    className={css({ width: 'auto', py: '1' })}
                                >
                                    <option value="">{t('Admin.Items.filter_all')}</option>
                                    {state.auctions.map((auction) => (
                                        <option key={auction.id} value={auction.id}>
                                            {auction.auctionDate} (ID: {auction.id})
                                        </option>
                                    ))}
                                </Select>
                            </HStack>
                        </Box>

                        {state.isItemsLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.loading)}</Box>
                        ) : !state.filterAuctionId ? (
                            <Box p="10" textAlign="center">
                                <Text className={css({ color: 'gray.500' })}>{t('Admin.Items.placeholder_select_auction')}</Text>
                            </Box>
                        ) : state.items.length === 0 ? (
                            <EmptyState
                                message={t(COMMON_TEXT_KEYS.no_data)}
                                icon={<span role="img" aria-label="item">🐟</span>}
                            />
                        ) : (
                            <Box overflowX="auto">
                                <Table>
                                    <Thead>
                                        <Tr>
                                            <Th width="80px">{t('Admin.Items.sort_order')}</Th>
                                            <Th>{t('Admin.Items.fish_type')}</Th>
                                            <Th>{t('Admin.Items.fisherman')}</Th>
                                            <Th>{t('Admin.Items.quantity')}</Th>
                                            <Th className={css({ textAlign: 'right' })}>{t('Admin.Auctions.action')}</Th>
                                        </Tr>
                                    </Thead>
                                    <Tbody>
                                        {state.items.map((item) => {
                                            const fisherman = state.fishermen.find(f => f.id === item.fishermanId);
                                            return (
                                                <Tr key={item.id}>
                                                    <Td>
                                                        <Text fontSize="sm" fontWeight="bold" className={css({ color: 'gray.500' })}>
                                                            #{item.sortOrder}
                                                        </Text>
                                                    </Td>
                                                    <Td>
                                                        <Text fontSize="sm" fontWeight="medium" className={css({ color: 'gray.900' })}>{item.fishType}</Text>
                                                    </Td>
                                                    <Td>
                                                        <Text fontSize="sm" className={css({ color: 'gray.900' })}>{fisherman?.name || `ID: ${item.fishermanId}`}</Text>
                                                    </Td>
                                                    <Td>
                                                        <Text fontSize="sm" className={css({ color: 'gray.900' })}>{item.quantity} {item.unit}</Text>
                                                    </Td>
                                                    <Td className={css({ textAlign: 'right' })}>
                                                        <HStack justify="end" spacing="2">
                                                            <Button
                                                                size="sm"
                                                                variant="outline"
                                                                onClick={() => actions.onEdit(item)}
                                                            >
                                                                {t(COMMON_TEXT_KEYS.edit)}
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                className={css({ bg: 'red.50', color: 'red.600', _hover: { bg: 'red.100' } })}
                                                                onClick={() => actions.onDelete(item.id)}
                                                                disabled={state.isDeleting}
                                                            >
                                                                {t(COMMON_TEXT_KEYS.delete)}
                                                            </Button>
                                                        </HStack>
                                                    </Td>
                                                </Tr>
                                            );
                                        })}
                                    </Tbody>
                                </Table>
                            </Box>
                        )}
                    </Card>
                </Box>
            </Box>
        </Box>
    );
}

export default function AdminItemsPage() {
    return (
        <Suspense fallback={
            <Box maxW="6xl" mx="auto" p="6" textAlign="center">
                <Text>Loading...</Text>
            </Box>
        }>
            <ItemsPageContent />
        </Suspense>
    );
}
