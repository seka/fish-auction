      <Td className={css({ fontSize: 'sm', color: 'gray.500', fontFamily: 'mono' })}>
        {invoice.buyerId}
      </Td>
      <Td className={css({ fontSize: 'sm', fontWeight: 'bold', color: 'gray.900' })}>
        {invoice.buyerName}
      </Td>
      <Td
        className={css({
          textAlign: 'right',
          fontWeight: 'bold',
          color: 'indigo.700',
          fontSize: 'lg',
        })}
      >
        ¥{invoice.totalAmount.toLocaleString()}
      </Td>
    </Tr>
  );
};
