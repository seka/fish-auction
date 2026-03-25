import { ItemList } from './ItemList';

export const ItemManagementContainer = () => {
  const { state, form, actions, t } = useItemManagement();

  return (
    <Box maxW="6xl" mx="auto" p="6">
      <Text
        as="h1"
        variant="h2"
        className={css({ color: 'gray.800' })}
        mb="8"
        pb="4"
        borderBottom="1px solid"
        borderColor="gray.200"
      >
        {t('Admin.Items.title')}
      </Text>

      {state.message && (
        <Box
          bg="blue.50"
          borderLeft="4px solid"
          borderColor="blue.500"
          className={css({ color: 'blue.700' })}
          p="4"
          mb="8"
          borderRadius="sm"
          shadow="sm"
          role="alert"
        >
          <Text fontWeight="bold">{t('Common.notification')}</Text>
          <Text>{state.message}</Text>
        </Box>
      )}

      <Box display="grid" gridTemplateColumns={{ base: '1fr', md: '1fr 2fr' }} gap="8">
        {/* Form Section */}
        <Box>
          <ItemForm
            form={form}
            onSubmit={actions.onSubmit}
            onCancelEdit={actions.onCancelEdit}
            isCreating={state.isCreating}
            isUpdating={state.isUpdating}
            editingItem={state.editingItem}
            auctions={state.auctions}
            fishermen={state.fishermen}
          />
        </Box>

        {/* List Section */}
        <Box>
          <Card padding="none" overflow="hidden">
            <Box
              p="6"
              borderBottom="1px solid"
              borderColor="gray.200"
              bg="white"
              display="flex"
              justifyContent="space-between"
              alignItems="center"
              flexWrap="wrap"
              gap="4"
            >
              <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">
                {t('Admin.Items.list_title')}
              </Text>
              <HStack spacing="2">
                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>
                  {t('Admin.Items.filter_auction')}
                </Text>
                <Select
                  value={state.filterAuctionId || ''}
                  onChange={(e) =>
                    actions.setFilterAuctionId(e.target.value ? Number(e.target.value) : undefined)
                  }
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

            <ItemList
              items={state.items}
              fishermen={state.fishermen}
              isItemsLoading={state.isItemsLoading}
              filterAuctionId={state.filterAuctionId}
              isDeleting={state.isDeleting}
              onDragEnd={actions.onDragEnd}
              onEdit={actions.onEdit}
              onDelete={actions.onDelete}
              t={t}
            />
          </Card>
        </Box>
      </Box>
    </Box>
  );
};
