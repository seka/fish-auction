'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { itemSchema, ItemFormData } from '@/src/models/schemas/admin';
import { useItemMutations } from './_hooks/useItem';
import { useFishermen } from '../fishermen/_hooks/useFisherman';
import { useAuctionQuery } from '@/src/repositories/auction';
import { Box, Stack, HStack, Text, Card, Button, Input, Select } from '@/src/core/ui';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { useTranslations } from 'next-intl';
import { css } from 'styled-system/css';
import { styled } from 'styled-system/jsx';


export default function AdminItemsPage() {
    const t = useTranslations();
    const [message, setMessage] = useState('');

    const { fishermen } = useFishermen();
    const { auctions } = useAuctionQuery({});
    const { createItem, isCreating } = useItemMutations();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<ItemFormData>({
        resolver: zodResolver(itemSchema),
    });

    const onSubmit = async (data: ItemFormData) => {
        try {
            await createItem({
                auctionId: parseInt(data.auctionId),
                fishermanId: parseInt(data.fishermanId),
                fishType: data.fishType,
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
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                出品管理
            </Text>

            {message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" className={css({ color: 'blue.700' })} p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">通知</Text>
                    <Text>{message}</Text>
                </Box>
            )}

            <Card padding="lg">
                <HStack mb="6">
                    <Box w="2" h="6" bg="orange.500" mr="3" borderRadius="full" />
                    <Text as="h2" variant="h4" className={css({ color: 'orange.900' })} fontWeight="bold">
                        新規出品登録
                    </Text>
                </HStack>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Box display="grid" gridTemplateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap="6">
                        <Box className={css({ md: { gridColumn: 'span 2' } })}>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                セリ
                            </Text>
                            <Select
                                {...register('auctionId')}
                            >
                                <option value="">セリを選択してください</option>
                                {auctions.map((auction) => (
                                    <option key={auction.id} value={auction.id}>
                                        {auction.auctionDate} {auction.startTime?.substring(0, 5)} - {auction.endTime?.substring(0, 5)} (ID: {auction.id})
                                    </option>
                                ))}
                            </Select>
                            {errors.auctionId && (
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.auctionId.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
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
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.fishermanId.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
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
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.fishType.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
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
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.quantity.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
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
                                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.unit.message}</Text>
                            )}
                        </Box>
                        <Box className={css({ md: { gridColumn: 'span 2' }, pt: '4' })}>
                            <HStack spacing="2">
                                <Button
                                    type="submit"
                                    disabled={isCreating}
                                    width="full"
                                    className={css({ flex: '1' })}
                                    variant="primary"
                                >
                                    {isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register)}
                                </Button>
                            </HStack>
                        </Box>
                    </Box>
                </form>
            </Card>
        </Box>
    );
}
