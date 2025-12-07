'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { auctionSchema, AuctionFormData } from '@/src/models/schemas/auction';
import { useAuctions, useAuctionMutations } from './_hooks/useAuction';
import { useVenues } from '../venues/_hooks/useVenue';
import { Auction } from '@/src/models/auction';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { styled } from 'styled-system/jsx';

// Select component with similar styling to Input
const Select = styled('select', {
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
        appearance: 'none',
        backgroundImage: 'url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns=\'http://www.w3.org/2000/svg\' fill=\'none\' viewBox=\'0 0 20 20\'%3E%3Cpath stroke=\'%236b7280\' stroke-linecap=\'round\' stroke-linejoin=\'round\' stroke-width=\'1.5\' d=\'M6 8l4 4 4-4\'/%3E%3C/svg%3E")',
        backgroundPosition: 'right 0.5rem center',
        backgroundRepeat: 'no-repeat',
        backgroundSize: '1.5em 1.5em',
        paddingRight: '2.5rem',
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

const Table = styled('table', { base: { minW: 'full', divideY: '1px', divideColor: 'gray.200' } });
const Thead = styled('thead', { base: { bg: 'gray.50' } });
const Tbody = styled('tbody', { base: { bg: 'white', divideY: '1px', divideColor: 'gray.200' } });
const Tr = styled('tr', { base: { _hover: { bg: 'gray.50' } } });
const Th = styled('th', { base: { px: '6', py: '3', textAlign: 'left', fontSize: 'xs', fontWeight: 'medium', color: 'gray.500', textTransform: 'uppercase', letterSpacing: 'wider' } });
const Td = styled('td', { base: { px: '6', py: '4', whiteSpace: 'nowrap' } });

export default function AuctionsPage() {
    const [message, setMessage] = useState('');
    const [editingAuction, setEditingAuction] = useState<Auction | null>(null);
    const [filterVenueId, setFilterVenueId] = useState<number | undefined>(undefined);

    const { venues } = useVenues();
    const { auctions, isLoading } = useAuctions({ venueId: filterVenueId });
    const { createAuction, updateAuction, updateStatus, deleteAuction, isCreating, isUpdating, isUpdatingStatus, isDeleting } = useAuctionMutations();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<AuctionFormData>({
        resolver: zodResolver(auctionSchema),
    });

    const onSubmit = async (data: AuctionFormData) => {
        try {
            const payload = {
                ...data,
                venueId: Number(data.venueId),
            };

            if (editingAuction) {
                await updateAuction({ id: editingAuction.id, data: payload });
                setMessage('セリ情報を更新しました');
                setEditingAuction(null);
            } else {
                await createAuction(payload);
                setMessage('セリを作成しました');
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage('エラーが発生しました');
        }
    };

    const onEdit = (auction: Auction) => {
        setEditingAuction(auction);
        setValue('venueId', auction.venueId);
        setValue('auctionDate', auction.auctionDate);
        setValue('startTime', auction.startTime || '');
        setValue('endTime', auction.endTime || '');
        setValue('status', auction.status);
    };

    const onCancelEdit = () => {
        setEditingAuction(null);
        reset();
    };

    const onDelete = async (id: number) => {
        if (confirm('本当に削除しますか？')) {
            try {
                await deleteAuction(id);
                setMessage('セリを削除しました');
            } catch (e) {
                console.error(e);
                setMessage('削除に失敗しました');
            }
        }
    };

    const onStatusChange = async (id: number, status: string) => {
        try {
            await updateStatus({ id, status });
            setMessage(`ステータスを ${status} に更新しました`);
        } catch (e) {
            console.error(e);
            setMessage('ステータス更新に失敗しました');
        }
    };

    const getStatusBadge = (status: string) => {
        const baseStyle = { fontSize: 'xs', fontWeight: 'medium', px: '2.5', py: '0.5', borderRadius: 'md' };
        switch (status) {
            case 'scheduled':
                return <span className={css(baseStyle, { bg: 'blue.100', color: 'blue.800' })}>予定</span>;
            case 'in_progress':
                return <span className={css(baseStyle, { bg: 'green.100', color: 'green.800', animation: 'pulse 2s infinite' })}>開催中</span>;
            case 'completed':
                return <span className={css(baseStyle, { bg: 'gray.100', color: 'gray.800' })}>終了</span>;
            case 'cancelled':
                return <span className={css(baseStyle, { bg: 'red.100', color: 'red.800' })}>中止</span>;
            default:
                return <span className={css(baseStyle, { bg: 'gray.100', color: 'gray.800' })}>{status}</span>;
        }
    };

    return (
        <Box maxW="6xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                セリ管理
            </Text>

            {message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" className={css({ color: 'blue.700' })} p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">通知</Text>
                    <Text>{message}</Text>
                </Box>
            )}

            <Box display="grid" gridTemplateColumns={{ base: '1fr', lg: '1fr 2fr' }} gap="8">
                {/* Form Section */}
                <Box>
                    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
                        <HStack mb="6">
                            <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'indigo.900' })} fontWeight="bold">
                                {editingAuction ? 'セリ編集' : '新規セリ登録'}
                            </Text>
                        </HStack>
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        会場
                                    </Text>
                                    <Select
                                        {...register('venueId', { valueAsNumber: true })}
                                    >
                                        <option value="">会場を選択してください</option>
                                        {venues.map((venue) => (
                                            <option key={venue.id} value={venue.id}>
                                                {venue.name}
                                            </option>
                                        ))}
                                    </Select>
                                    {errors.venueId && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.venueId.message}</Text>
                                    )}
                                </Box>
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        開催日
                                    </Text>
                                    <Input
                                        type="date"
                                        {...register('auctionDate')}
                                    />
                                    {errors.auctionDate && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.auctionDate.message}</Text>
                                    )}
                                </Box>
                                <Box display="grid" gridTemplateColumns="repeat(2, 1fr)" gap="4">
                                    <Box>
                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                            開始時間
                                        </Text>
                                        <Input
                                            type="time"
                                            {...register('startTime')}
                                        />
                                    </Box>
                                    <Box>
                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                            終了時間
                                        </Text>
                                        <Input
                                            type="time"
                                            {...register('endTime')}
                                        />
                                    </Box>
                                </Box>

                                <HStack spacing="2" pt="4">
                                    <Button
                                        type="submit"
                                        disabled={isCreating || isUpdating}
                                        width="full"
                                        className={css({ flex: '1' })}
                                        variant="primary"
                                    >
                                        {editingAuction ? (isUpdating ? '更新中...' : '更新する') : (isCreating ? '登録中...' : '登録する')}
                                    </Button>
                                    {editingAuction && (
                                        <Button
                                            type="button"
                                            onClick={onCancelEdit}
                                            variant="outline"
                                        >
                                            キャンセル
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
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200" display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap="4">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">セリ一覧</Text>
                            <HStack spacing="2">
                                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>会場絞り込み:</Text>
                                <Select
                                    value={filterVenueId || ''}
                                    onChange={(e) => setFilterVenueId(e.target.value ? Number(e.target.value) : undefined)}
                                    className={css({ width: 'auto', py: '1' })}
                                >
                                    <option value="">すべて</option>
                                    {venues.map((venue) => (
                                        <option key={venue.id} value={venue.id}>
                                            {venue.name}
                                        </option>
                                    ))}
                                </Select>
                            </HStack>
                        </Box>
                        {isLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>読み込み中...</Box>
                        ) : auctions.length === 0 ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>セリが登録されていません</Box>
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
                                        {auctions.map((auction) => {
                                            const venue = venues.find(v => v.id === auction.venueId);
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
                                                                    onClick={() => onStatusChange(auction.id, 'in_progress')}
                                                                    disabled={isUpdatingStatus}
                                                                    className={css({ color: 'green.600', bg: 'green.50', borderColor: 'transparent', _hover: { bg: 'green.100', color: 'green.900' } })}
                                                                >
                                                                    開始
                                                                </Button>
                                                            )}
                                                            {auction.status === 'in_progress' && (
                                                                <Button
                                                                    size="sm"
                                                                    onClick={() => onStatusChange(auction.id, 'completed')}
                                                                    disabled={isUpdatingStatus}
                                                                    className={css({ color: 'blue.600', bg: 'blue.50', borderColor: 'transparent', _hover: { bg: 'blue.100', color: 'blue.900' } })}
                                                                >
                                                                    終了
                                                                </Button>
                                                            )}
                                                            <Button
                                                                size="sm"
                                                                onClick={() => onEdit(auction)}
                                                                className={css({ color: 'indigo.600', bg: 'indigo.50', borderColor: 'transparent', _hover: { bg: 'indigo.100', color: 'indigo.900' } })}
                                                            >
                                                                編集
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                onClick={() => onDelete(auction.id)}
                                                                disabled={isDeleting}
                                                                className={css({ color: 'red.600', bg: 'red.50', borderColor: 'transparent', _hover: { bg: 'red.100', color: 'red.900' } })}
                                                            >
                                                                削除
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
