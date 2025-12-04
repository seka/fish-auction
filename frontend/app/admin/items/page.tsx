'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { itemSchema, ItemFormData } from '@/src/models/schemas/admin';
import { useItemMutations } from './_hooks/useItem';
import { useFishermen } from '../fishermen/_hooks/useFisherman';
import { useAuctions } from '../auctions/_hooks/useAuction';
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
        border: '1px solid',
        borderColor: 'gray.300',
        borderRadius: 'md',
        fontSize: 'sm',
        outline: 'none',
        transition: 'border-color 0.2s',
        appearance: 'none', // Remove default arrow
        backgroundImage: 'url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns=\'http://www.w3.org/2000/svg\' fill=\'none\' viewBox=\'0 0 20 20\'%3E%3Cpath stroke=\'%236b7280\' stroke-linecap=\'round\' stroke-linejoin=\'round\' stroke-width=\'1.5\' d=\'M6 8l4 4 4-4\'/%3E%3C/svg%3E")',
        backgroundPosition: 'right 0.5rem center',
        backgroundRepeat: 'no-repeat',
        backgroundSize: '1.5em 1.5em',
        paddingRight: '2.5rem',
        _focus: {
            borderColor: 'orange.500',
            ring: '1px',
            ringColor: 'orange.500',
        },
        _disabled: {
            bg: 'gray.50',
            cursor: 'not-allowed',
        },
    }
});

export default function ItemsPage() {
    const [message, setMessage] = useState('');

    const { fishermen } = useFishermen();
    const { auctions } = useAuctions({});
    const { createItem, isCreating } = useItemMutations();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<ItemFormData>({
        resolver: zodResolver(itemSchema),
    });

    const onSubmit = async (data: ItemFormData) => {
        try {
            await createItem({
                auction_id: parseInt(data.auctionId),
                fisherman_id: parseInt(data.fishermanId),
                fish_type: data.fishType,
                quantity: parseInt(data.quantity),
                unit: data.unit,
            });
            setMessage('出品を登録しました');
            reset();
        } catch (e) {
            console.error(e);
            setMessage('登録に失敗しました');
        }
    };

    return (
        <Box maxW="6xl" mx="auto" p="6">
            <Text as="h1" variant="h2" color="gray.800" mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                出品管理
            </Text>

            {message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" color="blue.700" p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">通知</Text>
                    <Text>{message}</Text>
                </Box>
            )}

            <Card p="8">
                <HStack mb="6">
                    <Box w="2" h="6" bg="orange.500" mr="3" borderRadius="full" />
                    <Text as="h2" variant="h4" color="orange.900" fontWeight="bold">
                        新規出品登録
                    </Text>
                </HStack>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Box display="grid" gridTemplateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap="6">
                        <Box className={css({ md: { gridColumn: 'span 2' } })}>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" color="gray.700" mb="1">
                                セリ
                            </Text>
                            <Select
                                {...register('auctionId')}
                            >
                                <option value="">セリを選択してください</option>
                                {auctions.map((auction) => (
                                    <option key={auction.id} value={auction.id}>
                                        {auction.auction_date} {auction.start_time?.substring(0, 5)} - {auction.end_time?.substring(0, 5)} (ID: {auction.id})
                                    </option>
                                ))}
                            </Select>
                            {errors.auctionId && (
                                <Text color="red.500" fontSize="sm" mt="1">{errors.auctionId.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" color="gray.700" mb="1">
                                漁師
                            </Text>
                            <Select
                                {...register('fishermanId')}
                            >
                                <option value="">漁師を選択してください</option>
                                {fishermen.map((fisherman) => (
                                    <option key={fisherman.id} value={fisherman.id}>
                                        {fisherman.name}
                                    </option>
                                ))}
                            </Select>
                            {errors.fishermanId && (
                                <Text color="red.500" fontSize="sm" mt="1">{errors.fishermanId.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" color="gray.700" mb="1">
                                魚種
                            </Text>
                            <Input
                                type="text"
                                {...register('fishType')}
                                placeholder="例: マグロ"
                                error={!!errors.fishType}
                                className={css({ _focus: { borderColor: 'orange.500', ringColor: 'orange.500' } })}
                            />
                            {errors.fishType && (
                                <Text color="red.500" fontSize="sm" mt="1">{errors.fishType.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" color="gray.700" mb="1">
                                数量
                            </Text>
                            <Input
                                type="number"
                                {...register('quantity')}
                                placeholder="例: 10"
                                error={!!errors.quantity}
                                className={css({ _focus: { borderColor: 'orange.500', ringColor: 'orange.500' } })}
                            />
                            {errors.quantity && (
                                <Text color="red.500" fontSize="sm" mt="1">{errors.quantity.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" color="gray.700" mb="1">
                                単位
                            </Text>
                            <Input
                                type="text"
                                {...register('unit')}
                                placeholder="例: kg, 匹, 箱"
                                error={!!errors.unit}
                                className={css({ _focus: { borderColor: 'orange.500', ringColor: 'orange.500' } })}
                            />
                            {errors.unit && (
                                <Text color="red.500" fontSize="sm" mt="1">{errors.unit.message}</Text>
                            )}
                        </Box>
                        <Box className={css({ md: { gridColumn: 'span 2' }, pt: '4' })}>
                            <Button
                                type="submit"
                                disabled={isCreating}
                                width="full"
                                className={css({
                                    bg: 'orange.600',
                                    _hover: { bg: 'orange.700' },
                                    _focus: { ringColor: 'orange.500' }
                                })}
                            >
                                {isCreating ? '出品中...' : '出品する'}
                            </Button>
                        </Box>
                    </Box>
                </form>
            </Card>
        </Box>
    );
}
