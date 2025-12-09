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
        <Box maxW="6xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                セリ管理
            </Text>

            {state.message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" className={css({ color: 'blue.700' })} p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">通知</Text>
                    <Text>{state.message}</Text>
                </Box>
            )}

            <Box display="grid" gridTemplateColumns={{ base: '1fr', lg: '1fr 2fr' }} gap="8">
                {/* Form Section */}
                <Box>
                    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
                        <HStack mb="6">
                            <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'indigo.900' })} fontWeight="bold">
                                {state.editingAuction ? 'セリ編集' : '新規セリ登録'}
                            </Text>
                        </HStack>
                        <form onSubmit={actions.onSubmit}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        会場
                                    </Text>
                                    <Select
                                        {...form.register('venueId', { valueAsNumber: true })}
                                    >
                                        <option value="">会場を選択してください</option>
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
                                        開催日
                                    </Text>
                                    <Input
                                        type="date"
                                        {...form.register('auctionDate')}
                                    />
                                    {form.errors.auctionDate && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{form.errors.auctionDate.message}</Text>
                                    )}
                                </Box>
                                <Box display="grid" gridTemplateColumns="repeat(2, 1fr)" gap="4">
                                    <Box>
                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                            開始時間
                                        </Text>
                                        <Input
                                            type="time"
                                            {...form.register('startTime')}
                                        />
                                    </Box>
                                    <Box>
                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                            終了時間
                                        </Text>
                                        <Input
                                            type="time"
                                            {...form.register('endTime')}
                                        />
                                    </Box>
                                </Box>

                                <HStack spacing="2" pt="4">
                                    {state.editingAuction && (
                                        <Button
                                            type="button"
                                            variant="outline"
                                            onClick={actions.onCancelEdit}
                                            disabled={state.isCreating || state.isUpdating}
                                            className={css({ flex: '1' })}
                                        >
                                            {t(COMMON_TEXT_KEYS.cancel)}
                                        </Button>
                                    )}
                                    <Button
                                        type="submit"
                                        disabled={state.isCreating || state.isUpdating}
                                        width="full"
                                        className={css({ flex: '1' })}
                                        variant="primary"
                                    >
                                        {state.editingAuction ? (state.isUpdating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.update)) : (state.isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register))}
                                    </Button>
                                </HStack>
                            </Stack>
                        </form>
                    </Card>
                </Box>

                {/* List Section */}
                <Box>
                    <Card padding="none" overflow="hidden">
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200" display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap="4">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">セリ一覧</Text>
                            <HStack spacing="2">
                                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>会場絞り込み:</Text>
                                <Select
                                    value={state.filterVenueId || ''}
                                    onChange={(e) => actions.setFilterVenueId(e.target.value ? Number(e.target.value) : undefined)}
                                    className={css({ width: 'auto', py: '1' })}
                                >
                                    <option value="">すべて</option>
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
                                            <Th>開催日 / 時間</Th>
                                            <Th>会場</Th>
                                            <Th>ステータス</Th>
                                            <Th className={css({ textAlign: 'right' })}>操作</Th>
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
                                                                    開始
                                                                </Button>
                                                            )}
                                                            {auction.status === 'in_progress' && (
                                                                <Button
                                                                    size="sm"
                                                                    onClick={() => actions.onStatusChange(auction.id, 'completed')}
                                                                    disabled={state.isUpdatingStatus}
                                                                    className={css({ color: 'blue.600', bg: 'blue.50', borderColor: 'transparent', _hover: { bg: 'blue.100', color: 'blue.900' } })}
                                                                >
                                                                    終了
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
