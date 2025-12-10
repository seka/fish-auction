'use client';

import { useAuctionPage } from './_hooks/useAuctionPage';
import { AUCTION_STATUS_KEYS, AuctionStatus } from '@/src/core/assets/status';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { Box, Stack, HStack, Text, Card, Button, Input, Select, Table, Thead, Tbody, Tr, Th, Td } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function AuctionsPage() {
    const { state, actions, form, t } = useAuctionPage();

    const getStatusBadge = (status: string) => {
        const baseStyle = { fontSize: 'xs', fontWeight: 'medium', px: '2.5', py: '0.5', borderRadius: 'md' };
        const statusKey = status as AuctionStatus;
        const label = AUCTION_STATUS_KEYS[statusKey] ? t(AUCTION_STATUS_KEYS[statusKey]) : status;

        switch (status) {
            case 'scheduled':
                return <span className={css(baseStyle, { bg: 'blue.100', color: 'blue.800' })}>{label}</span>;
            case 'in_progress':
                return <span className={css(baseStyle, { bg: 'orange.100', color: 'orange.800', animation: 'pulse 2s infinite' })}>🔥 {label}</span>;
            case 'completed':
                return <span className={css(baseStyle, { bg: 'green.100', color: 'green.800' })}>{label}</span>;
            case 'cancelled':
                return <span className={css(baseStyle, { bg: 'red.100', color: 'red.800' })}>{label}</span>;
            default:
                return <span className={css(baseStyle, { bg: 'gray.100', color: 'gray.800' })}>{label}</span>;
        }
    };

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                {t('Admin.Auctions.title')}
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
                            <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'indigo.900' })} fontWeight="bold">
                                {state.editingAuction ? t('Admin.Auctions.edit_title') : t('Admin.Auctions.register_title')}
                            </Text>
                        </HStack>
                        <form onSubmit={actions.onSubmit}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Auctions.venue')}
                                    </Text>
                                    <Select
                                        {...form.register('venueId')}
                                    >
                                        <option value="">{t('Admin.Auctions.placeholder_select_venue')}</option>
                                        {state.venues.map((venue) => (
                                            <option key={venue.id} value={venue.id}>
                                                {venue.name}
                                            </option>
                                        ))}
                                    </Select>
                                    {form.errors.venueId && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.venueId.message}</Text>
                                    )}
                                </Box>
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        {t('Admin.Auctions.date')}
                                    </Text>
                                    <Input
                                        type="date"
                                        {...form.register('auctionDate')}
                                        error={!!form.errors.auctionDate}
                                    />
                                    {form.errors.auctionDate && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.auctionDate.message}</Text>
                                    )}
                                </Box>
                                <HStack spacing="4">
                                    <Box flex="1">
                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                            {t('Admin.Auctions.start_time')}
                                        </Text>
                                        <Input
                                            type="time"
                                            {...form.register('startTime')}
                                            error={!!form.errors.startTime}
                                        />
                                    </Box>
                                    <Box flex="1">
                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                            {t('Admin.Auctions.end_time')}
                                        </Text>
                                        <Input
                                            type="time"
                                            {...form.register('endTime')}
                                            error={!!form.errors.endTime}
                                        />
                                    </Box>
                                </HStack>

                                <HStack spacing="2" pt="4">
                                    <Button
                                        type="submit"
                                        disabled={state.isCreating || state.isUpdating}
                                        width="full"
                                        className={css({ flex: '1' })}
                                        variant="primary"
                                    >
                                        {state.editingAuction ? (state.isUpdating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.update)) : (state.isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register))}
                                    </Button>
                                    {state.editingAuction && (
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
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200" bg="white" display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap="4">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">{t('Admin.Auctions.list_title')}</Text>
                            <HStack spacing="2">
                                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>{t('Admin.Auctions.filter_venue')}</Text>
                                <Select
                                    value={state.filterVenueId || ''}
                                    onChange={(e) => actions.setFilterVenueId(e.target.value ? Number(e.target.value) : undefined)}
                                    className={css({ width: 'auto', py: '1' })}
                                >
                                    <option value="">{t('Admin.Auctions.filter_all')}</option>
                                    {state.venues.map((venue) => (
                                        <option key={venue.id} value={venue.id}>
                                            {venue.name}
                                        </option>
                                    ))}
                                </Select>
                            </HStack>
                        </Box>
                        {state.isLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.loading)}</Box>
                        ) : state.auctions.length === 0 ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.no_data)}</Box>
                        ) : (
                            <Box overflowX="auto">
                                <Table>
                                    <Thead>
                                        <Tr>
                                            <Th>{t('Admin.Auctions.date_time')}</Th>
                                            <Th>{t('Admin.Auctions.venue')}</Th>
                                            <Th>{t('Admin.Auctions.status')}</Th>
                                            <Th className={css({ textAlign: 'right' })}>{t('Admin.Auctions.action')}</Th>
                                        </Tr>
                                    </Thead>
                                    <Tbody>
                                        {state.auctions.map((auction) => {
                                            const venue = state.venues.find(v => v.id === auction.venueId);
                                            return (
                                                <Tr key={auction.id}>
                                                    <Td>
                                                        <Text fontSize="sm" fontWeight="medium" className={css({ color: 'gray.900' })}>{auction.auctionDate}</Text>
                                                        <Text fontSize="sm" className={css({ color: 'gray.500' })}>
                                                            {auction.startTime ? auction.startTime.substring(0, 5) : '--:--'} - {auction.endTime ? auction.endTime.substring(0, 5) : '--:--'}
                                                        </Text>
                                                    </Td>
                                                    <Td>
                                                        <Text fontSize="sm" className={css({ color: 'gray.900' })}>{venue?.name || `ID: ${auction.venueId}`}</Text>
                                                    </Td>
                                                    <Td>
                                                        {getStatusBadge(auction.status)}
                                                    </Td>
                                                    <Td className={css({ textAlign: 'right' })}>
                                                        <HStack justify="end" spacing="2">
                                                            {auction.status === 'scheduled' && (
                                                                <Button
                                                                    size="sm"
                                                                    onClick={() => actions.onStatusChange(auction.id, 'in_progress')}
                                                                    disabled={state.isUpdatingStatus}
                                                                    className={css({ color: 'green.600', bg: 'green.50', borderColor: 'transparent', _hover: { bg: 'green.100', color: 'green.900' } })}
                                                                >
                                                                    {t('Admin.Auctions.start')}
                                                                </Button>
                                                            )}
                                                            {auction.status === 'in_progress' && (
                                                                <Button
                                                                    size="sm"
                                                                    onClick={() => actions.onStatusChange(auction.id, 'completed')}
                                                                    disabled={state.isUpdatingStatus}
                                                                    className={css({ color: 'blue.600', bg: 'blue.50', borderColor: 'transparent', _hover: { bg: 'blue.100', color: 'blue.900' } })}
                                                                >
                                                                    {t('Admin.Auctions.finish')}
                                                                </Button>
                                                            )}
                                                            <Button size="sm" variant="outline" onClick={() => actions.onEdit(auction)}>
                                                                {t(COMMON_TEXT_KEYS.edit)}
                                                            </Button>
                                                            <Button size="sm" className={css({ bg: 'red.50', color: 'red.600', _hover: { bg: 'red.100' } })} onClick={() => actions.onDelete(auction.id)}>
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
